//go:build !windows

package config

import (
	"os"
	"testing"
)

func checkPerms(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("file perms = %o, want 0600", info.Mode().Perm())
	}
}
