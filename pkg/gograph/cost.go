package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type CostFunction interface {
	Eval(from, to gomath.Spatial) float64
}

func Add(one, other CostFunction) CostFunction {
	return AdditiveCostFunction{Functions: []CostFunction{one, other}}
}

func Multiply(one, other CostFunction) CostFunction {
	return MultiplicativeCostFunction{Functions: []CostFunction{one, other}}
}

type EuclideanDistanceCostFunction struct{}

func (f EuclideanDistanceCostFunction) Eval(from, to gomath.Spatial) float64 {
	return gomath.ToPoint(from).DistanceTo(gomath.ToPoint(to), gomath.EuclideanDistance{})
}

type HaversineDistanceCostFunction struct{}

func (f HaversineDistanceCostFunction) Eval(from, to gomath.Spatial) float64 {
	return gomath.ToPoint(from).DistanceTo(gomath.ToPoint(to), gomath.HaversineDistance{})
}

type ManhattanDistanceCostFunction struct{}

func (f ManhattanDistanceCostFunction) Eval(from, to gomath.Spatial) float64 {
	return gomath.ToPoint(from).DistanceTo(gomath.ToPoint(to), gomath.ManhattanDistance{})
}

type AdditiveCostFunction struct {
	Functions []CostFunction
}

func (f AdditiveCostFunction) Eval(from, to gomath.Spatial) float64 {
	totalCost := 0.0
	for _, function := range f.Functions {
		totalCost += function.Eval(from, to)
	}
	return totalCost
}

type MultiplicativeCostFunction struct {
	Functions []CostFunction
}

func (f MultiplicativeCostFunction) Eval(from, to gomath.Spatial) float64 {
	totalCost := 1.0
	for _, function := range f.Functions {
		totalCost *= function.Eval(from, to)
	}
	return totalCost
}
