package gograph

import "github.com/mtresnik/gomath/pkg/gomath"

type Algorithm interface {
	GetStart() gomath.Spatial
	GetDestination() gomath.Spatial
	IdSet() map[int64]bool
	GetCostFunction() CostFunction
}
