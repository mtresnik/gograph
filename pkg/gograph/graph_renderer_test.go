package gograph

import (
	"image/png"
	"os"
	"testing"
)

func TestGraphRenderer_Render(t *testing.T) {
	response := AldousBroderMazeGenerator{}.Build(NewMazeGeneratorRequest(25, 25))
	graph := MazeToGraphProvider{}.Build(response.Maze)
	renderer := NewGraphRenderer(1250, 1250)
	renderer.AddGraph(graph)
	img := renderer.Render()
	file, err := os.Create("TestGraphRenderer_Render.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
