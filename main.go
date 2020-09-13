package main

import (
	"bufio"
	"log"
	"os"
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
	data := readFile(IMAGEINFODIR)
}
