package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"math"
	"strings"
)

const (
	COST_TYPE_DISTANCE = "distance"
	COST_TYPE_TIME     = "time"
)

type CostFunction interface {
	GetType() string
	Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64
}

type MultiCostFunction interface {
	GetTypes() []string
	GetCostFunctions() []CostFunction
	Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64
}
type EuclideanDistanceCostFunction struct{}

func (f EuclideanDistanceCostFunction) GetType() string {
	return COST_TYPE_DISTANCE
}

func (f EuclideanDistanceCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.EuclideanDistance{})
}

type HaversineDistanceCostFunction struct{}

func (f HaversineDistanceCostFunction) GetType() string {
	return COST_TYPE_DISTANCE
}

func (f HaversineDistanceCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.HaversineDistance{})
}

type ManhattanDistanceCostFunction struct{}

func (f ManhattanDistanceCostFunction) GetType() string {
	return COST_TYPE_DISTANCE
}

func (f ManhattanDistanceCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.ManhattanDistance{})
}

type AdditiveCostFunction struct {
	Functions []CostFunction
}

func (f AdditiveCostFunction) GetTypes() []string {
	retArray := make([]string, len(f.Functions))
	for i, function := range f.Functions {
		retArray[i] = function.GetType()
	}
	return retArray
}

func (f AdditiveCostFunction) GetCostFunctions() []CostFunction {
	return f.Functions
}

func (f AdditiveCostFunction) GetType() string {
	retArray := make([]string, len(f.Functions))
	for i, function := range f.Functions {
		retArray[i] = function.GetType()
	}
	return strings.Join(retArray, "+")
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

func (f MultiplicativeCostFunction) GetTypes() []string {
	retArray := make([]string, len(f.Functions))
	for i, function := range f.Functions {
		retArray[i] = function.GetType()
	}
	return retArray
}

func (f MultiplicativeCostFunction) GetCostFunctions() []CostFunction {
	return f.Functions
}

func (f MultiplicativeCostFunction) GetType() string {
	retArray := make([]string, len(f.Functions))
	for i, function := range f.Functions {
		retArray[i] = function.GetType()
	}
	return strings.Join(retArray, "*")
}

func (f MultiplicativeCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	totalCost := 1.0
	for _, function := range f.Functions {
		totalCost *= function.Eval(vertexWrapper, to)
	}
	return totalCost
}

type ConstantCostFunction struct {
	Type  string
	Value float64
}

func (f ConstantCostFunction) GetType() string {
	return f.Type
}

func (f ConstantCostFunction) Eval(_, _ gomath.Spatial) float64 {
	return f.Value
}

type PowerCostFunction struct {
	Base     CostFunction
	Exponent CostFunction
}

func (f PowerCostFunction) GetType() string {
	return f.Base.GetType() + "^" + f.Exponent.GetType()
}

func (f PowerCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return math.Pow(f.Base.Eval(vertexWrapper, to), f.Exponent.Eval(vertexWrapper, to))
}

type AbsoluteValueCostFunction struct {
	Inner CostFunction
}

func (f AbsoluteValueCostFunction) GetType() string {
	return f.Inner.GetType()
}

func (f AbsoluteValueCostFunction) Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64 {
	return math.Abs(f.Inner.Eval(vertexWrapper, to))
}
