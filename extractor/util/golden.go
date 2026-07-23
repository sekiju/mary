package util

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var updateGolden = flag.Bool("update", false, "update golden test files")

// AssertGolden compares actual against the contents of the golden file at
// path. Run `go test -update ./...` to (re)write the golden file from the
// current actual value instead of asserting.
func AssertGolden(t *testing.T, path string, actual []byte) {
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
			return
		}
		if err := os.WriteFile(path, actual, 0o644); err != nil {
			t.Fatal(err)
			return
		}
		return
	}

	expected, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, expected, actual)
}
