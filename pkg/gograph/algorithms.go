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
	CostFunctions     *map[string]CostFunction
}

type RoutingResponse struct {
	Costs   map[string]CostEntry
	Path    []Edge
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

func (b BFS) VisitVertex(vertex Vertex) {
	for _, visitor := range b.VertexListeners {
		visitor.Visit(vertex)
	}
}

func (b BFS) VisitEdge(edge Edge) {
	for _, visitor := range b.EdgeListeners {
		visitor.Visit(edge)
	}
}

func Backtrack(vertex VertexWrapper) []Edge {
	var path []Edge
	var currentWrapper = vertex
	for currentWrapper.Previous != nil {
		path = append(path, currentWrapper.Previous.Inner.GetEdge(currentWrapper.Inner))
		currentWrapper = *currentWrapper.Previous
	}
	return path
}

func (b BFS) Evaluate(parameters RoutingRequest) (RoutingResponse, *error) {
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
	visited[start.Id()] = true
	var curr = startWrapper
	for len(queue) > 0 {
		curr = queue[0]
		if curr.Hash() == destination.Hash() {
			break
		}
		for _, edge := range curr.Inner.GetEdges() {
			toVertex := ToVertex(edge.To())
			hashOrId := HashOrId(toVertex)
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
					successor.Previous = &curr
					queue = append(queue, successor)
					visited[hashOrId] = true
					b.VisitVertex(toVertex)
				}
			}
		}
		if len(queue) == 0 {
			err := errors.New("no path found, try relaxing the constraints")
			return RoutingResponse{}, &err
		}
	}
	edges := Backtrack(curr)
	return RoutingResponse{
		Costs:   curr.Costs,
		Path:    edges,
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
