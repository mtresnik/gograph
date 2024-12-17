package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"math"
)

const (
	COST_TYPE_DISTANCE = "distance"
	COST_TYPE_TIME     = "time"
)

type CostFunction interface {
	Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64
}

type MultiCostFunction interface {
	GetCostFunctions() map[string]CostFunction
	Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64
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

func (f ManhattanDistanceCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.ManhattanDistance{})
}

type AdditiveCostFunction struct {
	Functions map[string]CostFunction
}

func (f AdditiveCostFunction) GetCostFunctions() map[string]CostFunction {
	return f.Functions
}

func (f AdditiveCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	totalCost := 0.0
	for _, function := range f.Functions {
		totalCost += function.Eval(vertexWrapper, to)
	}
	return totalCost
}

type MultiplicativeCostFunction struct {
	Functions map[string]CostFunction
}

func (f *MultiplicativeCostFunction) GetCostFunctions() map[string]CostFunction {
	return f.Functions
}

func (f *MultiplicativeCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
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

type CostCombiner interface {
	Calculate(costs map[string]CostEntry) float64
}

type SumCostCombiner struct{}

func (c SumCostCombiner) Calculate(costs map[string]CostEntry) float64 {
	total := 0.0
	for _, cost := range costs {
		total += cost.Total
	}
	return total
}

type MinCostCombiner struct{}

func (c MinCostCombiner) Calculate(costs map[string]CostEntry) float64 {
	minCost := math.MaxFloat64
	for _, cost := range costs {
		if cost.Total < minCost {
			minCost = cost.Total
		}
	}
	return minCost
}

type MaxCostCombiner struct{}

func (c MaxCostCombiner) Calculate(costs map[string]CostEntry) float64 {
	maxCost := 0.0
	for _, cost := range costs {
		if cost.Total > maxCost {
			maxCost = cost.Total
		}
	}
	return maxCost
}

type MultiplicativeCostCombiner struct{}

func (c MultiplicativeCostCombiner) Calculate(costs map[string]CostEntry) float64 {
	total := 1.0
	for _, cost := range costs {
		total *= cost.Total
	}
	return total
}
