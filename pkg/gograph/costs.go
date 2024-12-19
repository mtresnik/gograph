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
	Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64
}

func GenerateInitialCosts(pCostFunctions *map[string]CostFunction) (map[string]CostFunction, map[string]CostEntry) {
	var costFunctions map[string]CostFunction
	if pCostFunctions == nil {
		costFunctions = map[string]CostFunction{COST_TYPE_DISTANCE: EuclideanDistanceCostFunction{}}
	}
	initialCosts := map[string]CostEntry{}
	for key, _ := range costFunctions {
		initialCosts[key] = CostEntry{
			Accumulated: 0,
			Current:     0,
			Total:       0,
		}
	}
	return costFunctions, initialCosts
}

func GenerateNextCosts(currWrapper *VertexWrapper, toVertex Vertex, costFunctions map[string]CostFunction) map[string]CostEntry {
	if currWrapper == nil {
		return map[string]CostEntry{}
	}
	nextCosts := map[string]CostEntry{}
	for key, _ := range costFunctions {
		nextCostByKey := GetCostOrEvaluate(currWrapper, toVertex, key, costFunctions)
		nextCosts[key] = CostEntry{
			Accumulated: currWrapper.Costs[key].Total,
			Current:     nextCostByKey,
			Total:       currWrapper.Costs[key].Total + nextCostByKey,
		}
	}
	return nextCosts
}

func GenerateWorstCosts(costFunctions map[string]CostFunction, optionalMax ...float64) map[string]CostEntry {
	maxValue := math.MaxFloat64
	if len(optionalMax) > 0 {
		maxValue = optionalMax[0]
	}
	worstCosts := map[string]CostEntry{}
	for key, _ := range costFunctions {
		worstCosts[key] = CostEntry{
			Accumulated: maxValue,
			Current:     maxValue,
			Total:       maxValue,
		}
	}
	return worstCosts
}

func maxDifference(one, two float64) float64 {
	if one > two {
		return math.Abs(one - two)
	}
	return math.Abs(two - one)
}

func GenerateCostDifference(one, two map[string]CostEntry) map[string]CostEntry {
	delta := map[string]CostEntry{}
	for key, _ := range one {
		delta[key] = CostEntry{
			Accumulated: maxDifference(one[key].Accumulated, two[key].Accumulated),
			Current:     maxDifference(one[key].Current, two[key].Current),
			Total:       maxDifference(one[key].Total, two[key].Total),
		}
	}
	return delta
}

func MultiplyCosts(costs map[string]CostEntry) CostEntry {
	accumulated := 1.0
	current := 1.0
	total := 1.0
	for _, cost := range costs {
		accumulated *= cost.Accumulated
		current *= cost.Current
		total *= cost.Total
	}
	return CostEntry{
		Accumulated: accumulated,
		Current:     current,
		Total:       total,
	}
}

func AddCosts(costs map[string]CostEntry) CostEntry {
	accumulated := 0.0
	current := 0.0
	total := 0.0
	for _, cost := range costs {
		accumulated += cost.Accumulated
		current += cost.Current
		total += cost.Total
	}
	return CostEntry{
		Accumulated: accumulated,
		Current:     current,
		Total:       total,
	}
}

func GetCostOrEvaluate(currWrapper *VertexWrapper, toVertex Vertex, key string, costFunctions map[string]CostFunction) float64 {
	edge := GetEdge(currWrapper.Inner, toVertex)
	if edge == nil {
		return 0.0
	}
	costMap := edge.Cost()
	if costMap != nil {
		cost, ok := (*costMap)[key]
		if ok {
			return cost
		}
	}

	costFunction := costFunctions[key]
	return costFunction.Eval(currWrapper, toVertex)
}

type MultiCostFunction interface {
	GetCostFunctions() map[string]CostFunction
	Eval(vertexWrapper VertexWrapper, to gomath.Spatial) float64
}
type EuclideanDistanceCostFunction struct{}

func (f EuclideanDistanceCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.EuclideanDistance)
}

type HaversineDistanceCostFunction struct{}

func (f HaversineDistanceCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.HaversineDistance)
}

type ManhattanDistanceCostFunction struct{}

func (f ManhattanDistanceCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
	return gomath.ToPoint(vertexWrapper.Inner).DistanceTo(gomath.ToPoint(to), gomath.ManhattanDistance)
}

type AdditiveCostFunction struct {
	Functions map[string]CostFunction
}

func (f AdditiveCostFunction) GetCostFunctions() map[string]CostFunction {
	return f.Functions
}

func (f AdditiveCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
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

func (f *MultiplicativeCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
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

func (f PowerCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
	return math.Pow(f.Base.Eval(vertexWrapper, to), f.Exponent.Eval(vertexWrapper, to))
}

type AbsoluteValueCostFunction struct {
	Inner CostFunction
}

func (f AbsoluteValueCostFunction) Eval(vertexWrapper *VertexWrapper, to gomath.Spatial) float64 {
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
