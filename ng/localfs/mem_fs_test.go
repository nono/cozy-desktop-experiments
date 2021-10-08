package localfs

import (
	"testing"
	"testing/fstest"
)

func TestMemFS(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fs := MemFS()
		if err := fstest.TestFS(fs); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("basic", func(t *testing.T) {
		fs := MemFS()
		if err := fs.Mkdir("foo"); err != nil {
			t.Fatal(err)
		}
		if err := fs.Mkdir("foo/bar"); err != nil {
			t.Fatal(err)
		}
		if err := fs.Mkdir("foo/bar/baz"); err != nil {
			t.Fatal(err)
		}
		if err := fstest.TestFS(fs, "./foo", "foo/bar", "foo/bar/baz"); err != nil {
			t.Fatal(err)
		}
	})
}
