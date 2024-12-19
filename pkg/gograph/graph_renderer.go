package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"github.com/mtresnik/goutils/pkg/goutils"
	"image"
	"image/color"
	"image/gif"
)

type IGraphRenderer interface {
	Render() *image.RGBA
	AddPoint(point gomath.Spatial, color ...color.RGBA)
	AddGraph(graph Graph, color ...color.RGBA)
	AddPath(path Path, color ...color.RGBA)
	AddEdge(edge Edge, color ...color.RGBA)
}

type GraphRenderer struct {
	Graphs        map[int64]Graph
	Points        map[int64]gomath.Spatial
	Edges         map[int64]Edge
	Colors        map[int64]color.RGBA
	Paths         map[int64]Path
	bounds        gomath.BoundingBox
	Width         int
	Height        int
	Padding       int
	lineThickness int
	pointRadius   int
}

func NewGraphRenderer(width, height int) *GraphRenderer {
	return &GraphRenderer{
		Graphs:        make(map[int64]Graph),
		Points:        make(map[int64]gomath.Spatial),
		Colors:        make(map[int64]color.RGBA),
		Paths:         make(map[int64]Path),
		bounds:        gomath.BoundingBox{0, 0, 0, 0},
		Width:         width,
		Height:        height,
		Padding:       50,
		lineThickness: 5,
		pointRadius:   10,
	}
}

func (g *GraphRenderer) AddPoint(point gomath.Spatial, color ...color.RGBA) {
	hash := gomath.HashSpatial(point)
	g.Points[hash] = point
	if len(color) > 0 {
		g.Colors[hash] = color[0]
	}
	points := make([]gomath.Spatial, 0)
	if g.bounds.Area() > 0 {
		gPoints := g.bounds.GetPoints()
		for _, point := range gPoints {
			points = append(points, point)
		}
	}
	points = append(points, point)
	newBounds := gomath.NewBoundingBox(points...)
	g.bounds = newBounds
}

func (g *GraphRenderer) AddGraph(graph Graph, color ...color.RGBA) {
	hashOrId := GraphHashOrId(graph)
	g.Graphs[hashOrId] = graph
	if len(color) > 0 {
		g.Colors[hashOrId] = color[0]
	}
	points := make([]gomath.Spatial, 0)
	if g.bounds.Area() > 0 {
		gPoints := g.bounds.GetPoints()
		for _, point := range gPoints {
			points = append(points, point)
		}
	}
	vertices := graph.GetVertices()
	edges := graph.GetEdges()
	for _, v := range vertices {
		points = append(points, v)
	}
	for _, e := range edges {
		points = append(points, e.From())
		points = append(points, e.To())
	}
	newBounds := gomath.NewBoundingBox(points...)
	g.bounds = newBounds
}

func (g *GraphRenderer) AddPath(path Path, color ...color.RGBA) {
	hashOrId := PathHashOrId(path)
	g.Paths[hashOrId] = path
	if len(color) > 0 {
		g.Colors[hashOrId] = color[0]
	}
	points := make([]gomath.Spatial, 0)
	if g.bounds.Area() > 0 {
		gPoints := g.bounds.GetPoints()
		for _, point := range gPoints {
			points = append(points, point)
		}
	}
	for _, edge := range path.GetEdges() {
		points = append(points, edge.From())
		points = append(points, edge.To())
	}
	newBounds := gomath.NewBoundingBox(points...)
	g.bounds = newBounds
}

func (g *GraphRenderer) AddEdge(edge Edge, color ...color.RGBA) {
	hashOrId := EdgeHashOrId(edge)
	g.Edges[hashOrId] = edge
	if len(color) > 0 {
		g.Colors[hashOrId] = color[0]
	}
	points := make([]gomath.Spatial, 0)
	if g.bounds.Area() <= 0 {
		gPoints := g.bounds.GetPoints()
		for _, point := range gPoints {
			points = append(points, point)
		}
	}
	points = append(points, edge.From())
	points = append(points, edge.To())
	newBounds := gomath.NewBoundingBox(points...)
	g.bounds = newBounds
}

func (g *GraphRenderer) convertPixels(point gomath.Spatial) image.Point {
	if g.bounds.Area() <= 0 {
		panic("Bounds not set")
	}
	if !g.bounds.Contains(point) {
		return image.Point{}
	}
	x := point.GetValues()[0]
	y := point.GetValues()[1]

	relX := (x - g.bounds.MinX) / (g.bounds.MaxX - g.bounds.MinX)
	relY := (y - g.bounds.MinY) / (g.bounds.MaxY - g.bounds.MinY)
	return image.Point{
		X: int(relX * float64(g.Width)),
		Y: int(relY * float64(g.Height)),
	}
}

func (g *GraphRenderer) drawLine(img *image.RGBA, from, to gomath.Spatial, color color.RGBA) {
	fromCoords := g.convertPixels(from)
	toCoords := g.convertPixels(to)
	goutils.DrawLine(img, fromCoords.X, fromCoords.Y, toCoords.X, toCoords.Y, color, g.lineThickness)
}

func (g *GraphRenderer) drawPoint(img *image.RGBA, point gomath.Spatial, color color.RGBA) {
	coords := g.convertPixels(point)
	goutils.FillCircle(img, coords.X, coords.Y, g.pointRadius, color)
}

func (g *GraphRenderer) padBounds() {
	minPoint := gomath.Point{Values: []float64{g.bounds.MinX, g.bounds.MinY}}
	maxPoint := gomath.Point{Values: []float64{g.bounds.MaxX, g.bounds.MaxY}}
	dx := maxPoint.X() - minPoint.X()
	dy := maxPoint.Y() - minPoint.Y()

	dt := dx * (float64(g.Padding) / float64(g.Width))
	dv := dy * (float64(g.Padding) / float64(g.Height))

	minPoint = gomath.Point{Values: []float64{g.bounds.MinX - dt, g.bounds.MinY - dv}}
	maxPoint = gomath.Point{Values: []float64{g.bounds.MaxX + dt, g.bounds.MaxY + dv}}

	g.bounds.MinX = minPoint.X()
	g.bounds.MinY = minPoint.Y()
	g.bounds.MaxX = maxPoint.X()
	g.bounds.MaxY = maxPoint.Y()
}

func (g *GraphRenderer) Render() *image.RGBA {
	if g.bounds.Area() <= 0 {
		return nil
	}
	g.padBounds()
	img := image.NewRGBA(image.Rect(0, 0, g.Width, g.Height))
	goutils.FillRectangle(img, 0, 0, g.Width, g.Height, color.White)
	for hash, graph := range g.Graphs {
		graphColor, ok := g.Colors[hash]
		if !ok {
			graphColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
		}
		for _, edge := range graph.GetEdges() {
			g.drawLine(img, edge.From(), edge.To(), graphColor)
		}
		for _, vertex := range graph.GetVertices() {
			g.drawPoint(img, vertex, graphColor)
		}
	}
	for hash, point := range g.Points {
		pointColor, ok := g.Colors[hash]
		if !ok {
			pointColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
		}
		g.drawPoint(img, point, pointColor)
	}
	for hash, edge := range g.Edges {
		edgeColor, ok := g.Colors[hash]
		if !ok {
			edgeColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
		}
		g.drawLine(img, edge.From(), edge.To(), edgeColor)
	}

	for hash, path := range g.Paths {
		pathColor, ok := g.Colors[hash]
		if !ok {
			pathColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		}
		for _, edge := range path.GetEdges() {
			g.drawLine(img, edge.From(), edge.To(), pathColor)
			g.drawPoint(img, edge.From(), pathColor)
			g.drawPoint(img, edge.To(), pathColor)
		}
	}
	return img
}

type LiveGraphRenderer struct {
	Graph         Graph
	Request       RoutingAlgorithmRequest
	Frames        []*image.Paletted
	Width         int
	Height        int
	Padding       int
	lineThickness int
	pointRadius   int
}

func NewLiveGraphRenderer(graph Graph, request RoutingAlgorithmRequest, width, height int) *LiveGraphRenderer {
	return &LiveGraphRenderer{
		Graph:         graph,
		Request:       request,
		Frames:        make([]*image.Paletted, 0),
		Width:         width,
		Height:        height,
		Padding:       50,
		lineThickness: 5,
		pointRadius:   10,
	}
}

func (g *LiveGraphRenderer) Build() {
	g.Request.UpdateListeners = &[]RoutingAlgorithmUpdateListener{g}
	EvaluateRoutingAlgorithm(g.Request)
}

func (g *LiveGraphRenderer) Update(response RoutingAlgorithmResponse) {
	points := make([]gomath.Spatial, 0)
	pointColor := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	visitedColor := color.RGBA{R: 111, G: 111, B: 255, A: 255}
	points = append(points, g.Request.Start, g.Request.Destination)
	pointMap := make(map[int64]gomath.Spatial)
	colorMap := make(map[int64]color.RGBA)
	for _, point := range points {
		hash := gomath.HashSpatial(point)
		pointMap[hash] = point
		colorMap[hash] = pointColor
	}
	for _, vertex := range g.Graph.GetVertices() {
		hash := vertex.Hash()
		_, ok := response.Visited[hash]
		if ok {
			colorMap[hash] = visitedColor
			pointMap[hash] = vertex
		}
	}
	internalGraphRenderer := GraphRenderer{
		Graphs:        map[int64]Graph{},
		Points:        pointMap,
		Colors:        colorMap,
		Paths:         map[int64]Path{},
		bounds:        gomath.BoundingBox{},
		Width:         g.Width,
		Height:        g.Height,
		Padding:       g.Padding,
		lineThickness: g.lineThickness,
		pointRadius:   g.pointRadius,
	}
	internalGraphRenderer.AddGraph(g.Graph)
	internalGraphRenderer.AddPath(response.Path)
	internalGraphRenderer.AddPoint(g.Request.Start, pointColor)
	internalGraphRenderer.AddPoint(g.Request.Destination, pointColor)
	img := internalGraphRenderer.Render()
	paletted := goutils.ConvertImageToPaletted(img)
	g.Frames = append(g.Frames, paletted)
}

func (g *LiveGraphRenderer) RenderFrames() *gif.GIF {
	g.Build()
	if len(g.Frames) == 0 {
		return nil
	}
	images := g.Frames
	delays := make([]int, len(images))
	for i := 0; i < len(images); i++ {
		delays[i] = 5
	}
	retGif := &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: -1,
	}
	return retGif

}
