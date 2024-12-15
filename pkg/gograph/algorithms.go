package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type Algorithm interface {
	GetInitialCosts() map[string]CostEntry
	GetStart() gomath.Spatial
	GetDestination() gomath.Spatial
	IdSet() map[int64]bool
	GetCostFunction() CostFunction
	GetConstraint() Constraint
}
