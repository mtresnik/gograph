package gograph

import (
	"container/heap"
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
	CostCombiner      *CostCombiner
	UpdateListeners   *[]RoutingAlgorithmUpdateListener
	Algorithm         RoutingAlgorithm
	ExplorationFactor float64
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
	costCombiner := MultiplicativeCostCombiner
	if parameters.CostCombiner != nil {
		costCombiner = *parameters.CostCombiner
	}
	startWrapper := NewVertexWrapper(start, initialCosts, costCombiner)
	startWrapper.Previous = nil
	startWrapper.Combined.Accumulated = 0
	startWrapper.Combined.Current = costCombiner(GenerateNextCosts(startWrapper, destination, costFunctions)).Current
	startWrapper.Combined.Total = costCombiner(GenerateNextCosts(startWrapper, destination, costFunctions)).Current

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
		queue = queue[1:]
		if goutils.SetContains(visited, VertexHashOrId(curr)) {
			continue
		}
		visited[VertexHashOrId(curr)] = true
		currCombined := costCombiner(GenerateNextCosts(curr, destination, costFunctions)).Current
		if currCombined < bestCombined {
			best = NewVertexWrapper(curr.Inner, curr.Costs, costCombiner)
			best.Previous = curr.Previous
			bestCombined = currCombined
			if len(updateListeners) > 0 {
				VisitRoutingAlgorithmUpdateListeners(updateListeners, RoutingAlgorithmResponse{
					best.Costs,
					NewSimplePath(Backtrack(best)),
					visited,
					false})
			}
		}
		if curr.Hash() == destination.Hash() {
			break
		}
		for _, edge := range curr.Inner.GetEdges() {
			toVertex := ToVertex(edge.To())
			hashOrId := VertexHashOrId(toVertex)
			nextCosts := GenerateNextCosts(curr, toVertex, costFunctions)
			successor := NewVertexWrapper(toVertex, nextCosts, costCombiner)
			successor.Combined.Total = successor.Combined.Accumulated
			if !goutils.SetContains(visited, hashOrId) {
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
					g := curr.Combined.Accumulated + successor.Combined.Current
					f := g
					if successor.Previous == nil || g < successor.Combined.Accumulated {
						successor.Previous = previousWrapper
						successor.Combined.Accumulated = g
						successor.Combined.Total = f
						queue = append(queue, successor)
					}
				}
			}
		}
	}

	edges := Backtrack(curr)
	path := NewSimplePath(edges)
	if len(edges) == 0 {
		path = NewSimplePath(Backtrack(best))
	}
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

type PathState struct {
	vertex        Vertex
	edges         []Edge
	nextEdgeIndex int
	costs         map[string]CostEntry
	accumulated   float64
	current       float64
	total         float64
	previous      *PathState
}

func PathStateToVertexWrapper(state *PathState) *VertexWrapper {
	ret := NewVertexWrapper(state.vertex, state.costs, MultiplicativeCostCombiner)
	ret.Combined.Accumulated = state.accumulated
	ret.Combined.Current = state.current
	ret.Combined.Total = state.total
	return ret
}

func NewPathState(vertex Vertex, costs map[string]CostEntry, accumulated, current, total float64, previous *PathState) *PathState {
	return &PathState{
		vertex:        vertex,
		edges:         vertex.GetEdges(),
		nextEdgeIndex: 0,
		costs:         costs,
		accumulated:   accumulated,
		current:       current,
		total:         total,
		previous:      previous,
	}
}

var DFS RoutingAlgorithm = func(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse {
	start := parameters.Start
	destination := parameters.Destination
	constraints := parameters.Constraints
	costFunctions, initialCosts := GenerateInitialCosts(parameters.CostFunctions)
	costCombiner := MultiplicativeCostCombiner
	if parameters.CostCombiner != nil {
		costCombiner = *parameters.CostCombiner
	}

	// Initialize the starting PathState
	startCosts := GenerateNextCosts(nil, start, costFunctions)
	startCombined := costCombiner(startCosts)
	startState := NewPathState(
		start,
		startCosts,
		0,
		startCombined.Current,
		startCombined.Current,
		nil,
	)

	stack := []*PathState{startState}
	visited := make(map[int64]bool)
	var curr *PathState
	bestCombined := math.MaxFloat64
	var best = startState

	updateListeners := make([]RoutingAlgorithmUpdateListener, 0)
	if parameters.UpdateListeners != nil {
		updateListeners = *parameters.UpdateListeners
	}

	for len(stack) > 0 {
		curr = stack[len(stack)-1]

		if curr.nextEdgeIndex >= len(curr.edges) {
			stack = stack[:len(stack)-1]
			continue
		}

		vertexHash := VertexHashOrId(curr.vertex)
		if !visited[vertexHash] {
			visited[vertexHash] = true

			currCombined := costCombiner(GenerateNextCosts(PathStateToVertexWrapper(curr), destination, costFunctions)).Current
			if currCombined < bestCombined {
				best = NewPathState(
					curr.vertex,
					curr.costs,
					curr.accumulated,
					curr.current,
					curr.total,
					curr.previous,
				)
				bestCombined = currCombined
				if len(updateListeners) > 0 {
					VisitRoutingAlgorithmUpdateListeners(updateListeners, RoutingAlgorithmResponse{
						best.costs,
						NewSimplePath(BacktrackPathState(best)),
						visited,
						false,
					})
				}
			}

		}

		if curr.vertex.Hash() == destination.Hash() {
			break
		}

		edge := curr.edges[curr.nextEdgeIndex]
		curr.nextEdgeIndex++

		toVertex := ToVertex(edge.To())
		hashOrId := VertexHashOrId(toVertex)

		if !visited[hashOrId] {
			nextCosts := GenerateNextCosts(PathStateToVertexWrapper(curr), toVertex, costFunctions)
			combined := costCombiner(nextCosts)

			pass := true
			if constraints != nil {
				keys := maps.Keys(*constraints)
				for key := range keys {
					_, constraintsExist := (*constraints)[key]
					currCost, costExist := nextCosts[key]
					if constraintsExist && costExist {
						pass = CheckAllConstraints(PathStateToVertexWrapper(curr), currCost, key, *constraints)
						if !pass {
							break
						}
					}
				}
			}

			if pass {
				g := curr.accumulated + combined.Current
				f := g

				successor := NewPathState(
					toVertex,
					nextCosts,
					g,
					combined.Current,
					f,
					curr,
				)

				stack = append(stack, successor)
			}
		}
	}

	// Build final response
	edges := BacktrackPathState(curr)
	path := NewSimplePath(edges)
	if len(edges) == 0 {
		path = NewSimplePath(BacktrackPathState(best))
	}

	finalCosts := initialCosts
	if curr != nil {
		finalCosts = curr.costs
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

func BacktrackPathState(state *PathState) []Edge {
	if state == nil {
		return make([]Edge, 0)
	}

	path := make([]Edge, 0)
	curr := state
	for curr.previous != nil {
		for _, edge := range curr.previous.vertex.GetEdges() {
			if ToVertex(edge.To()).Hash() == curr.vertex.Hash() {
				path = append([]Edge{edge}, path...)
				break
			}
		}
		curr = curr.previous
	}
	return path
}

var AStar RoutingAlgorithm = func(parameters RoutingAlgorithmRequest) RoutingAlgorithmResponse {
	start := parameters.Start
	destination := parameters.Destination
	constraints := parameters.Constraints
	costFunctions, initialCosts := GenerateInitialCosts(parameters.CostFunctions)
	costCombiner := MultiplicativeCostCombiner
	if parameters.CostCombiner != nil {
		costCombiner = *parameters.CostCombiner
	}
	startWrapper := NewVertexWrapper(start, initialCosts, costCombiner)
	startWrapper.Previous = nil
	startWrapper.Combined.Accumulated = 0
	startWrapper.Combined.Current = costCombiner(GenerateNextCosts(startWrapper, destination, costFunctions)).Current
	startWrapper.Combined.Total = costCombiner(GenerateNextCosts(startWrapper, destination, costFunctions)).Current

	explorationFactor := parameters.ExplorationFactor
	if explorationFactor == 0 {
		explorationFactor = math.Max(costCombiner(GenerateNextCosts(startWrapper, destination, costFunctions)).Current*5.0, 5.0)
	}
	open := &PriorityQueue{}
	heap.Init(open)
	visited := make(map[int64]bool)
	var curr *VertexWrapper
	bestCombined := math.MaxFloat64
	var best = startWrapper
	updateListeners := make([]RoutingAlgorithmUpdateListener, 0)
	if parameters.UpdateListeners != nil {
		updateListeners = *parameters.UpdateListeners
	}
	PushPriorityQueue(open, startWrapper, startWrapper.Combined.Total)
	for open.Len() > 0 {
		curr, ok := PollPriorityQueue(open).(*VertexWrapper)
		if ok {
			if goutils.SetContains(visited, VertexHashOrId(curr)) {
				continue
			}
			visited[VertexHashOrId(curr)] = true
			currCombined := costCombiner(GenerateNextCosts(curr, destination, costFunctions)).Current
			if currCombined < bestCombined {
				best = NewVertexWrapper(curr.Inner, curr.Costs, costCombiner)
				best.Previous = curr.Previous
				bestCombined = currCombined
				if len(updateListeners) > 0 {
					VisitRoutingAlgorithmUpdateListeners(updateListeners, RoutingAlgorithmResponse{
						best.Costs,
						NewSimplePath(Backtrack(best)),
						visited,
						false})
				}
			}
			if curr.Hash() == destination.Hash() {
				break
			}
			for _, edge := range curr.Inner.GetEdges() {
				toVertex := ToVertex(edge.To())
				hashOrId := VertexHashOrId(toVertex)
				nextCosts := GenerateNextCosts(curr, toVertex, costFunctions)
				successor := NewVertexWrapper(toVertex, nextCosts, costCombiner)
				successor.Combined.Total = successor.Combined.Accumulated + explorationFactor*costCombiner(GenerateNextCosts(successor, destination, costFunctions)).Current
				if !goutils.SetContains(visited, hashOrId) {
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
						g := curr.Combined.Accumulated + successor.Combined.Current
						h := costCombiner(GenerateNextCosts(successor, destination, costFunctions)).Current
						f := g * math.Pow(h+1.0, explorationFactor)
						if successor.Previous == nil || g < successor.Combined.Accumulated {
							successor.Previous = previousWrapper
							successor.Combined.Accumulated = g
							successor.Combined.Current = h
							successor.Combined.Total = f
							PushPriorityQueue(open, successor, f)
						}
					}
				}
			}
		}
	}

	edges := Backtrack(curr)
	path := NewSimplePath(edges)
	if len(edges) == 0 {
		path = NewSimplePath(Backtrack(best))
	}
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
