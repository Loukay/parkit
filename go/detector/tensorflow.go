package main

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	tf "github.com/galeone/tensorflow/tensorflow/go"
)

func createTensor(image io.ReadCloser) (*tf.Tensor, error) {
	srcImage, _ := jpeg.Decode(image)
	img := imaging.Fill(srcImage, 128, 128, imaging.Center, imaging.Lanczos)
	return imageToTensor(img, 128, 128)
}

func imageToTensor(img image.Image, imageHeight, imageWidth int) (tfTensor *tf.Tensor, err error) {
	var tfImage [1][128][128][3]float32

	for i := 0; i < imageWidth; i++ {
		for j := 0; j < imageHeight; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			tfImage[0][j][i][0] = float32(r) / float32(255)
			tfImage[0][j][i][1] = float32(g) / float32(255)
			tfImage[0][j][i][2] = float32(b) / float32(255)
		}
	}
	return tf.NewTensor(tfImage)
}
