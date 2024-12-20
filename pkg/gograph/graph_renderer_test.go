package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"image/gif"
	"image/png"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestGraphRenderer_Render(t *testing.T) {
	response := AldousBroderMazeGenerator(NewMazeGeneratorRequest(25, 25))
	graph := MazeToGraphProvider{response.Maze}.Build()
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

func testLiveGraphRenderer_RenderFrames(t *testing.T) {
	size := 50
	randomPruneProvider := RandomPruneGraphProvider{BoundedGridGraphProvider{
		BoundingBox: gomath.BoundingBox{10, 10, 20, 20},
		Width:       size,
		Height:      size,
	}, 0.10}
	graph := randomPruneProvider.Build()
	vertices := graph.GetVertices()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	startVertex := vertices[random.Intn(len(vertices))]
	endVertex := vertices[random.Intn(len(vertices))]
	algorithm := AStar
	println("start:", gomath.SpatialString(startVertex), "\tend:", gomath.SpatialString(endVertex))
	println("distance:", algorithm(RoutingAlgorithmRequest{
		Start:       startVertex,
		Destination: endVertex,
		Algorithm:   algorithm,
	}).Path.Length())
	renderer := NewLiveGraphRenderer(graph, RoutingAlgorithmRequest{
		Start:       startVertex,
		Destination: endVertex,
		Algorithm:   algorithm,
	}, size*30, size*30)
	g := renderer.RenderFrames()
	f, err := os.Create("TestLiveGraphRenderer_RenderFrames.gif")
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

func TestBoundedGraphProvider_UI(t *testing.T) {
	provider := BoundedGridGraphProvider{
		BoundingBox: gomath.BoundingBox{10, 10, 20, 20},
		Width:       40,
		Height:      40,
	}
	graph := provider.Build()
	renderer := NewGraphRenderer(4000, 4000)
	renderer.AddGraph(graph)
	img := renderer.Render()
	file, err := os.Create("TestBoundedGraphProvider_UI.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func TestRandomGraphProvider_UI(t *testing.T) {
	randomPruneProvider := RandomPruneGraphProvider{BoundedGridGraphProvider{
		BoundingBox: gomath.BoundingBox{10, 10, 20, 20},
		Width:       40,
		Height:      40,
	}, 0.10}
	graph := randomPruneProvider.Build()
	renderer := NewGraphRenderer(4000, 4000)
	renderer.AddGraph(graph)
	img := renderer.Render()
	file, err := os.Create("TestRandomGraphProvider_UI.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
