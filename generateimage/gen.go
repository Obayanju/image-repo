package generateimage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var API_KEY = os.Getenv("PIXABAY_KEY")

const (
	PixabayURL   = "https://pixabay.com/api/"
	IMAGEINFODIR = "./images.txt"
)

type Image struct {
	Url  string
	Tags []string
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

//  getImage queries the Pixabay API and returns results images matching
// category
func getImage(category, results string) []Image {
	i, err := strconv.Atoi(results)
	if err != nil {
		log.Fatalf("cannot convert %s to an int: %v", results, err)
	}
	if i < 3 {
		results = "3"
	}
	param := url.Values{}
	param.Set("key", API_KEY)
	param.Set("category", category)
	param.Set("per_page", results)
	queryString := param.Encode()

	url := fmt.Sprintf("%s?%s", PixabayURL, queryString)
	byt := getDataFromURL(url)

	jsonData := unmarshalJSON(byt)
	hits := jsonData["hits"].([]interface{})

	images := []Image{}
	for _, v := range hits {
		v := v.(map[string]interface{})
		tags := strings.Split(v["tags"].(string), ", ")
		tags = append(tags, category)

		img := Image{
			Url:  v["previewURL"].(string),
			Tags: tags,
		}
		images = append(images, img)
	}
	return images
}

func appendToFile(fName, body string) {
	f, err := os.OpenFile(fName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.Write([]byte(body)); err != nil {
		log.Println(err)
	}
}

func truncate(filename string, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("could not open file %q for truncation: %v", filename, err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("could not close file handler for %q after truncation: %v", filename, err)
	}
	return nil
}

func GenerateImages(tags []string, amount []int) {
	categories := [][]string{}
	for i, tag := range tags {
		tagAmount := []string{tag, strconv.Itoa(amount[i])}
		categories = append(categories, tagAmount)
	}

	var images []Image
	for _, category := range categories {
		images = append(images, getImage(category[0], category[1])...)
	}

	truncate(IMAGEINFODIR, 0644)
	for _, img := range images {
		content := fmt.Sprintf("%s->%v\n", img.Url, strings.Join(img.Tags, ","))
		appendToFile(IMAGEINFODIR, content)
	}
}
