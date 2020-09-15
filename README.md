# image-repo
Search images

## QuickStart
- Clone the repo using `git clone https://github.com/Obayanju/image-repo.git`
- Get your Pixabay API key
	- Sign up at [url](https://pixabay.com/accounts/register/?source=main_nav)
	- Your API key is [here](https://pixabay.com/api/docs/#api_search_images)
- Add your API key in the `API_KEY` variable in [gen.go](https://github.com/Obayanju/image-repo/blob/master/generateimage/gen.go)
- run `go run main.go -tags nature,fashion,people,food,sports -amount 10,5,8,6,9` to generate image urls and its tags. The file is stored at `./images.txt`
- run `go run main.go -run` to start the server. The server is located at `http://localhost:3000`
- run `curl "http://localhost:3000/?tags=fashion,food,sports"`for a GET request for images matching the specified tags. It returns a JSON `{image: tag}`

### Help
run `go run main.go --help`
