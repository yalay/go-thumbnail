package util

import (
	"github.com/oliamb/cutter"
	"image"
)

func CropImg(srcImg image.Image, dstWidth, dstHeight int) image.Image {
	origBounds := srcImg.Bounds()
	origWidth := origBounds.Dx()
	origHeight := origBounds.Dy()

	dstImg, err := cutter.Crop(srcImg, cutter.Config{
		Height: dstHeight,      // height in pixel or Y ratio(see Ratio Option below)
		Width:  dstWidth,       // width in pixel or X ratio
		Mode:   cutter.TopLeft, // Accepted Mode: TopLeft, Centered
		Anchor: image.Point{
			origWidth / 12,
			origHeight / 8}, // Position of the top left point
		Options: 0, // Accepted Option: Ratio
	})
	if err != nil {
		Logln("[GIN] Cannot crop image:" + err.Error())
		return srcImg
	}
	return dstImg
}
