package util

// 修改自github.com/nfnt/resize/thumbnail.go
// 参数中最大长宽修改为最小长宽
import (
	"github.com/nfnt/resize"
	"image"
)

// 缩略图按照指定的最小长宽不失真缩放
func Thumbnail(minWidth, minHeight uint, img image.Image) image.Image {
	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	newWidth, newHeight := origWidth, origHeight

	// Return original image if it have same or smaller size as constraints
	if minWidth >= origWidth && minHeight >= origHeight {
		return img
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
	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}
