package gomaze

import (
	"image"
	"image/color"
	"image/gif"
	"os"
	"testing"
)

func testMazeRenderer_Build(t *testing.T) {
	request := NewMazeGeneratorRequest(50, 50)
	img := MazeRenderer{
		CellSize:     20,
		VisitedColor: nil,
		Request:      request,
		Generator:    AldousBroderMazeGenerator{},
	}.Build()
	// Display the image in a new window
	images := []*image.Paletted{img}
	delays := []int{0}
	g := &gif.GIF{
		Image: images,
		Delay: delays,
	}
	f, err := os.Create("TestMazeRenderer_Build.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Encode and save the GIF
	err = gif.EncodeAll(f, g)
	if err != nil {
		panic(err)
	}
}

func testMazeRenderer_RenderFrames(t *testing.T) {
	request := NewMazeGeneratorRequest(50, 50)
	renderer := LiveMazeRenderer{
		CellSize:     20,
		VisitedColor: &color.RGBA{R: 111, G: 111, B: 255, A: 255},
		Request:      request,
		Generator:    AldousBroderMazeGenerator{},
		Frames:       []*image.Paletted{},
	}
	g := renderer.RenderFrames()
	f, err := os.Create("TestMazeRenderer_RenderFrames.gif")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	// Encode and save the GIF
	err = gif.EncodeAll(f, g)
	if err != nil {
		panic(err)
	}
}
