package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_UnmarshalsOnlyFileConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.hcl")
	content := `
application {
  check_updates = false
  max_parallel_downloads = 7
  max_parallel_chapters = 3
}
output {
  directory = "custom-dir"
}
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	Params.Runtime = RuntimeFlags{ListChaptersMode: true}

	Load(path)

	if Params.File.Application.MaxParallelDownloads != 7 {
		t.Errorf("MaxParallelDownloads = %d, want 7", Params.File.Application.MaxParallelDownloads)
	}
	if Params.File.Application.MaxParallelChapters != 3 {
		t.Errorf("MaxParallelChapters = %d, want 3", Params.File.Application.MaxParallelChapters)
	}
	if Params.File.Output.Directory != "custom-dir" {
		t.Errorf("Output.Directory = %q, want %q", Params.File.Output.Directory, "custom-dir")
	}
	if !Params.Runtime.ListChaptersMode {
		t.Error("Runtime.ListChaptersMode was reset by Load, but Load must not touch RuntimeFlags")
	}
}

func TestLoad_MissingFileKeepsDefaults(t *testing.T) {
	Params.File = FileConfig{
		Application: application{MaxParallelDownloads: 4},
		Output:      output{Directory: "downloads", FileFormat: AutoOutputFormat},
	}

	Load(filepath.Join(t.TempDir(), "does-not-exist.hcl"))

	if Params.File.Application.MaxParallelDownloads != 4 {
		t.Errorf("MaxParallelDownloads = %d, want default 4", Params.File.Application.MaxParallelDownloads)
	}
	if Params.File.Output.Directory != "downloads" {
		t.Errorf("Output.Directory = %q, want default %q", Params.File.Output.Directory, "downloads")
	}
}
