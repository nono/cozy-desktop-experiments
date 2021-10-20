package localfs

import (
	"fmt"
	"os"
	"path/filepath"
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
	dir     dirFS
	mem     *memFS
	parents []string
}

func (cmp *cmpFS) Init(t *rapid.T) {
	count++
	baseDir := fmt.Sprintf("%s/%d", tempDir, count)
	require.NoError(t, os.Mkdir(baseDir, 0750))
	cmp.dir = NewDirFS(baseDir).(dirFS)
	cmp.mem = NewMemFS().(*memFS)
	cmp.parents = []string{"."}
}

func (cmp *cmpFS) Mkdir(t *rapid.T) {
	parent := rapid.SampledFrom(cmp.parents).Draw(t, "parent").(string)
	name := rapid.String().Draw(t, "name").(string)
	path := filepath.Join(parent, name)
	errl := cmp.dir.Mkdir(path)
	errr := cmp.mem.Mkdir(path)
	require.Equal(t, errl == nil, errr == nil)
	if errl == nil {
		cmp.parents = append(cmp.parents, path)
	}
}

func (cmp *cmpFS) Remove(t *rapid.T) {
	path := rapid.SampledFrom(cmp.parents).Draw(t, "path").(string)
	errl := cmp.dir.RemoveAll(path)
	errr := cmp.mem.RemoveAll(path)
	require.Equal(t, errl == nil, errr == nil)
}

func (cmp *cmpFS) Check(t *rapid.T) {
	require.NoError(t, cmp.mem.CheckInvariants())

	left, err := cmp.dir.ToMemFS()
	require.NoError(t, err)
	right := cmp.mem
	require.Equal(t, len(left.ByPath), len(right.ByPath))
	for k, v := range left.ByPath {
		require.Contains(t, right.ByPath, k)
		require.Equal(t, v.path, right.ByPath[k].path)
		require.Equal(t, v.info.name, right.ByPath[k].info.name)
		require.Equal(t, v.info.mode, right.ByPath[k].info.mode)
	}
}
