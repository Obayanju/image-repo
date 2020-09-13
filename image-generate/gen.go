package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var API_KEY = os.Getenv("PIXABAY_KEY")

const (
	PixabayURL = "https://pixabay.com/api/"
)

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

func getImage(category, downloadDir string) error {
	param := url.Values{}
	param.Set("key", API_KEY)
	param.Set("category", category)
	param.Set("per_page", "3")
	queryString := param.Encode()

	url := fmt.Sprintf("%s?%s", PixabayURL, queryString)
	fmt.Println(url)
	byt := getDataFromURL(url)

	jsonData := unmarshalJSON(byt)
	return nil
}

func main() {
	category := "fashion"
	downloadDir := "."
	if err := getImage(category, downloadDir); err != nil {
		log.Fatalf("There was an error downloading an image -> %v", err)
	}
}
