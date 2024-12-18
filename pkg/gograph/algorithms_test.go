package gograph

import (
	"image/png"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestBFS_Evaluate(t *testing.T) {
	mazeResponse := AldousBroderMazeGenerator{}.Build(NewMazeGeneratorRequest(25, 25))
	graph := MazeToGraphProvider{}.Build(mazeResponse.Maze)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	vertices := graph.GetVertices()
	startVertex := vertices[random.Intn(len(vertices))]
	endVertex := vertices[random.Intn(len(vertices))]
	bfs := BFS{}
	response := bfs.Evaluate(RoutingAlgorithmRequest{
		Start:       startVertex,
		Destination: endVertex,
	})
	renderer := NewGraphRenderer(1250, 1250)
	renderer.AddGraph(graph)
	renderer.AddPath(response.Path)
	img := renderer.Render()
	file, err := os.Create("TestBFS_Evaluate.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
