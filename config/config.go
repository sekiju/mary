package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

var Params = config{
	File: FileConfig{
		Version: CurrentConfigVersion,
		Settings: settings{
			CheckForUpdates:             true,
			AutoInstallUpdates:          true,
			MaxParallelPageFetches:      4,
			MaxParallelChaptersDownload: 1,
			OutputDirectory:             "downloads",
		},
	},
}

func Load(explicitPath string) error {
	path, err := ResolvePath(explicitPath)
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Info().Str("path", path).Msg("Config doesn't exist, generating defaults")
		return Save(path)
	} else if err != nil {
		return fmt.Errorf("stat config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return &InvalidConfigError{Path: path, Err: fmt.Errorf("parse config: %w", err)}
	}

	changed, err := RunMigrations(raw)
	if err != nil {
		return fmt.Errorf("migrate config: %w", err)
	}

	if changed {
		migrated, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal migrated config: %w", err)
		}
		var fc FileConfig
		if err := json.Unmarshal(migrated, &fc); err != nil {
			return fmt.Errorf("unmarshal migrated config: %w", err)
		}
		Params.File = fc
		Params.File.Schema = schemaURL
		return Save(path)
	}

	var fc FileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return &InvalidConfigError{Path: path, Err: fmt.Errorf("unmarshal config: %w", err)}
	}
	Params.File = fc
	return nil
}

// InvalidConfigError indicates the config file at Path exists but could not be
// parsed. Callers may offer to overwrite it with defaults via Save(Path).
type InvalidConfigError struct {
	Path string
	Err  error
}

func (e *InvalidConfigError) Error() string {
	return fmt.Sprintf("invalid config file %q: %s", e.Path, e.Err)
}

func (e *InvalidConfigError) Unwrap() error {
	return e.Err
}

// SetSiteCookie sets (or, if cookie is nil, clears) the stored cookie for
// hostname in the in-memory config. Callers must call Save to persist it.
func SetSiteCookie(hostname string, cookie *string) {
	if Params.File.Sites == nil {
		Params.File.Sites = make(map[string]Site)
	}
	if cookie == nil {
		delete(Params.File.Sites, hostname)
		return
	}
	Params.File.Sites[hostname] = Site{Cookie: cookie}
}

func Save(path string) error {
	if path == "" {
		resolved, err := DefaultConfigPath()
		if err != nil {
			return fmt.Errorf("resolve default config path: %w", err)
		}
		path = resolved
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	Params.File.Schema = schemaURL
	Params.File.Version = CurrentConfigVersion

	data, err := json.MarshalIndent(Params.File, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
