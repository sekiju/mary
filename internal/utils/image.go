package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
)

func ReadImageBytes(filePath string) ([]byte, error) {
	imageBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading image file: %s", err)
	}

	return imageBytes, nil
}

func WriteImageBytes(filePath string, b []byte) error {
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	if ext := filepath.Ext(filePath); ext == "" {
		return fmt.Errorf("file has no extension: %s", filePath)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func WriteImage(filePath string, img image.Image) error {
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
		return fmt.Errorf("file without extension: %s", filePath)
	}

	switch ext {
	case ".jpg", ".jpeg":
		err := jpeg.Encode(file, img, nil)
		if err != nil {
			return fmt.Errorf("failed to jpeg encode: %s", err)
		}
	case ".png":
		err := png.Encode(file, img)
		if err != nil {
			return fmt.Errorf("failed to png encode: %s", err)
		}
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
