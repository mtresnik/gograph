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

func testLiveGraphRenderer_RenderFrames(t *testing.T) {
	response := AldousBroderMazeGenerator{}.Build(NewMazeGeneratorRequest(25, 25))
	graph := MazeToGraphProvider{}.Build(response.Maze)
	vertices := graph.GetVertices()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	startVertex := vertices[random.Intn(len(vertices))]
	endVertex := vertices[random.Intn(len(vertices))]
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
	}, 1250, 1250)
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
