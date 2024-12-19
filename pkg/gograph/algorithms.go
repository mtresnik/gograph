package gograph

import (
	"github.com/mtresnik/goutils/pkg/goutils"
	"maps"
	"math"
)

type RoutingAlgorithmRequest struct {
	Start             Vertex
	Destination       Vertex
	Constraints       *map[string][]Constraint
	MultiCostFunction *MultiCostFunction
	CostFunctions     *map[string]CostFunction
	UpdateListeners   *[]RoutingAlgorithmUpdateListener
	Algorithm         RoutingAlgorithm
}

type RoutingAlgorithmResponse struct {
	Costs     map[string]CostEntry
	Path      Path
	Visited   map[int64]bool
	Completed bool
}

type RoutingAlgorithmUpdateListener interface {
	Update(response RoutingAlgorithmResponse)
}

func VisitRoutingAlgorithmUpdateListeners(listeners []RoutingAlgorithmUpdateListener, response RoutingAlgorithmResponse) {
	for _, listener := range listeners {
		listener.Update(response)
	}
}

type RoutingAlgorithm func(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse

func EvaluateRoutingAlgorithm(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse {
	return parameters.Algorithm(parameters)
}

func Backtrack(vertex *VertexWrapper) []Edge {
	if vertex == nil {
		return []Edge{}
	}
	var path []Edge
	var currentWrapper = vertex
	for currentWrapper.Previous != nil && currentWrapper.Hash() != currentWrapper.Previous.Hash() {
		path = append(path, currentWrapper.Previous.Inner.GetEdge(currentWrapper.Inner))
		currentWrapper = currentWrapper.Previous
	}
	return path
}

var BFS RoutingAlgorithm = func(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse {
	start := parameters.Start
	destination := parameters.Destination
	constraints := parameters.Constraints
	costFunctions, initialCosts := GenerateInitialCosts(parameters.CostFunctions)
	startWrapper := NewVertexWrapper(start, initialCosts)
	startWrapper.Previous = nil
	queue := []*VertexWrapper{startWrapper}
	visited := make(map[int64]bool)
	var curr *VertexWrapper
	bestCombined := math.MaxFloat64
	var best = startWrapper
	updateListeners := make([]RoutingAlgorithmUpdateListener, 0)
	if parameters.UpdateListeners != nil {
		updateListeners = *parameters.UpdateListeners
	}
	for len(queue) > 0 {
		curr = queue[0]
		currCombined := MultiplyCosts(GenerateNextCosts(curr, destination, costFunctions)).Current
		if currCombined < bestCombined {
			best = NewVertexWrapper(curr.Inner, curr.Costs)
			best.Previous = curr.Previous
			bestCombined = currCombined
		}
		VisitRoutingAlgorithmUpdateListeners(updateListeners, RoutingAlgorithmResponse{
			best.Costs,
			NewSimplePath(Backtrack(best)),
			visited,
			false})
		queue = queue[1:]
		visited[VertexHashOrId(curr)] = true
		if curr.Hash() == destination.Hash() {
			break
		}
		for _, edge := range curr.Inner.GetEdges() {
			toVertex := ToVertex(edge.To())
			hashOrId := VertexHashOrId(toVertex)
			nextCosts := GenerateNextCosts(curr, toVertex, costFunctions)
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
					previousWrapper := curr
					successor.Previous = previousWrapper
					queue = append(queue, successor)
					visited[hashOrId] = true
				}
			}
		}
		if len(queue) == 0 {
			println("no path found, try relaxing the constraints")
			return RoutingAlgorithmResponse{
				best.Costs,
				NewSimplePath(Backtrack(best)),
				visited,
				false}
		}

	}
	edges := Backtrack(curr)
	path := NewSimplePath(edges)
	finalCosts := initialCosts
	if curr != nil {
		finalCosts = curr.Costs
	}
	response := RoutingAlgorithmResponse{
		Costs:     finalCosts,
		Path:      path,
		Visited:   visited,
		Completed: true,
	}
	VisitRoutingAlgorithmUpdateListeners(updateListeners, response)

	return response
}

var DFS RoutingAlgorithm = func(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse {
	panic("implement me")
}

var AStar RoutingAlgorithm = func(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse {
	panic("implement me")
}
