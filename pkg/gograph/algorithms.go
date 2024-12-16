package gograph

import (
	"errors"
	"github.com/mtresnik/goutils/pkg/goutils"
)

type RoutingRequest struct {
	Start             Vertex
	Destination       Vertex
	Constraints       *[]Constraint
	MultiCostFunction *MultiCostFunction
	CostFunctions     *[]CostFunction
}

type RoutingResponse struct {
	Costs   []CostEntry
	Path    []Edge
	Visited map[int64]bool
}

type RoutingAlgorithm interface {
	VisitVertex(vertex Vertex)
	VisitEdge(edge Edge)
	Evaluate(parameters RoutingRequest) (RoutingResponse, *error)
}

type BFS struct {
	VertexVisitors []VertexVisitor
	EdgeVisitors   []EdgeVisitor
}

func (b BFS) VisitVertex(vertex Vertex) {
	for _, visitor := range b.VertexVisitors {
		visitor.Visit(vertex)
	}
}

func (b BFS) VisitEdge(edge Edge) {
	for _, visitor := range b.EdgeVisitors {
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
	var costFunctions []CostFunction
	if parameters.CostFunctions != nil {
		costFunctions = *parameters.CostFunctions
	} else {
		costFunctions = []CostFunction{EuclideanDistanceCostFunction{}}
	}
	initialCosts := make([]CostEntry, len(costFunctions))
	for i, function := range costFunctions {
		initialCosts[i] = CostEntry{
			Type:        function.GetType(),
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
			costs := make([]CostEntry, len(costFunctions))
			for i, costFunction := range costFunctions {
				cost := costFunction.Eval(curr, toVertex)
				costs[i] = CostEntry{
					Type:        costFunction.GetType(),
					Accumulated: curr.Costs[i].Total,
					Current:     cost,
					Total:       curr.Costs[i].Total + cost,
				}
			}
			if !goutils.SetContains(visited, hashOrId) {
				nextWrapper := NewVertexWrapper(toVertex, costs)
				pass := true
				if constraints != nil {
					pass = CheckAllConstraints(curr, costs, *constraints...)
				}
				if pass {
					b.VisitEdge(edge)
					nextWrapper.Previous = &curr
					queue = append(queue, nextWrapper)
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
