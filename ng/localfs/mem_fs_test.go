package localfs

import (
	"testing"
	"testing/fstest"
)

func TestMemFS(t *testing.T) {
	fs := MemFS()
	if err := fstest.TestFS(fs); err != nil {
		t.Fatal(err)
	}
}
