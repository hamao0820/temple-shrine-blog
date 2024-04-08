package util

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
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
	fileName := fmt.Sprintf("%s.webp", uid.String())
	p := filepath.Join(dir, fileName)

	file, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = webp.Encode(file, img, nil)
	if err != nil {
		return "", err
	}

	return p, nil
}

// func resizeImage(img image.Image, width, height int) image.Image {
// 	// 欲しいサイズの画像を新しく作る
// 	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

// 	// サイズを変更しながら画像をコピーする
// 	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)

// 	return newImage
// }

// func resizeImageKeepAspect(img image.Image, size float64) image.Image {
// 	// 画像のサイズを取得する
// 	width := float64(img.Bounds().Max.X)
// 	height := float64(img.Bounds().Max.Y)

// 	// 結果となる画像のサイズを計算する
// 	if width > height {
// 		height = height * size / width
// 		width = size
// 	} else {
// 		width = width * size / height
// 		height = size
// 	}

// 	// 先ほどの関数を使って画像をリサイズする
// 	return resizeImage(img, int(width), int(height))
// }
