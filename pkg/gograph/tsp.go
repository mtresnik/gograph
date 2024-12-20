package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"github.com/mtresnik/goutils/pkg/goutils"
	"math"
	"math/rand"
	"time"
)

type TSPRequest struct {
	Graph            Graph
	DistanceFunction *gomath.DistanceFunction
	MaxIterations    int
}

type TSPResponse struct {
	Path Path
}

type TSP func(TSPRequest) TSPResponse

var GreedyTSP TSP = func(request TSPRequest) TSPResponse {
	visited := make([]Vertex, 0)
	visitedSet := make(map[int64]bool)
	visitedEdges := make([]Edge, 0)

	vertices := request.Graph.GetVertices()
	numVertices := len(vertices)
	if numVertices < 2 {
		return TSPResponse{
			Path: NewSimplePath([]Edge{}),
		}
	}

	distanceFunction := gomath.EuclideanDistance
	if request.DistanceFunction != nil {
		distanceFunction = *request.DistanceFunction
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	startIndex := random.Intn(numVertices)
	startVertex := vertices[startIndex]
	currentVertex := startVertex

	visited = append(visited, currentVertex)
	visitedSet[currentVertex.Hash()] = true

	for len(visited) < numVertices {
		currentEdges := currentVertex.GetEdges()

		nextEdges := goutils.Filter(currentEdges, func(edge Edge) bool {
			return !goutils.SetContains(visitedSet, ToVertex(edge.To()).Hash())
		})

		if len(nextEdges) == 0 {
			var closestVertex Vertex
			minDistance := math.MaxFloat64

			for _, vertex := range vertices {
				if !visitedSet[vertex.Hash()] {
					distance := distanceFunction(currentVertex, vertex)
					if distance < minDistance {
						minDistance = distance
						closestVertex = vertex
					}
				}
			}

			newEdge := NewEdge(currentVertex, closestVertex)
			visitedEdges = append(visitedEdges, newEdge)
			currentVertex = closestVertex
			visited = append(visited, currentVertex)
			visitedSet[currentVertex.Hash()] = true
		} else {
			closestEdge := goutils.MinBy(nextEdges, func(edge Edge) float64 {
				return edge.DistanceCached(distanceFunction)
			})
			currentVertex = ToVertex(closestEdge.To())
			visited = append(visited, currentVertex)
			visitedSet[currentVertex.Hash()] = true
			visitedEdges = append(visitedEdges, closestEdge)
		}
	}

	path := NewSimplePath(visitedEdges).Wrap()
	return TSPResponse{Path: path}
}

type RepeatTSP struct {
	InternalTSP TSP
	Repeat      int
}

func (tsp RepeatTSP) Eval(request TSPRequest) TSPResponse {
	var minPath Path
	for i := 0; i < tsp.Repeat; i++ {
		response := tsp.InternalTSP(request)
		if minPath == nil || response.Path.Distance() < minPath.Distance() {
			minPath = response.Path
		}
	}
	return TSPResponse{
		Path: minPath,
	}
}

var RandomTSP TSP = func(request TSPRequest) TSPResponse {
	vertices := request.Graph.GetVertices()
	numVertices := len(vertices)
	if numVertices < 2 {
		return TSPResponse{Path: NewSimplePath([]Edge{})}
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	startIndex := random.Intn(numVertices)
	startVertex := vertices[startIndex]
	currentVertex := startVertex

	visitedSet := make(map[int64]bool)
	visitedEdges := make([]Edge, 0)
	visitedOrder := make([]Vertex, 0)

	visitedSet[currentVertex.Hash()] = true
	visitedOrder = append(visitedOrder, currentVertex)

	for len(visitedOrder) < numVertices {
		currentEdges := currentVertex.GetEdges()

		availableEdges := goutils.Filter(currentEdges, func(edge Edge) bool {
			return !goutils.SetContains(visitedSet, ToVertex(edge.To()).Hash())
		})

		if len(availableEdges) == 0 {
			unvisitedVertex := vertices[random.Intn(numVertices)]
			for visitedSet[unvisitedVertex.Hash()] {
				unvisitedVertex = vertices[random.Intn(numVertices)]
			}
			newEdge := NewEdge(currentVertex, unvisitedVertex)
			visitedEdges = append(visitedEdges, newEdge)
			currentVertex = unvisitedVertex
		} else {
			randomEdge := availableEdges[random.Intn(len(availableEdges))]
			currentVertex = ToVertex(randomEdge.To())
			visitedEdges = append(visitedEdges, randomEdge)
		}
		visitedOrder = append(visitedOrder, currentVertex)
		visitedSet[currentVertex.Hash()] = true
	}
	path := NewSimplePath(visitedEdges).Wrap()
	return TSPResponse{Path: path}
}
