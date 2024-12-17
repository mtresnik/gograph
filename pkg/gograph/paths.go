package gograph

import "hash/fnv"

type Path interface {
	GetEdges() []Edge
	GetCost() float64
	Length() int
	Hash() int64
	Id() int64
}

func PathHashOrId(path Path) int64 {
	if path.Id() > 0 {
		return path.Id()
	}
	return path.Hash()
}

type SimplePath struct {
	Edges []Edge
	Cost  float64
	id    int64
	hash  int64
}

func NewSimplePath(edges []Edge, cost float64) *SimplePath {
	return &SimplePath{edges, cost, -1, -1}
}

func (p *SimplePath) Length() int {
	return len(p.Edges)
}

func (p *SimplePath) GetEdges() []Edge {
	return p.Edges
}

func (p *SimplePath) GetCost() float64 {
	return p.Cost
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
