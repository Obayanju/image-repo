package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/obayanju/image-repo/generateimage"
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

func startServer() {
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

func main() {
	var tags, amount string
	var runServer bool

	flag.StringVar(&tags, "tags", "", "tags of images to match\nvalue must be comma sepearated with no space")
	flag.StringVar(&amount, "amount", "", "number of images of each tag to generate\nvalue must be comma sepearated with no space")
	flag.BoolVar(&runServer, "run", false, "run server")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	var tagSl []string
	var amountSl []int
	if tags != "" {
		tagSl = strings.Split(tags, ",")
	}
	if amount != "" {
		for _, v := range strings.Split(amount, ",") {
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
			}
			amountSl = append(amountSl, i)
		}

		remaining := len(tagSl) - len(amountSl)
		for i := 0; i < remaining; i++ {
			amountSl = append(amountSl, 10)
		}
	}
	if tags != "" {
		generateimage.GenerateImages(tagSl, amountSl)
	}

	if runServer {
		startServer()
	}
}
