package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type Constraint interface {
	Check(currentVertex VertexWrapper, nextCost []CostEntry) bool
}

type NegationConstraint struct {
	Inner Constraint
}

func (b NegationConstraint) Check(currentVertex VertexWrapper, nextCost []CostEntry) bool {
	return !b.Inner.Check(currentVertex, nextCost)
}

type ForAllConstraint struct {
	Constraints []Constraint
}

func CheckAllConstraints(currentVertex VertexWrapper, nextCost []CostEntry, constraints ...Constraint) bool {
	for _, constraint := range constraints {
		if !constraint.Check(currentVertex, nextCost) {
			return false
		}
	}
	return true
}

func (b ForAllConstraint) Check(currentVertex VertexWrapper, nextCost []CostEntry) bool {
	return CheckAllConstraints(currentVertex, nextCost, b.Constraints...)
}

type ForEachConstraint struct {
	Constraints []Constraint
}

func CheckAnyConstraint(currentVertex VertexWrapper, nextCost []CostEntry, constraints ...Constraint) bool {
	for _, constraint := range constraints {
		if constraint.Check(currentVertex, nextCost) {
			return true
		}
	}
	return false
}

func (b ForEachConstraint) Check(currentVertex VertexWrapper, nextCost []CostEntry) bool {
	return CheckAnyConstraint(currentVertex, nextCost, b.Constraints...)
}

type ShapeContainsConstraint struct {
	Shape            gomath.Shape
	DistanceFunction *gomath.DistanceFunction
}

func (b ShapeContainsConstraint) Check(currentVertex VertexWrapper, _ []CostEntry) bool {
	var distanceFunction gomath.DistanceFunction
	if b.DistanceFunction != nil {
		distanceFunction = *b.DistanceFunction
	} else {
		distanceFunction = gomath.EuclideanDistance{}
	}
	return b.Shape.Contains(gomath.ToPoint(currentVertex.Inner), distanceFunction)
}

type MaximumCostConstraint struct {
	Index   int
	Maximum float64
}

func (b MaximumCostConstraint) Check(_ VertexWrapper, nextCost []CostEntry) bool {
	return nextCost[b.Index].Total <= b.Maximum && nextCost[b.Index].Current <= b.Maximum
}
