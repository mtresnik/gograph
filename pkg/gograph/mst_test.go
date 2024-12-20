package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"image/png"
	"os"
	"testing"
	"time"
)

func TestKruskalMST(t *testing.T) {
	size := 35
	provider := BoundedRandomGraphProvider{
		BoundingBox:    gomath.BoundingBox{10, 10, 20, 20},
		NumPoints:      size * 5,
		NumConnections: 5,
	}
	graph := provider.Build()

	start := time.Now().UnixMilli()
	response := KruskalMST(MSTRequest{Graph: graph})
	end := time.Now().UnixMilli()
	println("KruskalMST Time:", end-start)
	newGraph := response.Graph

	renderer := NewGraphRenderer(size*40, size*40)
	//renderer.AddGraph(graph, color.RGBA{
	//	R: 255,
	//	G: 0,
	//	B: 0,
	//	A: 255,
	//})
	renderer.AddGraph(newGraph)
	img := renderer.Render()
	file, err := os.Create("TestKruskalMST.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func TestPrimsMST(t *testing.T) {
	size := 35
	provider := BoundedRandomGraphProvider{
		BoundingBox:    gomath.BoundingBox{10, 10, 20, 20},
		NumPoints:      size * 5,
		NumConnections: 5,
	}
	graph := provider.Build()
	start := time.Now().UnixMilli()
	response := PrimsMST(MSTRequest{Graph: graph})
	end := time.Now().UnixMilli()
	println("PrimsMST Time:", end-start)
	newGraph := response.Graph

	renderer := NewGraphRenderer(size*40, size*40)
	renderer.AddGraph(newGraph)
	img := renderer.Render()
	file, err := os.Create("TestPrimsMST.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
