package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"math"
)

type CostFunction interface {
	Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64
}

func AddCost(one, other CostFunction) CostFunction {
	return AdditiveCostFunction{Functions: []CostFunction{one, other}}
}

func MultiplyCost(one, other CostFunction) CostFunction {
	return MultiplicativeCostFunction{Functions: []CostFunction{one, other}}
}

func PowerCost(one, other CostFunction) CostFunction {
	return PowerCostFunction{Base: one, Exponent: other}
}

type EuclideanDistanceCostFunction struct{}

func (f EuclideanDistanceCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.EuclideanDistance{})
}

type HaversineDistanceCostFunction struct{}

func (f HaversineDistanceCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.HaversineDistance{})
}

type ManhattanDistanceCostFunction struct{}

func (f ManhattanDistanceCostFunction) Eval(_ map[string]CostEntry, from, to gomath.Spatial) float64 {
	return gomath.ToPoint(from).DistanceTo(gomath.ToPoint(to), gomath.ManhattanDistance{})
}

type AdditiveCostFunction struct {
	Functions []CostFunction
}

func (f AdditiveCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	totalCost := 0.0
	for _, function := range f.Functions {
		totalCost += function.Eval(vertexWrapper, to)
	}
	return totalCost
}

type MultiplicativeCostFunction struct {
	Functions []CostFunction
}

func (f MultiplicativeCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	totalCost := 1.0
	for _, function := range f.Functions {
		totalCost *= function.Eval(vertexWrapper, to)
	}
	return totalCost
}

type ConstantCostFunction struct {
	Value float64
}

func (f ConstantCostFunction) Eval(_, _ gomath.Spatial) float64 {
	return f.Value
}

type PowerCostFunction struct {
	Base     CostFunction
	Exponent CostFunction
}

func (f PowerCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return math.Pow(f.Base.Eval(vertexWrapper, to), f.Exponent.Eval(vertexWrapper, to))
}

type AbsoluteValueCostFunction struct {
	Inner CostFunction
}

func (f AbsoluteValueCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return math.Abs(f.Inner.Eval(vertexWrapper, to))
}
