package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogging() (*os.File, error) {
	dir, err := DefaultLogDir()
	if err != nil {
		return nil, fmt.Errorf("resolve log dir: %w", err)
	}

	if err = os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.jsonl", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	log.Logger = zerolog.New(file).With().Timestamp().Logger()

	return file, nil
}
