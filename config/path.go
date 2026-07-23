package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func DefaultConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), ".config", "mdl"), nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "mdl"), nil
}

func DefaultConfigPath() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// ResolvePath returns the effective config path, handling flag and env overrides.
// If explicitPath is non-empty, it is used verbatim (absolute or relative).
// Otherwise MDL_CONFIG env var is checked, then the OS-specific default.
func ResolvePath(explicitPath string) (string, error) {
	if explicitPath != "" {
		return explicitPath, nil
	}
	if env := os.Getenv("MDL_CONFIG"); env != "" {
		return env, nil
	}
	return DefaultConfigPath()
}
