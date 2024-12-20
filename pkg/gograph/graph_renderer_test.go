package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"github.com/mtresnik/goutils/pkg/goutils"
	"image"
	"image/gif"
	"image/png"
	"math"
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

type CustomCostFunction struct {
	BoundingBox    gomath.BoundingBox
	TimeScalar     float64
	PositionScalar float64
	ResultScalar   float64
	SplitNumber    int
}

func (c CustomCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
	edge := GetEdge(vertexWrapper.Inner, ToVertex(to))
	cost := edge.Cost()
	dx := c.BoundingBox.MaxX - c.BoundingBox.MinX
	dy := c.BoundingBox.MaxY - c.BoundingBox.MinY
	if cost == nil {
		cost = &map[string]float64{}
		(*cost)[COST_TYPE_TIME] = max(gomath.EuclideanDistance(vertexWrapper.Inner, to), 11)
	}
	currTime, ok := (*cost)[COST_TYPE_TIME]
	sum := 0.0
	if ok {
		segments := edge.Split(c.SplitNumber)
		for i, segment := range segments {
			p1 := gomath.NewPoint(segment.From().GetValues()...)
			p2 := gomath.NewPoint(segment.To().GetValues()...)
			midX := ((p1.X()+p2.X())/2.0 - c.BoundingBox.MinX) / dx
			midY := ((p1.Y()+p2.Y())/2.0 - c.BoundingBox.MinY) / dy
			x := c.PositionScalar * midX
			y := c.PositionScalar * midY
			z := c.TimeScalar * (vertexWrapper.Costs[COST_TYPE_TIME].Total + currTime*float64(i)/float64(len(segments)))
			sum += c.ResultScalar * max(gomath.PerlinNoise(x, y, z), 0.0)
		}
		sum = sum / float64(len(segments))
		sum += 1.0
		sum = math.Pow(sum, 2)
	}
	return max(1.0, sum)
}

func testAnimatedGraphProvider_UI(t *testing.T) {
	size := 25
	width := size * 30
	height := size * 30
	resultScalar := 2.0
	timeScalar := 0.01
	positionScalar := 2.0
	explorationFactor := 0.5
	splitNumber := 100
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	initialCosts := map[string]gomath.DistanceFunction{}
	initialCosts[COST_TYPE_DISTANCE] = gomath.EuclideanDistance
	initialCosts[COST_TYPE_TIME] = func(one, other gomath.Spatial) float64 {
		return rand.Float64()*5.0 + 1.0
	}
	boundingBox := gomath.BoundingBox{10, 10, 20, 20}
	// dx := boundingBox.MaxX - boundingBox.MinX
	// dy := boundingBox.MaxY - boundingBox.MinY

	randomPruneProvider := BoundedRandomGraphProvider{
		BoundingBox:    boundingBox,
		NumPoints:      size * size,
		NumConnections: 4,
		CostFunctions:  &initialCosts,
	}
	graph := randomPruneProvider.Build()
	//randomPruneProvider := BoundedGridGraphProvider{
	//	BoundingBox: gomath.BoundingBox{10, 10, 20, 20},
	//	Width:       size,
	//	Height:      size,
	//}
	//graph := randomPruneProvider.Build()
	vertices := graph.GetVertices()
	topLeft := gomath.NewPoint(boundingBox.MinX, boundingBox.MinY+2.5)
	topRight := gomath.NewPoint(boundingBox.MaxX, boundingBox.MinY+2.5)
	startVertex := vertices[random.Intn(len(vertices))]
	startVertex = goutils.MinBy(vertices, func(vertex Vertex) float64 {
		return gomath.EuclideanDistance(topLeft, vertex)
	})
	endVertex := vertices[random.Intn(len(vertices))]
	endVertex = goutils.MinBy(vertices, func(vertex Vertex) float64 {
		return gomath.EuclideanDistance(topRight, vertex)
	})
	algorithm := AStar

	var costFunctions = &map[string]CostFunction{}
	(*costFunctions)[COST_TYPE_DISTANCE] = EuclideanDistanceCostFunction{}
	(*costFunctions)[COST_TYPE_TIME] = InitialCostFunction{
		Default: 1.0,
		Type:    COST_TYPE_TIME,
	}

	(*costFunctions)["custom"] = CustomCostFunction{
		BoundingBox:    boundingBox,
		TimeScalar:     timeScalar,
		PositionScalar: positionScalar,
		ResultScalar:   resultScalar,
		SplitNumber:    splitNumber,
	}

	var backgroundGenerator func(time float64) *image.RGBA = func(time float64) *image.RGBA {
		retImage := image.NewRGBA(image.Rect(0, 0, width, height))
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				xPrime := positionScalar * float64(x+1) / float64(width)
				yPrime := positionScalar * float64(y+1) / float64(height)
				z := time * timeScalar
				inner := resultScalar * gomath.PerlinNoise(xPrime, yPrime, z)
				retImage.Set(x, y, goutils.GradientGreenToRed(inner))
			}
		}
		return retImage
	}

	println("start:", gomath.SpatialString(startVertex), "\tend:", gomath.SpatialString(endVertex))
	println("distance:", algorithm(RoutingAlgorithmRequest{
		Start:             startVertex,
		Destination:       endVertex,
		Algorithm:         algorithm,
		CostFunctions:     costFunctions,
		ExplorationFactor: explorationFactor,
	}).Path.Length())
	renderer := NewLiveGraphRenderer(graph, RoutingAlgorithmRequest{
		Start:             startVertex,
		Destination:       endVertex,
		Algorithm:         algorithm,
		CostFunctions:     costFunctions,
		ExplorationFactor: explorationFactor,
	}, width, height, &backgroundGenerator)
	g := renderer.RenderFrames()
	f, err := os.Create("TestAnimatedGraphProvider_UI.gif")
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
