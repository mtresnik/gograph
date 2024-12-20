package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"image/png"
	"os"
	"testing"
	"time"
)

func TestGreedyTSP_Eval(t *testing.T) {
	size := 30
	provider := BoundedRandomGraphProvider{
		BoundingBox:    gomath.BoundingBox{10, 10, 20, 20},
		NumPoints:      size * 5,
		NumConnections: 5,
	}
	graph := provider.Build()

	start := time.Now().UnixMilli()
	response := GreedyTSP(TSPRequest{Graph: graph})
	end := time.Now().UnixMilli()
	println("GreedyTSP Time:", end-start)

	renderer := NewGraphRenderer(max(size*40, 500), max(size*40, 500))
	renderer.AddGraph(graph)
	renderer.AddPath(response.Path)
	img := renderer.Render()
	file, err := os.Create("TestGreedyTSP_Eval.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func TestRandomTSP_Eval(t *testing.T) {
	size := 30
	provider := BoundedRandomGraphProvider{
		BoundingBox:    gomath.BoundingBox{10, 10, 20, 20},
		NumPoints:      size * 5,
		NumConnections: 5,
	}
	graph := provider.Build()

	start := time.Now().UnixMilli()
	response := RandomTSP(TSPRequest{Graph: graph})
	end := time.Now().UnixMilli()
	println("RandomTSP Time:", end-start)

	renderer := NewGraphRenderer(max(size*40, 500), max(size*40, 500))
	renderer.AddGraph(graph)
	renderer.AddPath(response.Path)
	img := renderer.Render()
	file, err := os.Create("TestRandomTSP_Eval.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func TestRepeatTSP_Eval(t *testing.T) {
	size := 30
	provider := BoundedRandomGraphProvider{
		BoundingBox:    gomath.BoundingBox{10, 10, 20, 20},
		NumPoints:      size * 5,
		NumConnections: 5,
	}
	graph := provider.Build()

	start := time.Now().UnixMilli()
	var tsp TSP = RandomTSP
	response :=
		RepeatTSP{
			tsp,
			100,
		}.Eval(TSPRequest{Graph: graph})
	end := time.Now().UnixMilli()
	println("RepeatTSP Time:", end-start)

	renderer := NewGraphRenderer(max(size*40, 500), max(size*40, 500))
	renderer.AddGraph(graph)
	renderer.AddPath(response.Path)
	img := renderer.Render()
	file, err := os.Create("TestRepeatTSP_Eval.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
