package util

import (
	"github.com/nfnt/resize"
	"image"
)

// 简单的缩放,指定最大宽和高
func ThumbnailSimple(maxWidth, maxHeight uint, img image.Image) image.Image {
	return resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
}
