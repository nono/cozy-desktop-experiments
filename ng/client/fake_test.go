package client

import (
	"fmt"
	"testing"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestFakeClient(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		client := NewFake("http://cozy.localhost:8080/")
		require.NoError(t, client.(*Fake).CheckInvariants())
	})

	t.Run("basic", func(t *testing.T) {
		client := NewFake("http://cozy.localhost:8080/")
		foo, err := client.CreateDir(remote.RootID, "foo")
		require.NoError(t, err, "CreateDir foo")
		bar, err := client.CreateDir(foo.ID, "bar")
		require.NoError(t, err, "CreateDir bar")
		_, err = client.CreateDir(bar.ID, "baz")
		require.NoError(t, err, "CreateDir baz")
		require.NoError(t, client.(*Fake).CheckInvariants())
	})

	t.Run("invalid name", func(t *testing.T) {
		client := NewFake("http://cozy.localhost:8080/")
		_, err := client.CreateDir(remote.RootID, "foo/")
		require.Error(t, err)
		_, err = client.CreateDir(remote.RootID, "֏\ufeff+$ª!/")
		require.Error(t, err)
		_, err = client.CreateDir(remote.RootID, "foo\nbar")
		require.Error(t, err)
	})
}

func TestConsistency(t *testing.T) {
	stack = NewStack(8081, t.TempDir())
	require.NoError(t, stack.Start())
	t.Cleanup(func() {
		_ = stack.Stop()
	})
	rapid.Check(t, rapid.Run(&cmpClient{}))
}

var stack *Stack
var count int

type cmpClient struct {
	inst    *Instance
	client  *Client
	fake    *Fake
	parents []*remote.Doc
}

func (cmp *cmpClient) Init(t *rapid.T) {
	count++
	inst, err := stack.CreateInstance(fmt.Sprintf("test%03d", count))
	require.NoError(t, err)
	cmp.inst = inst
	addr := inst.Address()
	cmp.client = New(addr).(*Client)
	require.NoError(t, cmp.client.Register())
	require.NoError(t, inst.CreateAccessToken(cmp.client))
	cmp.fake = NewFake(addr).(*Fake)
	cmp.parents = []*remote.Doc{
		{ID: remote.RootID, Type: remote.Directory},
	}
}

func (cmp *cmpClient) Cleanup() {
	_ = cmp.inst.Remove()
}

// TODO compare calls to the changes feed

func (cmp *cmpClient) CreateDir(t *rapid.T) {
	parent := rapid.SampledFrom(cmp.parents).Draw(t, "parent").(*remote.Doc)
	name := rapid.String().Draw(t, "name").(string)
	docl, errl := cmp.client.CreateDir(parent.ID, name)
	if errl == nil {
		cmp.fake.GenerateID = func() remote.ID { return docl.ID }
		cmp.fake.GenerateRev = func(_gen int) remote.Rev { return docl.Rev }
	}
	docr, errr := cmp.fake.CreateDir(parent.ID, name)
	require.Equal(t, errl == nil, errr == nil)
	if errl == nil && errr == nil {
		require.Equal(t, docl.ID, docr.ID)
		require.Equal(t, docl.Rev, docr.Rev)
		require.Equal(t, docl.Type, docr.Type)
		require.Equal(t, docl.Name, docr.Name)
		require.Equal(t, docl.DirID, docr.DirID)
		cmp.parents = append(cmp.parents, docl)
	}
}

func (cmp *cmpClient) Trash(t *rapid.T) {
	dir := rapid.SampledFrom(cmp.parents).Draw(t, "dir").(*remote.Doc)
	if dir.ID == remote.RootID {
		return
	}
	_, errl := cmp.client.Trash(dir)
	_, errr := cmp.fake.Trash(dir)
	require.Equal(t, errl == nil, errr == nil)
	// TODO compare docl and docr
}

func (cmp *cmpClient) Check(t *rapid.T) {
	require.NoError(t, cmp.fake.CheckInvariants())
	// TODO compare the trees
}
