package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"testing"
)

func TestSimpleEdge_String(t *testing.T) {
	edge := NewEdge(gomath.Point{Values: []float64{0.0, 0.0}}, gomath.Point{Values: []float64{1.0, 1.0}})
	println(edge.String())
}

func TestSimpleEdge_Distance(t *testing.T) {
	edge := NewEdge(gomath.Point{Values: []float64{0.0, 0.0}}, gomath.Point{Values: []float64{1.0, 1.0}})
	distance := edge.Distance(gomath.ManhattanDistance)
	println(distance)

	distance = edge.Distance()
	println(distance)
}

func TestSimpleEdge_Split(t *testing.T) {
	edge := NewEdge(gomath.Point{Values: []float64{0.0, 0.0}}, gomath.Point{Values: []float64{100.0, 100.0}})
	split := edge.Split(5)
	cast := CastToEdges(split...)
	println(EdgesToString(cast...))
}

func TestPolyEdge_String(t *testing.T) {
	polyEdge := NewEdge(gomath.Point{Values: []float64{0.0, 0.0}}, gomath.Point{Values: []float64{0.5, 1.0}}, gomath.Point{Values: []float64{2.0, 2.0}})
	println(polyEdge.String())
	println(polyEdge.Distance(gomath.ManhattanDistance))
}
