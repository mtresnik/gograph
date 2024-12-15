package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type Constraint interface {
	Check(currentVertex VertexWrapper, nextCost map[string]CostEntry) bool
}

type NegationConstraint struct {
	Inner Constraint
}

func (b NegationConstraint) Check(currentVertex VertexWrapper, nextCost map[string]CostEntry) bool {
	return !b.Inner.Check(currentVertex, nextCost)
}

type ForAllConstraint struct {
	Constraints []Constraint
}

func (b ForAllConstraint) Check(currentVertex VertexWrapper, nextCost map[string]CostEntry) bool {
	for _, constraint := range b.Constraints {
		if !constraint.Check(currentVertex, nextCost) {
			return false
		}
	}
	return true
}

type ForEachConstraint struct {
	Constraints []Constraint
}

func (b ForEachConstraint) Check(currentVertex VertexWrapper, nextCost map[string]CostEntry) bool {
	for _, constraint := range b.Constraints {
		if constraint.Check(currentVertex, nextCost) {
			return true
		}
	}
	return false
}

type ShapeContainsConstraint struct {
	Shape            gomath.Shape
	DistanceFunction *gomath.DistanceFunction
}

func (b ShapeContainsConstraint) Check(currentVertex VertexWrapper, _ map[string]CostEntry) bool {
	var distanceFunction gomath.DistanceFunction
	if b.DistanceFunction != nil {
		distanceFunction = *b.DistanceFunction
	} else {
		distanceFunction = gomath.EuclideanDistance{}
	}
	return b.Shape.Contains(gomath.ToPoint(currentVertex.Inner), distanceFunction)
}

type MaxCostConstraint struct {
	Maximum float64
	Type    string
}

func (b MaxCostConstraint) Check(_ VertexWrapper, nextCost map[string]CostEntry) bool {
	cost, ok := nextCost[b.Type]
	if !ok {
		return false
	}
	return cost.Current <= b.Maximum && cost.Accumulated <= b.Maximum
}
