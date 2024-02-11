package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func WriteImageBytes(filePath string, b []byte) error {
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	ext := filepath.Ext(filePath)
	if ext == "" {
		return fmt.Errorf("file has no extension: %s", filePath)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
