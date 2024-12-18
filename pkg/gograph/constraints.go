package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type Constraint interface {
	Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool
}

type NegationConstraint struct {
	Inner Constraint
}

func (b NegationConstraint) Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool {
	return !b.Inner.Check(currentVertex, nextCost)
}

type ForAllConstraint struct {
	Key         string
	Constraints []Constraint
}

func CheckAllConstraints(currentVertex *VertexWrapper, nextCost CostEntry, key string, constraintMap map[string][]Constraint) bool {
	cost := map[string]CostEntry{key: nextCost}
	constraints, ok := constraintMap[key]
	if !ok {
		return true
	}
	for _, constraint := range constraints {
		if !constraint.Check(currentVertex, cost) {
			return false
		}
	}
	return true
}

func (b ForAllConstraint) Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool {
	constraintMap := map[string][]Constraint{b.Key: b.Constraints}
	return CheckAllConstraints(currentVertex, nextCost[b.Key], b.Key, constraintMap)
}

type ForEachConstraint struct {
	Key         string
	Constraints []Constraint
}

func CheckAnyConstraints(currentVertex *VertexWrapper, nextCost CostEntry, key string, constraintMap map[string][]Constraint) bool {
	cost := map[string]CostEntry{key: nextCost}
	constraints, ok := constraintMap[key]
	if !ok {
		return true
	}
	for _, constraint := range constraints {
		if constraint.Check(currentVertex, cost) {
			return true
		}
	}
	return false
}

func (b ForEachConstraint) Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool {
	constraintMap := map[string][]Constraint{b.Key: b.Constraints}
	return CheckAnyConstraints(currentVertex, nextCost[b.Key], b.Key, constraintMap)
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

type MaximumCostConstraint struct {
	Key     string
	Maximum float64
}

func (b MaximumCostConstraint) Check(_ VertexWrapper, nextCost map[string]CostEntry) bool {
	return nextCost[b.Key].Total <= b.Maximum && nextCost[b.Key].Current <= b.Maximum
}
