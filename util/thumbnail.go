package util

// 修改自github.com/nfnt/resize/thumbnail.go
// 参数中最大长宽修改为最小长宽
import (
	"image"

	"github.com/disintegration/imaging"
)

// 缩略图按照指定的宽和高非失真缩放裁剪
func ThumbnailCrop(minWidth, minHeight int, srcImage image.Image) image.Image {
	return imaging.Fill(srcImage, minWidth, minHeight, imaging.TopLeft, imaging.Lanczos)
}

// 简单的缩放,指定最大宽和高
func ThumbnailSimple(maxWidth, maxHeight int, srcImage image.Image) image.Image {
	return imaging.Fit(srcImage, maxWidth, maxHeight, imaging.Lanczos)
}
