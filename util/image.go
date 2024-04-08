package util

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	dir = "images"
)

func SaveImage(img image.Image) (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s.jpeg", uid.String())
	p := filepath.Join(dir, fileName)

	file, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return "", err
	}

	return p, nil
}
