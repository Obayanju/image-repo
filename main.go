package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/obayanju/image-repo/graph"
	"github.com/obayanju/image-repo/set"
)

const IMAGEINFODIR = "./images.txt"

type ImageTags struct {
	items map[string]*set.StringSet
}

type Params struct {
	Tags []string `query:"tags"`
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

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		params := new(Params)

		if err := c.QueryParser(params); err != nil {
			return err
		}

		// workaround for a bug in fiber
		// https://github.com/gofiber/fiber/issues/782
		if len(params.Tags) == 1 {
			params.Tags = strings.Split(params.Tags[0], ",")
		}

		urlTagMap := fiber.Map{}
		imageMatches := ImageTags{}
		imageMatches = getImageMatch(params.Tags, &graph)
		for url, tagSet := range imageMatches.items {
			urlTagMap[url] = tagSet.Items()
		}

		return c.JSON(urlTagMap)
	})

	app.Listen(":3000")
}
