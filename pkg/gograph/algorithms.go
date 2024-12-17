package gograph

import (
	"errors"
	"github.com/mtresnik/goutils/pkg/goutils"
	"maps"
)

type RoutingRequest struct {
	Start             Vertex
	Destination       Vertex
	Constraints       *map[string][]Constraint
	MultiCostFunction *MultiCostFunction
	CostCombiner      *CostCombiner
	CostFunctions     *map[string]CostFunction
}

type RoutingResponse struct {
	Costs   map[string]CostEntry
	Path    Path
	Visited map[int64]bool
}

type RoutingAlgorithm interface {
	VisitVertex(vertex Vertex)
	VisitEdge(edge Edge)
	Evaluate(parameters RoutingRequest) (RoutingResponse, *error)
}

type BFS struct {
	VertexListeners []VertexListener
	EdgeListeners   []EdgeListener
}

func (b *BFS) VisitVertex(vertex Vertex) {
	for _, visitor := range b.VertexListeners {
		visitor.Visit(vertex)
	}
}

func (b *BFS) VisitEdge(edge Edge) {
	for _, visitor := range b.EdgeListeners {
		visitor.Visit(edge)
	}
}

func Backtrack(vertex VertexWrapper) []Edge {
	var path []Edge
	var currentWrapper = vertex
	for currentWrapper.Previous != nil && currentWrapper.Hash() != currentWrapper.Previous.Hash() {
		path = append(path, currentWrapper.Previous.Inner.GetEdge(currentWrapper.Inner))
		currentWrapper = *currentWrapper.Previous
	}
	return path
}

func (b *BFS) Evaluate(parameters RoutingRequest) (RoutingResponse, *error) {
	start := parameters.Start
	destination := parameters.Destination
	constraints := parameters.Constraints
	var costFunctions map[string]CostFunction
	if parameters.CostFunctions != nil {
		costFunctions = *parameters.CostFunctions
	} else {
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
	startWrapper := NewVertexWrapper(start, initialCosts)
	startWrapper.Previous = nil
	queue := []VertexWrapper{startWrapper}
	visited := make(map[int64]bool)
	var curr VertexWrapper
	for len(queue) > 0 {
		curr = queue[0]
		queue = queue[1:]
		visited[VertexHashOrId(curr)] = true
		if curr.Hash() == destination.Hash() {
			println("found path")
			break
		}
		for _, edge := range curr.Inner.GetEdges() {
			toVertex := ToVertex(edge.To())
			hashOrId := VertexHashOrId(toVertex)
			nextCosts := map[string]CostEntry{}
			for key, costFunction := range costFunctions {
				nextCostByKey := costFunction.Eval(curr, toVertex)
				nextCosts[key] = CostEntry{
					Accumulated: curr.Costs[key].Total,
					Current:     nextCostByKey,
					Total:       curr.Costs[key].Total + nextCostByKey,
				}
			}
			if !goutils.SetContains(visited, hashOrId) {
				successor := NewVertexWrapper(toVertex, nextCosts)
				pass := true
				if constraints != nil {
					keys := maps.Keys(*constraints)
					for key := range keys {
						_, constraintsExist := (*constraints)[key]
						currCost, costExist := nextCosts[key]
						if constraintsExist && costExist {
							pass = CheckAllConstraints(curr, currCost, key, *constraints)
							if !pass {
								break
							}
						}
					}
				}
				if pass {
					b.VisitEdge(edge)
					previousWrapper := curr
					successor.Previous = &previousWrapper
					queue = append(queue, successor)
					visited[hashOrId] = true
					b.VisitVertex(toVertex)
				}
			}
		}
		if len(queue) == 0 {
			println("no path found, try relaxing the constraints")
			err := errors.New("no path found, try relaxing the constraints")
			return RoutingResponse{}, &err
		}
	}
	println("backtracking")
	edges := Backtrack(curr)
	println("edge count:", len(edges))

	costCombinerPtr := parameters.CostCombiner
	var costCombiner CostCombiner
	if costCombinerPtr == nil {
		costCombiner = MultiplicativeCostCombiner{}
	} else {
		costCombiner = *costCombinerPtr
	}

	path := NewSimplePath(edges, costCombiner.Calculate(curr.Costs))
	return RoutingResponse{
		Costs:   curr.Costs,
		Path:    path,
		Visited: visited,
	}, nil
}

type DFS struct {
}

func (b DFS) Evaluate(parameters RoutingRequest) (RoutingResponse, *error) {
	panic("implement me")
}

type AStar struct {
}

func (b AStar) Evaluate(parameters RoutingRequest) (RoutingResponse, *error) {
	panic("implement me")
}
