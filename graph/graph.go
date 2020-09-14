package graph

import "fmt"

type Graph struct {
	edges map[string][]string
}

func (g *Graph) AddEdge(key, value string) {
	if g.edges == nil {
		g.edges = make(map[string][]string)
	}
	g.edges[key] = append(g.edges[key], value)
}

func (g *Graph) GetValues(key string) []string {
	return g.edges[key]
}

func (g *Graph) String() {
	for k, v := range g.edges {
		fmt.Printf("%s -> %v\n", k, v)
	}
}
