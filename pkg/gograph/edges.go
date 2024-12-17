package gograph

import (
	"fmt"
	"github.com/mtresnik/gomath/pkg/gomath"
	"hash/fnv"
	"math"
	"strings"
)

// Edge <editor-fold>
type Edge interface {
	From() gomath.Spatial
	To() gomath.Spatial
	Reverse() Edge
	Distance(distanceFunction ...gomath.DistanceFunction) float64
	DistanceCached(distanceFunction ...gomath.DistanceFunction) float64
	Scale(t float64, distanceFunction ...gomath.DistanceFunction) gomath.Spatial
	Split(size int, distanceFunction ...gomath.DistanceFunction) []gomath.Segment
	Id() int64
	String() string
	Hash() int64
}

func EdgeHashOrId(e Edge) int64 {
	if e.Id() != -1 {
		return e.Id()
	}
	return e.Hash()
}

type EdgeListener interface {
	Visit(e Edge)
}

func EdgesToString(edges ...Edge) string {
	retArray := make([]string, len(edges))
	for i, edge := range edges {
		retArray[i] = edge.String()
	}
	return "[" + strings.Join(retArray, ",") + "]"
}

func New3dEdgeFromFloat64(floats ...float64) Edge {
	numPoints := len(floats) / 3
	if numPoints < 2 {
		panic("New3dEdgeFromFloat64 needs at least two points, or six values")
	}
	points := make([]gomath.Spatial, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = gomath.Point{Values: []float64{floats[3*i], floats[3*i+1], floats[3*i+2]}}
	}
	return NewEdge(points...)
}

func New2dEdgeFromFloat64(floats ...float64) Edge {
	numPoints := len(floats) / 2
	if numPoints < 2 {
		panic("New2dEdgeFromFloat64 needs at least two points, or four values")
	}
	points := make([]gomath.Spatial, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = gomath.Point{Values: []float64{floats[2*i], floats[2*i+1]}}
	}
	return NewEdge(points...)
}

func NewEdge(points ...gomath.Spatial) Edge {
	if len(points) < 2 {
		panic("At least two points are required to create an edge")
	}
	if len(points) == 2 {
		return SimpleEdge{points[0], points[1], -1, -1, -1}
	}
	return PolyEdge{points, -1, -1, -1}
}

func CastToEdges(segments ...gomath.Segment) []Edge {
	var edges = make([]Edge, 0)
	for _, segment := range segments {
		edge, ok := segment.(Edge)
		if ok {
			edges = append(edges, edge)
		}
	}
	return edges
}

func Contract(firstEdge gomath.Segment, edges ...gomath.Segment) gomath.Segment {
	if len(edges) == 0 {
		return firstEdge
	}
	// put all points into one array and then return a PolyEdge
	allPoints := []gomath.Spatial{firstEdge.From(), firstEdge.To()}
	for _, edge := range edges {
		allPoints = append(allPoints, edge.To())
	}
	id := int64(-1)
	asEdge, ok := firstEdge.(Edge)
	if ok {
		id = asEdge.Id()
	}
	return PolyEdge{Points: allPoints, id: id}
}

// </editor-fold>

// SimpleEdge <editor-fold>
type SimpleEdge struct {
	from     gomath.Spatial
	to       gomath.Spatial
	id       int64
	distance float64
	hash     int64
}

func NewSimpleEdge(from gomath.Spatial, to gomath.Spatial, id int64) SimpleEdge {
	return SimpleEdge{from, to, id, -1.0, -1}
}

func (e SimpleEdge) From() gomath.Spatial {
	return e.from
}

func (e SimpleEdge) To() gomath.Spatial {
	return e.to
}

func (e SimpleEdge) Reverse() Edge {
	return SimpleEdge{e.to, e.from, e.id, -1.0, -1}
}

func (e SimpleEdge) Id() int64 {
	return e.id
}

func (e SimpleEdge) Distance(distanceFunction ...gomath.DistanceFunction) float64 {
	return gomath.ToPoint(e.from).DistanceTo(gomath.ToPoint(e.to), distanceFunction...)
}

func (e SimpleEdge) DistanceCached(distanceFunction ...gomath.DistanceFunction) float64 {
	if e.distance > 0 {
		return e.distance
	}
	e.distance = e.Distance(distanceFunction...)
	return e.distance
}

func (e SimpleEdge) Scale(t float64, _ ...gomath.DistanceFunction) gomath.Spatial {
	p1, p2 := gomath.ToPoint(e.from), gomath.ToPoint(e.to)
	v1 := p2.Subtract(p1)
	v2 := v1.Scale(t)
	return p1.AddVector(v2)
}

func (e SimpleEdge) MidPoint() gomath.Spatial {
	return e.Scale(0.5)
}

func (e SimpleEdge) Split(size int, distanceFunction ...gomath.DistanceFunction) []gomath.Segment {
	if size <= 0 {
		return []gomath.Segment{}
	}
	if size == 1 {
		return []gomath.Segment{e}
	}
	delta := 1.0 / float64(size)
	retArray := make([]gomath.Segment, size)
	previous := e.From()
	for i := 0; i < size; i++ {
		scalar := float64(i+1) * delta
		current := e.Scale(scalar, distanceFunction...)
		retArray[i] = SimpleEdge{previous, current, -1, -1, -1}
		previous = current
	}
	return retArray
}

func (e SimpleEdge) String() string {
	if e.id == -1 {
		return fmt.Sprintf("[%v -> %v]", gomath.ToPoint(e.From()).String(), gomath.ToPoint(e.To()).String())
	}
	return fmt.Sprintf("[%v -> %v]:%v", gomath.ToPoint(e.From()).String(), gomath.ToPoint(e.To()).String(), e.Id())
}

func (e SimpleEdge) ToPolyEdge(size int) PolyEdge {
	if size <= 0 {
		panic("At least one edge is required to create a poly edge")
	}
	if size == 1 {
		return PolyEdge{[]gomath.Spatial{e.From(), e.To()}, e.Id(), e.distance, -1}
	}
	split := e.Split(size) // use Euclidean distance
	contracted := Contract(split[0], split[1:]...)
	polyEdge, ok := contracted.(PolyEdge)
	if ok {
		return polyEdge
	}
	// Shouldn't get here.
	return PolyEdge{[]gomath.Spatial{e.From(), e.To()}, e.Id(), e.distance, e.Hash()}
}

func (e SimpleEdge) Hash() int64 {
	if e.hash != -1 {
		return e.hash
	}
	fromHash := VertexHashOrId(ToVertex(e.From()))
	toHash := VertexHashOrId(ToVertex(e.To()))
	values := []float64{float64(fromHash), float64(toHash)}
	hasher := fnv.New64a()

	for _, value := range values {
		bits := math.Float64bits(value)
		buf := make([]byte, 8)
		for i := 0; i < 8; i++ {
			buf[i] = byte(bits >> (i * 8))
		}
		_, _ = hasher.Write(buf)
	}

	e.hash = int64(hasher.Sum64())
	return e.hash
}

// </editor-fold>

// PolyEdge <editor-fold>
type PolyEdge struct {
	Points   []gomath.Spatial
	id       int64
	distance float64
	hash     int64
}

func (e PolyEdge) From() gomath.Spatial {
	return e.Points[0]
}

func (e PolyEdge) To() gomath.Spatial {
	return e.Points[len(e.Points)-1]
}

func (e PolyEdge) Reverse() Edge {
	reversedVertices := make([]gomath.Spatial, len(e.Points))
	for i := 0; i < len(e.Points); i++ {
		reversedVertices[i] = e.Points[len(e.Points)-1-i]
	}
	return PolyEdge{Points: reversedVertices, id: e.id}
}

func (e PolyEdge) Id() int64 {
	return e.id
}

func (e PolyEdge) Distance(distanceFunction ...gomath.DistanceFunction) float64 {
	retSum := 0.0
	for i := 0; i < len(e.Points)-1; i++ {
		curr, next := gomath.ToPoint(e.Points[i]), gomath.ToPoint(e.Points[i+1])
		retSum += curr.DistanceTo(next, distanceFunction...)
	}
	return retSum
}

func (e PolyEdge) DistanceCached(distanceFunction ...gomath.DistanceFunction) float64 {
	if e.distance > 0 {
		return e.distance
	}
	e.distance = e.Distance(distanceFunction...)
	return e.distance
}

func (e PolyEdge) getRelativeDistances(distanceFunction ...gomath.DistanceFunction) []float64 {
	totalDistance := e.Distance(distanceFunction...)
	if totalDistance == 0 {
		return []float64{}
	}
	relativeDistances := make([]float64, len(e.Points))
	relativeDistances[0] = 0.0
	accumulatedDistance := 0.0
	for i := 0; i < len(e.Points)-1; i++ {
		curr, next := gomath.ToPoint(e.Points[i]), gomath.ToPoint(e.Points[i+1])
		currDistance := curr.DistanceTo(next, distanceFunction...)
		accumulatedDistance += currDistance
		relativeDistances[i+1] = accumulatedDistance / totalDistance
	}
	return relativeDistances
}

func (e PolyEdge) scaleInternal(t float64, relativeDistances []float64) gomath.Spatial {
	if len(relativeDistances) == 0 {
		return e.From()
	}
	for i := 0; i < len(relativeDistances)-1; i++ {
		u, v := relativeDistances[i], relativeDistances[i+1]
		if t > u && t < v {
			scale := math.Abs((t - u) / (v - u))
			tempEdge := SimpleEdge{
				from: e.Points[i],
				to:   e.Points[i+1],
				id:   -1,
			}
			return tempEdge.Scale(scale)
		}
	}
	// Illegal state exception
	panic(fmt.Sprintf("Illegal state: cannot scale to %f", t))
}

func (e PolyEdge) Scale(t float64, distanceFunction ...gomath.DistanceFunction) gomath.Spatial {
	if t <= 0.0 {
		return e.From()
	}
	if t >= 1.0 {
		return e.To()
	}
	relativeDistances := e.getRelativeDistances(distanceFunction...)
	return e.scaleInternal(t, relativeDistances)
}

func (e PolyEdge) Split(size int, distanceFunction ...gomath.DistanceFunction) []gomath.Segment {
	if size <= 0 {
		return []gomath.Segment{}
	}
	if size == 1 {
		return []gomath.Segment{e}
	}
	// Short circuit for when size is the same
	if size == len(e.Points)-1 {
		retArray := make([]gomath.Segment, size)
		for i := 0; i < len(e.Points)-1; i++ {
			curr, next := e.Points[i], e.Points[i+1]
			retArray[i] = SimpleEdge{curr, next, -1, -1, -1}
		}
		return retArray
	}
	relativeDistances := e.getRelativeDistances(distanceFunction...)
	delta := 1.0 / float64(size)
	retArray := make([]gomath.Segment, size)
	previous := e.From()
	for i := 1; i < size; i++ {
		scalar := float64(i) * delta
		current := e.scaleInternal(scalar, relativeDistances)
		retArray[i] = SimpleEdge{previous, current, -1, -1, -1}
		previous = current
	}
	retArray[0] = SimpleEdge{e.From(), retArray[1].From(), -1, -1, -1}
	return retArray
}

func (e PolyEdge) String() string {
	retArray := make([]string, len(e.Points))
	for i, point := range e.Points {
		retArray[i] = fmt.Sprintf("%v", gomath.ToPoint(point).String())
	}
	if e.id == -1 {
		return fmt.Sprintf("[%v]", strings.Join(retArray, "->"))
	}
	return fmt.Sprintf("[%v]:%v", strings.Join(retArray, "->"), e.id)
}

func (e PolyEdge) Hash() int64 {
	if e.hash != -1 {
		return e.hash
	}
	hash := int64(31)
	for _, point := range e.Points {
		hash ^= VertexHashOrId(ToVertex(point))
	}
	e.hash = hash
	return e.hash
}

// </editor-fold>
