package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var API_KEY = os.Getenv("PIXABAY_KEY")

const (
	PixabayURL = "https://pixabay.com/api/"
)

type image struct {
	url  string
	tags []string
}

func getDataFromURL(url string) []byte {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func unmarshalJSON(byt []byte) map[string]interface{} {
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		log.Fatal(err)
	}
	return dat
}

func getImage(category string) []image {
	param := url.Values{}
	param.Set("key", API_KEY)
	param.Set("category", category)
	param.Set("per_page", "3")
	queryString := param.Encode()

	url := fmt.Sprintf("%s?%s", PixabayURL, queryString)
	byt := getDataFromURL(url)

	jsonData := unmarshalJSON(byt)
	hits := jsonData["hits"].([]interface{})

	images := []image{}
	for _, v := range hits {
		v := v.(map[string]interface{})
		tags := strings.Split(v["tags"].(string), ",")
		tags = append(tags, category)

		img := image{
			url:  v["previewURL"].(string),
			tags: tags,
		}
		images = append(images, img)
	}
	return images
}

func main() {
	category := "fashion"
	var images []image
	images = append(images, getImage(category)...)
	for _, img := range images {
		fmt.Printf("%s %v\n", img.url, strings.Join(img.tags, ","))
	}
}
