package gograph

import (
	"fmt"
	"github.com/mtresnik/goutils/pkg/goutils"
	"math"
	"math/rand"
	"sort"
	"time"
)

type MSTRequest struct {
	Graph Graph
}

type MSTResponse struct {
	Graph Graph
}

type MST func(MSTRequest) MSTResponse

var KruskalMST MST = func(request MSTRequest) MSTResponse {
	retGraph := NewSimpleGraph()
	allSets := make([][]Vertex, 0)
	sortedEdges := make([]Edge, 0)
	visitedSortedEdges := map[int64]bool{}

	// Collect edges and initialize sets of vertices
	for _, vertex := range request.Graph.GetVertices() {
		for _, edge := range vertex.GetEdges() {
			// Ensure no duplicate edges are added
			if !goutils.SetContains(visitedSortedEdges, edge.Hash()) {
				sortedEdges = append(sortedEdges, edge)
				visitedSortedEdges[edge.Hash()] = true
			}
		}
		allSets = append(allSets, []Vertex{vertex})
	}

	// Sort edges
	sort.Slice(sortedEdges, func(i, j int) bool {
		edge1 := sortedEdges[i]
		edge2 := sortedEdges[j]
		edge1Value := edge1.DistanceCached() + 1
		edge2Value := edge2.DistanceCached() + 1
		return edge1Value < edge2Value
	})

	var validEdges []Edge = make([]Edge, 0)
	validEdgesSet := map[int64]bool{}

	for len(sortedEdges) > 0 && len(allSets) > 1 {
		currentEdge := sortedEdges[0]
		sortedEdges = sortedEdges[1:] // Remove the first edge

		// fromIndex := indexOfSetContainingVertex(allSets, currentEdge.From)
		fromIndex := indexOfSetContainingVertex(allSets, VertexFromSpatial(currentEdge.From()))
		if fromIndex != -1 {
			fromSet := allSets[fromIndex]
			allSets = append(allSets[:fromIndex], allSets[fromIndex+1:]...)

			toIndex := indexOfSetContainingVertex(allSets, VertexFromSpatial(currentEdge.To()))

			if toIndex != -1 {
				toSet := allSets[toIndex]
				allSets = append(allSets[:toIndex], allSets[toIndex+1:]...)

				// Merge the two sets
				joinedSet := map[int64]bool{}
				joined := make([]Vertex, 0)
				for _, vertex := range fromSet {
					if !goutils.SetContains(joinedSet, vertex.Hash()) {
						joinedSet[vertex.Hash()] = true
						joined = append(joined, vertex)
					}
				}
				for _, vertex := range toSet {
					if !goutils.SetContains(joinedSet, vertex.Hash()) {
						joinedSet[vertex.Hash()] = true
						joined = append(joined, vertex)
					}
				}
				allSets = append(allSets, joined)

				validEdges = append(validEdges, currentEdge)
				validEdgesSet[currentEdge.Hash()] = true
			} else {
				allSets = append(allSets, fromSet)
			}
		}
	}

	// Check for disjoint graph
	if len(allSets) != 1 {
		fmt.Println("Graph is disjoint!")
	}

	// Save valid edges and construct the resulting graph
	for _, edge := range validEdges {
		cloned := NewSimpleEdge(edge.From(), edge.To(), -1)
		retGraph.AddEdge(cloned)
		from := VertexFromSpatial(cloned.From())
		from.AddEdge(cloned)
		to := VertexFromSpatial(cloned.To())
		retGraph.AddVertex(from)
		retGraph.AddVertex(to)
	}
	return MSTResponse{Graph: retGraph}
}

func indexOfSetContainingVertex(sets [][]Vertex, vertex Vertex) int {
	for i, set := range sets {
		for _, v := range set {
			if v.Hash() == vertex.Hash() {
				return i
			}
		}
	}
	return -1
}

func PrimsMST(request MSTRequest) MSTResponse {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	retGraph := NewSimpleGraph()
	cheapestVertexCost := map[int64]float64{}
	edgeProvidingConnection := map[int64]*Edge{}
	graphVertices := request.Graph.GetVertices()
	notVisited := make([]Vertex, len(graphVertices))
	notVisitedSet := map[int64]bool{}
	copy(notVisited, graphVertices)
	for _, vertex := range notVisited {
		cheapestVertexCost[vertex.Hash()] = math.MaxFloat64
		edgeProvidingConnection[vertex.Hash()] = nil
		notVisitedSet[vertex.Hash()] = true
	}

	edgeForest := make([]Edge, 0)

	for len(notVisited) > 0 {
		toRemoveIndex := -1
		cheapestCost := math.MaxFloat64
		for i, vertex := range notVisited {
			if cheapestVertexCost[vertex.Hash()] < cheapestCost {
				cheapestCost = cheapestVertexCost[vertex.Hash()]
				toRemoveIndex = i
			}
		}
		if toRemoveIndex == -1 {
			toRemoveIndex = random.Intn(len(notVisited))
		}
		removed := notVisited[toRemoveIndex]
		notVisited = append(notVisited[:toRemoveIndex], notVisited[toRemoveIndex+1:]...)
		notVisitedSet[removed.Hash()] = false
		edgeRemoved, ok := edgeProvidingConnection[removed.Hash()]
		if ok && edgeRemoved != nil {
			edgeForest = append(edgeForest, *edgeRemoved)
		}
		for _, edge := range removed.GetEdges() {
			if goutils.SetContains(notVisitedSet, VertexFromSpatial(edge.To()).Hash()) {
				newCost := edge.DistanceCached()
				cheapestToCost, ok := cheapestVertexCost[VertexFromSpatial(edge.To()).Hash()]
				if !ok {
					cheapestToCost = math.MaxFloat64
				}
				if newCost < cheapestToCost {
					cheapestVertexCost[VertexFromSpatial(edge.To()).Hash()] = newCost
					edgeProvidingConnection[VertexFromSpatial(edge.To()).Hash()] = &edge
				}
			}
		}
	}
	for _, edge := range edgeForest {
		cloned := NewSimpleEdge(edge.From(), edge.To(), -1)
		retGraph.AddEdge(cloned)
		from := VertexFromSpatial(cloned.From())
		from.AddEdge(cloned)
		to := VertexFromSpatial(cloned.To())
		retGraph.AddVertex(from)
		retGraph.AddVertex(to)
	}

	return MSTResponse{Graph: retGraph}
}
