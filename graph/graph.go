package graph

import (
	"fmt"

	"github.com/obayanju/image-repo/set"
)

type Graph struct {
	edges map[string]*set.StringSet
}

func (g *Graph) AddEdge(key, value string) {
	if g.edges == nil {
		g.edges = make(map[string]*set.StringSet)
	}

	if g.edges[key] == nil {
		g.edges[key] = &set.StringSet{}
	}
	g.edges[key].Add(value)
}

func (g *Graph) GetValues(key string) *set.StringSet {
	return g.edges[key]
}

func (g *Graph) String() {
	for k, v := range g.edges {
		fmt.Printf("%s -> %v\n", k, v)
	}
}
