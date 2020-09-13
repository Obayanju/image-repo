package graph

import "fmt"

type Graph struct {
	edges map[string][]string
}

func (g *Graph) AddEdge(key, value string) {
	g.edges[key] = append(values, value)
}

func (g *Graph) GetValues(key string) []string {
	return g.edges[key]
}

func (g *Graph) String() {
	for k, v := range g.edges {
		fmt.Printf("%s -> %v\n", k, v)
	}
}
