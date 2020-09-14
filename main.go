package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/obayanju/image-repo/graph"
)

const IMAGEINFODIR = "./images.txt"

func readFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var out []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return out
}

// adds the images and tags from the text file to the graph
func addImageTagToGraph(lines []string, graph *graph.Graph) {
	for _, line := range lines {
		parts := strings.Split(line, "->")
		url := parts[0]
		tags := strings.Split(parts[1], ",")

		for _, tag := range tags {
			graph.AddEdge(tag, url)
		}
	}
}

func getImageMatch(tags []string, graph *graph.Graph) []string {
	imageMatches := []string{}
	for _, tag := range tags {
		sets := graph.GetValues(tag)
		if sets != nil {
			imageMatches = append(imageMatches, sets.Items()...)
		}
	}
	return imageMatches
}

func main() {
	var graph graph.Graph

	data := readFile(IMAGEINFODIR)
	addImageTagToGraph(data, &graph)

	tags := []string{"nature", "volleyball", "mountain", "sunset", "fashion"}
	fmt.Println(getImageMatch(tags, &graph))
}
