package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type RoutingConstraint interface {
	Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool
}

type NegationRoutingConstraint struct {
	Inner RoutingConstraint
}

func (b NegationRoutingConstraint) Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool {
	return !b.Inner.Check(currentVertex, nextCost)
}

type ForAllRoutingConstraint struct {
	Key         string
	Constraints []RoutingConstraint
}

func CheckAllRoutingConstraints(currentVertex *VertexWrapper, nextCost CostEntry, key string, constraintMap map[string][]RoutingConstraint) bool {
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

func (b ForAllRoutingConstraint) Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool {
	constraintMap := map[string][]RoutingConstraint{b.Key: b.Constraints}
	return CheckAllRoutingConstraints(currentVertex, nextCost[b.Key], b.Key, constraintMap)
}

type ForEachRoutingConstraint struct {
	Key         string
	Constraints []RoutingConstraint
}

func CheckAnyRoutingConstraints(currentVertex *VertexWrapper, nextCost CostEntry, key string, constraintMap map[string][]RoutingConstraint) bool {
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

func (b ForEachRoutingConstraint) Check(currentVertex *VertexWrapper, nextCost map[string]CostEntry) bool {
	constraintMap := map[string][]RoutingConstraint{b.Key: b.Constraints}
	return CheckAnyRoutingConstraints(currentVertex, nextCost[b.Key], b.Key, constraintMap)
}

type ShapeContainsRoutingConstraint struct {
	Shape            gomath.Shape
	DistanceFunction *gomath.DistanceFunction
}

func (b ShapeContainsRoutingConstraint) Check(currentVertex VertexWrapper, _ map[string]CostEntry) bool {
	var distanceFunction gomath.DistanceFunction
	if b.DistanceFunction != nil {
		distanceFunction = *b.DistanceFunction
	} else {
		distanceFunction = gomath.EuclideanDistance
	}
	return b.Shape.Contains(gomath.ToPoint(currentVertex.Inner), distanceFunction)
}

type MaximumCostRoutingConstraint struct {
	Key     string
	Maximum float64
}

func (b MaximumCostRoutingConstraint) Check(_ VertexWrapper, nextCost map[string]CostEntry) bool {
	return nextCost[b.Key].Total <= b.Maximum && nextCost[b.Key].Current <= b.Maximum
}
