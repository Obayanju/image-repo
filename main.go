package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/obayanju/image-repo/graph"
	"github.com/obayanju/image-repo/set"
)

const IMAGEINFODIR = "./images.txt"

type ImageTags struct {
	items map[string]*set.StringSet
}

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

func getImageMatch(tags []string, graph *graph.Graph) ImageTags {
	imageMatches := ImageTags{}
	for _, tag := range tags {
		urls := graph.GetValues(tag)
		if urls != nil {
			for _, url := range urls.Items() {
				if imageMatches.items == nil {
					imageMatches.items = make(map[string]*set.StringSet)
				}
				if imageMatches.items[url] == nil {
					imageMatches.items[url] = &set.StringSet{}
				}
				tagSet := imageMatches.items[url]
				tagSet.Add(tag)
			}
		}
	}
	return imageMatches
}

func main() {
	var graph graph.Graph

	data := readFile(IMAGEINFODIR)
	addImageTagToGraph(data, &graph)

	tags := []string{"nature", "volleyball", "mountains", "sunset", "fashion"}
	imageMatches := ImageTags{}
	imageMatches = getImageMatch(tags, &graph)
	for url, tagSet := range imageMatches.items {
		fmt.Printf("%s -> %v\n", url, tagSet.Items())
	}
}
