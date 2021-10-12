package localfs

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestMemFS(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fs := NewMemFS()
		require.NoError(t, fstest.TestFS(fs))
		require.NoError(t, fs.(*memFS).CheckInvariants())
	})

	t.Run("basic", func(t *testing.T) {
		fs := NewMemFS()
		require.NoError(t, fs.Mkdir("foo"))
		require.NoError(t, fs.Mkdir("foo/bar"))
		require.NoError(t, fs.Mkdir("foo/bar/baz"))
		require.NoError(t, fstest.TestFS(fs, "foo", "foo/bar", "foo/bar/baz"))
		require.NoError(t, fs.(*memFS).CheckInvariants())
	})
}

func TestConsistency(t *testing.T) {
	tempDir = t.TempDir()
	rapid.Check(t, rapid.Run(&cmpFS{}))
}

var tempDir string
var count int

type cmpFS struct {
	dir dirFS
	mem *memFS
}

func (cmp *cmpFS) Init(t *rapid.T) {
	count++
	baseDir := fmt.Sprintf("%s/%d", tempDir, count)
	require.NoError(t, os.Mkdir(baseDir, 0755))
	cmp.dir = NewDirFS(baseDir).(dirFS)
	cmp.mem = NewMemFS().(*memFS)
}

func (cmp *cmpFS) Mkdir(t *rapid.T) {
	name := rapid.String().Draw(t, "name").(string)
	errl := cmp.dir.Mkdir(name)
	errr := cmp.mem.Mkdir(name)
	require.Equal(t, errl == nil, errr == nil)
}

func (cmp *cmpFS) Check(t *rapid.T) {
	require.NoError(t, cmp.mem.CheckInvariants())

	left, err := cmp.dir.ReadDir(".")
	require.NoError(t, err)
	right, err := cmp.mem.ReadDir(".")
	require.NoError(t, err)
	require.Equal(t, len(left), len(right))
	// TODO improve comparison
}
