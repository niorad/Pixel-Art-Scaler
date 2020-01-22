package processing

import "image"

// BasicScaling Applies basic interpolation to an image.
// It's doubling the image-size in each dimension
func BasicScaling(srcImg image.Image) image.Image {

	oWidth := srcImg.Bounds().Max.X
	oHeight := srcImg.Bounds().Max.Y

	targetSizeRect := image.Rectangle{image.Point{0, 0}, image.Point{oWidth * 2, oHeight * 2}}
	newImg := image.NewRGBA(targetSizeRect)

	for y := 0; y < oHeight; y++ {
		for x := 0; x < oWidth; x++ {

			srcColor := srcImg.At(x, y)
			cursorX, cursorY := x*2, y*2

			newBottomRightColor := srcColor
			newBottomLeftColor := srcColor
			newTopLeftColor := srcColor
			newTopRightColor := srcColor

			if x < oWidth-2 && y < oHeight-2 {
				rightCol := srcImg.At(x+1, y)
				bottomCol := srcImg.At(x, y+1)
				bottomRightCol := srcImg.At(x+1, y+1)
				if rightCol == bottomCol && bottomCol == bottomRightCol {
					newBottomRightColor = rightCol
				}
			}

			if x > 0 && y < oHeight-2 {
				leftCol := srcImg.At(x-1, y)
				bottomCol := srcImg.At(x, y+1)
				bottomLeftCol := srcImg.At(x-1, y+1)
				if leftCol == bottomCol && bottomCol == bottomLeftCol {
					newBottomLeftColor = leftCol
				}
			}

			if x < oWidth-2 && y > 0 {
				rightCol := srcImg.At(x+1, y)
				topCol := srcImg.At(x, y-1)
				topRightCol := srcImg.At(x+1, y-1)
				if rightCol == topCol && topCol == topRightCol {
					newTopRightColor = rightCol
				}
			}

			if x > 0 && y > 0 {

				leftCol := srcImg.At(x-1, y)
				topCol := srcImg.At(x, y-1)
				topLeftCol := srcImg.At(x-1, y-1)
				if leftCol == topCol && topCol == topLeftCol {
					newTopLeftColor = leftCol
				}
			}

			newImg.Set(cursorX+1, cursorY+1, newBottomRightColor)
			newImg.Set(cursorX, cursorY+1, newBottomLeftColor)
			newImg.Set(cursorX+1, cursorY, newTopRightColor)
			newImg.Set(cursorX, cursorY, newTopLeftColor)

		}
	}

	return newImg

}
