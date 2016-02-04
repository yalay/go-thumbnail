package util

import (
	"fmt"
	"github.com/oliamb/cutter"
	"image"
)

func CropImg(srcImg image.Image, dstWidth, dstHeight int) image.Image {
	dstImg, err := cutter.Crop(srcImg, cutter.Config{
		Height: dstHeight,       // height in pixel or Y ratio(see Ratio Option below)
		Width:  dstWidth,        // width in pixel or X ratio
		Mode:   cutter.Centered, // Accepted Mode: TopLeft, Centered
		//Anchor:  image.Point{100, 100}, // Position of the top left point
		Options: 0, // Accepted Option: Ratio
	})
	if err != nil {
		fmt.Printf("Cannot crop image:", err)
		return srcImg
	}
	return dstImg
}
