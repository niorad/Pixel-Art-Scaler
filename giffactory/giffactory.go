package giffactory

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"pixelartscaler/processing"
)

// Generate creates an animated GIF out of multiple images
func Generate(srcImg image.Image, scalingIterationCount int, frameDelay int, frameCount int) gif.GIF {

	var frames []*image.Paletted
	anim := gif.GIF{}

	for i := 0; i < frameCount; i++ {

		processedImage := srcImg

		for j := 0; j < scalingIterationCount; j++ {
			processedImage = processing.BasicScaling(processedImage, true)
		}

		mygif := image.NewPaletted(processedImage.Bounds().Bounds(), palette.WebSafe)
		draw.FloydSteinberg.Draw(mygif, processedImage.Bounds().Bounds(), processedImage, image.ZP)

		frames = append(frames, mygif)
		anim.Delay = append(anim.Delay, frameDelay)
	}

	anim.Image = frames
	anim.LoopCount = 0

	return anim

}
