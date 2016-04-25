package util

import (
	"image"
	"image/draw"
)

// 加水印
func WaterMark(srcImg image.Image) (image.Image, error) {
	// 读取水印图片
	markImg, err := LoadImage(WaterMarkImg)
	if err != nil {
		return nil, err
	}

	//把水印写到右下角，并向0坐标各偏移10个像素
	srcBounds := srcImg.Bounds()
	offset := image.Pt(srcBounds.Dx()-markImg.Bounds().Dx()-10, srcBounds.Dy()-markImg.Bounds().Dy()-10)
	newImg := image.NewNRGBA(srcBounds)
	draw.Draw(newImg, srcBounds, srcImg, image.ZP, draw.Src)
	draw.Draw(newImg, markImg.Bounds().Add(offset), markImg, image.ZP, draw.Over)
	return newImg, nil
}
