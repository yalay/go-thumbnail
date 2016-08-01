package util

// 修改自github.com/nfnt/resize/thumbnail.go
// 参数中最大长宽修改为最小长宽
import (
	"github.com/nfnt/resize"
	"image"
)

// 缩略图按照指定的宽和高非失真缩放裁剪
func ThumbnailCrop(minWidth, minHeight uint, img image.Image) image.Image {
	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	// Return original image if it have same or smaller size as constraints
	if minWidth >= origWidth && minHeight >= origHeight {
		return img
	}

	if minWidth > origWidth {
		minWidth = origWidth
	}

	if minHeight > origHeight {
		minHeight = origHeight
	}

	// Preserve aspect ratio
	if origWidth > minWidth {
		newHeight = uint(origHeight * minWidth / origWidth)
		if newHeight < 1 {
			newHeight = 1
		}
		newWidth = minWidth
	}

	if newHeight < minHeight {
		newWidth = uint(newWidth * minHeight / newHeight)
		if newWidth < 1 {
			newWidth = 1
		}
		newHeight = minHeight
	}

	thumbImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return CropImg(thumbImg, int(minWidth), int(minHeight))
}

// 简单的缩放,指定最大宽和高
func ThumbnailSimple(maxWidth, maxHeight uint, img image.Image) image.Image {
	oriBounds := img.Bounds()
	oriWidth := uint(oriBounds.Dx())
	oriHeight := uint(oriBounds.Dy())

	if maxWidth == 0 {
		maxWidth = oriWidth
	}

	if maxHeight == 0 {
		maxHeight = oriHeight
	}
	return resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)
}
