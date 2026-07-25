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
  "settings": {
    "checkForUpdates": false,
    "maxParallelPageFetches": 7,
    "maxParallelChaptersDownload": 3,
    "outputDirectory": "custom-dir"
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	Params.Runtime = RuntimeFlags{ListChaptersMode: true}

	if err := Load(path); err != nil {
		t.Fatal(err)
	}

	if Params.File.Settings.MaxParallelPageFetches != 7 {
		t.Errorf("MaxParallelPageFetches = %d, want 7", Params.File.Settings.MaxParallelPageFetches)
	}
	if Params.File.Settings.MaxParallelChaptersDownload != 3 {
		t.Errorf("MaxParallelChaptersDownload = %d, want 3", Params.File.Settings.MaxParallelChaptersDownload)
	}
	if Params.File.Settings.OutputDirectory != "custom-dir" {
		t.Errorf("OutputDirectory = %q, want %q", Params.File.Settings.OutputDirectory, "custom-dir")
	}
	if !Params.Runtime.ListChaptersMode {
		t.Error("Runtime.ListChaptersMode was reset by Load, but Load must not touch RuntimeFlags")
	}
}

func TestLoad_MissingFileGeneratesDefaults(t *testing.T) {
	defaults := FileConfig{
		Version: CurrentConfigVersion,
		Settings: settings{
			CheckForUpdates:        true,
			MaxParallelPageFetches: 4,
			OutputDirectory:        "downloads",
		},
	}
	Params.File = defaults

	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	if err := Load(path); err != nil {
		t.Fatal(err)
	}

	if Params.File.Settings.MaxParallelPageFetches != 4 {
		t.Errorf("MaxParallelPageFetches = %d, want default 4", Params.File.Settings.MaxParallelPageFetches)
	}
	if Params.File.Settings.OutputDirectory != "downloads" {
		t.Errorf("OutputDirectory = %q, want default %q", Params.File.Settings.OutputDirectory, "downloads")
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
		Settings: settings{
			CheckForUpdates:             false,
			MaxParallelPageFetches:      8,
			MaxParallelChaptersDownload: 2,
			OutputDirectory:             "output-dir",
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

	if Params.File.Settings.MaxParallelPageFetches != 8 {
		t.Errorf("MaxParallelPageFetches = %d, want 8", Params.File.Settings.MaxParallelPageFetches)
	}
	if Params.File.Settings.OutputDirectory != "output-dir" {
		t.Errorf("OutputDirectory = %q, want %q", Params.File.Settings.OutputDirectory, "output-dir")
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
