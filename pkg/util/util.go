package util

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func SaveImages(images []image.Image, path string) error {
	for i, img := range images {

		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
		output_path := fmt.Sprintf("%sframe%d.jpg", path, i)
		out, err := os.Create(output_path)
		if err != nil {
			return err
		}
		defer out.Close()

		err = jpeg.Encode(out, img, nil)
		if err != nil {
			return err
		}
		log.Printf("Frame %d saved", i)
	}
	return nil
}
