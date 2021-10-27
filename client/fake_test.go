package client

import (
	"fmt"
	"testing"

	"github.com/nono/cozy-desktop-experiments/state/remote"
	"github.com/nono/cozy-desktop-experiments/state/types"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestFakeClient(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		client := NewFake("http://cozy.localhost:8080/").(*Fake)
		client.AddInitialDocs()
		require.NoError(t, client.CheckInvariants())
	})

	t.Run("basic", func(t *testing.T) {
		client := NewFake("http://cozy.localhost:8080/").(*Fake)
		client.AddInitialDocs()
		foo, err := client.CreateDir(remote.RootID, "foo")
		require.NoError(t, err, "CreateDir foo")
		bar, err := client.CreateDir(foo.ID, "bar")
		require.NoError(t, err, "CreateDir bar")
		_, err = client.CreateDir(bar.ID, "baz")
		require.NoError(t, err, "CreateDir baz")
		require.NoError(t, client.CheckInvariants())
	})

	t.Run("invalid name", func(t *testing.T) {
		client := NewFake("http://cozy.localhost:8080/").(*Fake)
		client.AddInitialDocs()
		_, err := client.CreateDir(remote.RootID, "foo/")
		require.Error(t, err)
		_, err = client.CreateDir(remote.RootID, "֏\ufeff+$ª!/")
		require.Error(t, err)
		_, err = client.CreateDir(remote.RootID, "foo\nbar")
		require.Error(t, err)
	})
}

func TestConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

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
	seq     *remote.Seq
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
	changes, err := cmp.client.Changes(nil)
	require.NoError(t, err)
	cmp.fake = NewFake(addr).(*Fake)
	cmp.fake.AddInitialDocs(changes.Docs...)
	cmp.fake.MatchSequence(changes.Seq)
	cmp.parents = []*remote.Doc{
		{ID: remote.RootID, Type: types.DirType},
	}
}

func (cmp *cmpClient) Cleanup() {
	_ = cmp.inst.Remove()
}

func (cmp *cmpClient) Changes(t *rapid.T) {
	left, errl := cmp.client.Changes(cmp.seq)
	right, errr := cmp.fake.Changes(cmp.seq)
	require.Equal(t, errl == nil, errr == nil)
	if errl == nil && errr == nil {
		require.Equal(t, len(left.Docs), len(right.Docs))
		// TODO not the same revisions, and probably not the same order
		// for i := range left.Docs {
		// 	require.Equal(t, left.Docs[i], right.Docs[i])
		// }
		require.Equal(t, left.Pending, right.Pending)
		require.Equal(t, left.Seq.ExtractGeneration(), right.Seq.ExtractGeneration())
		cmp.seq = &left.Seq
	}
}

func (cmp *cmpClient) CreateDir(t *rapid.T) {
	parent := rapid.SampledFrom(cmp.parents).Draw(t, "parent").(*remote.Doc)
	name := rapid.String().Draw(t, "name").(string)
	docl, errl := cmp.client.CreateDir(parent.ID, name)
	if errl == nil {
		cmp.fake.GenerateID = func() remote.ID { return docl.ID }
		cmp.fake.GenerateRev = func(id remote.ID, gen int) remote.Rev {
			if id == docl.ID {
				return docl.Rev
			}
			return newRev(id, gen)
		}
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
	docl, errl := cmp.client.Trash(dir)
	if errl == nil {
		cmp.fake.GenerateRev = func(id remote.ID, gen int) remote.Rev {
			if id == docl.ID {
				return docl.Rev
			}
			return newRev(id, gen)
		}
		cmp.fake.ConflictName = func(id remote.ID, name string) string {
			if id == docl.ID {
				return docl.Name
			}
			return conflictName(id, name)
		}
	}
	docr, errr := cmp.fake.Trash(dir)
	require.Equal(t, errl == nil, errr == nil)
	if errl == nil && errr == nil {
		require.Equal(t, docl.ID, docr.ID)
		require.Equal(t, docl.Rev, docr.Rev)
		require.Equal(t, docl.Type, docr.Type)
		require.Equal(t, docl.Name, docr.Name)
		require.Equal(t, docl.DirID, docr.DirID)
	}
}

func (cmp *cmpClient) Check(t *rapid.T) {
	require.NoError(t, cmp.fake.CheckInvariants())

	left := cmp.client.DocsByID()
	right := cmp.fake.DocsByID()
	require.Equal(t, len(left), len(right))
	for id := range left {
		docl := left[id]
		docr := right[id]
		require.NotNil(t, docr)
		require.Equal(t, docl.ID, docr.ID)
		// TODO we need to inject the right revisions for Trash
		// require.Equal(t, docl.Rev, docr.Rev)
		require.Equal(t, docl.Type, docr.Type)
		require.Equal(t, docl.Name, docr.Name)
		require.Equal(t, docl.DirID, docr.DirID)
	}
}
