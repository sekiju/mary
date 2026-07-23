package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_UnmarshalsJSONFileConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "version": 1,
  "application": {
    "check_updates": false,
    "max_parallel_downloads": 7,
    "max_parallel_chapters": 3
  },
  "output": {
    "directory": "custom-dir"
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	Params.Runtime = RuntimeFlags{ListChaptersMode: true}

	if err := Load(path); err != nil {
		t.Fatal(err)
	}

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

func TestLoad_MissingFileGeneratesDefaults(t *testing.T) {
	defaults := FileConfig{
		Version: CurrentConfigVersion,
		Application: application{
			CheckUpdates:         true,
			MaxParallelDownloads: 4,
		},
		Output: output{
			Directory:  "downloads",
			FileFormat: AutoOutputFormat,
		},
	}
	Params.File = defaults

	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	if err := Load(path); err != nil {
		t.Fatal(err)
	}

	if Params.File.Application.MaxParallelDownloads != 4 {
		t.Errorf("MaxParallelDownloads = %d, want default 4", Params.File.Application.MaxParallelDownloads)
	}
	if Params.File.Output.Directory != "downloads" {
		t.Errorf("Output.Directory = %q, want default %q", Params.File.Output.Directory, "downloads")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Config file was not generated on missing-file load")
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	Params.File = FileConfig{
		Version: CurrentConfigVersion,
		Application: application{
			CheckUpdates:         false,
			MaxParallelDownloads: 8,
			MaxParallelChapters:  2,
		},
		Output: output{
			Directory:  "output-dir",
			FileFormat: PngOutputFormat,
		},
		Sites: map[string]Site{
			"example.com": {Cookie: ptr("secret")},
		},
	}

	if err := Save(path); err != nil {
		t.Fatal(err)
	}

	checkPerms(t, path)

	Params.File = FileConfig{}

	if err := Load(path); err != nil {
		t.Fatal(err)
	}

	if Params.File.Application.MaxParallelDownloads != 8 {
		t.Errorf("MaxParallelDownloads = %d, want 8", Params.File.Application.MaxParallelDownloads)
	}
	if Params.File.Output.FileFormat != PngOutputFormat {
		t.Errorf("FileFormat = %q, want %q", Params.File.Output.FileFormat, PngOutputFormat)
	}
	if Params.File.Sites["example.com"].Cookie == nil || *Params.File.Sites["example.com"].Cookie != "secret" {
		t.Error("Site cookie not round-tripped")
	}
	if Params.File.Schema == "" {
		t.Error("$schema field was not injected on save")
	}
	if Params.File.Version != CurrentConfigVersion {
		t.Errorf("Version = %d, want %d", Params.File.Version, CurrentConfigVersion)
	}
}

func TestResolvePath_Override(t *testing.T) {
	path, err := ResolvePath("/custom/path/config.json")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/custom/path/config.json" {
		t.Errorf("path = %q, want %q", path, "/custom/path/config.json")
	}
}

func TestResolvePath_Default(t *testing.T) {
	t.Setenv("MDL_CONFIG", "")
	path, err := ResolvePath("")
	if err != nil {
		t.Fatal(err)
	}
	def, err := DefaultConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if path != def {
		t.Errorf("path = %q, want %q", path, def)
	}
}

func TestResolvePath_EnvOverride(t *testing.T) {
	t.Setenv("MDL_CONFIG", "/env/path/config.json")
	path, err := ResolvePath("")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/env/path/config.json" {
		t.Errorf("path = %q, want /env/path/config.json", path)
	}
}

func ptr(s string) *string { return &s }
