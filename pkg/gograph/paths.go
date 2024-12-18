package gograph

import (
	"github.com/mtresnik/gomath/pkg/gomath"
	"hash/fnv"
)

type Path interface {
	GetEdges() []Edge
	Length() int
	Hash() int64
	Id() int64
}

func GetPathCost(path Path, pCostFunctions *map[string]CostFunction) map[string]CostEntry {
	costFunctions, initialCosts := GenerateInitialCosts(pCostFunctions)
	if path.Length() == 0 {
		return initialCosts
	}
	startWrapper := NewVertexWrapper(ToVertex(path.GetEdges()[0].From()), initialCosts)
	var curr = startWrapper
	for _, edge := range path.GetEdges() {
		toVertex := ToVertex(edge.To())
		nextCosts := GenerateNextCosts(curr, toVertex, costFunctions)
		curr = NewVertexWrapper(toVertex, nextCosts)
	}
	return curr.Costs
}

func GetPathDistance(path Path, distanceFunction ...gomath.DistanceFunction) float64 {
	sum := 0.0
	for _, edge := range path.GetEdges() {
		sum += edge.Distance(distanceFunction...)
	}
	return sum
}

func PathHashOrId(path Path) int64 {
	if path.Id() > 0 {
		return path.Id()
	}
	return path.Hash()
}

type SimplePath struct {
	Edges []Edge
	id    int64
	hash  int64
}

func NewSimplePath(edges []Edge) *SimplePath {
	return &SimplePath{edges, -1, -1}
}

func (p *SimplePath) Length() int {
	return len(p.Edges)
}

func (p *SimplePath) GetEdges() []Edge {
	return p.Edges
}

func (p *SimplePath) Hash() int64 {
	if p.hash != -1 {
		return p.hash
	}
	edgeHashes := make([]int64, 0, len(p.Edges))
	hasher := fnv.New64a()
	for _, key := range edgeHashes {
		var buf [8]byte
		for i := 0; i < 8; i++ {
			buf[i] = byte(key >> (i * 8))
		}
		_, err := hasher.Write(buf[:])
		if err != nil {
			return 0
		}
	}
	p.hash = int64(hasher.Sum64())
	return p.hash
}

func (p *SimplePath) Id() int64 {
	return p.id
}
