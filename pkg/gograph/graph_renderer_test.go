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
	response := AldousBroderMazeGenerator{}.Build(NewMazeGeneratorRequest(25, 25))
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
	size := 20
	randomPruneProvider := RandomPruneGraphProvider{BoundedGraphProvider{
		BoundingBox: gomath.BoundingBox{10, 10, 20, 20},
		Width:       size,
		Height:      size,
	}, 0.10}
	graph := randomPruneProvider.Build()
	vertices := graph.GetVertices()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	startVertex := vertices[random.Intn(len(vertices))]
	endVertex := vertices[random.Intn(len(vertices))]
	distanceCheck := gomath.ManhattanDistance{}
	for distanceCheck.Eval(startVertex, endVertex) < float64(size/2) {
		startVertex = vertices[random.Intn(len(vertices))]
		endVertex = vertices[random.Intn(len(vertices))]
	}
	bfs := BFS{}
	println("start:", gomath.SpatialString(startVertex), "\tend:", gomath.SpatialString(endVertex))
	println("distance:", bfs.Evaluate(RoutingAlgorithmRequest{
		Start:         startVertex,
		Destination:   endVertex,
		Algorithm:     &bfs,
		CostFunctions: &map[string]CostFunction{COST_TYPE_DISTANCE: ManhattanDistanceCostFunction{}},
	}).Path.Length())
	renderer := NewLiveGraphRenderer(graph, RoutingAlgorithmRequest{
		Start:       startVertex,
		Destination: endVertex,
		Algorithm:   &bfs,
	}, 1000, 1000)
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
	provider := BoundedGraphProvider{
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
	randomPruneProvider := RandomPruneGraphProvider{BoundedGraphProvider{
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
