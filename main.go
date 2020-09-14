package main

import (
	"bufio"
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

func main() {
	//generateimage.GenerateImages()
	var graph graph.Graph

	data := readFile(IMAGEINFODIR)
	for _, line := range data {
		parts := strings.Split(line, "->")
		url := parts[0]
		tags := strings.Split(parts[1], ",")

		for _, tag := range tags {
			graph.AddEdge(tag, url)
		}
	}
	graph.String()
}
