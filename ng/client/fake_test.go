package client

import (
	"testing"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
	"github.com/stretchr/testify/require"
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
}
