package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"image"
	"image/color"
)

func HashColor(c color.Color) int64 {
	r, g, b, a := c.RGBA()
	values := []float64{float64(r), float64(g), float64(b), float64(a)}
	return gomath.HashSpatial(gomath.Point{Values: values})
}

func GetColors(img *image.RGBA) []color.Color {
	uniqueColors := make(map[int64]bool)
	retArray := make([]color.Color, 0)
	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dy(); y++ {
			col := img.At(x, y)
			hash := HashColor(col)
			if _, ok := uniqueColors[hash]; !ok {
				uniqueColors[hash] = true
				retArray = append(retArray, col)
			}
		}
	}
	return retArray
}

func ConvertImageToPaletted(img *image.RGBA, colors ...color.RGBA) *image.Paletted {
	palette := make([]color.Color, len(colors))
	if len(palette) > 0 {
		for i, c := range colors {
			palette[i] = c
		}
	} else {
		palette = GetColors(img)
	}
	retImage := image.NewPaletted(img.Rect, palette)
	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dy(); y++ {
			retImage.Set(x, y, img.At(x, y))
		}
	}
	return retImage
}
