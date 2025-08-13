package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ThumbnailInfo struct {
	URL string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type Thumbnail struct {
	Default *ThumbnailInfo `json:"default"`
	Medium *ThumbnailInfo `json:"medium"`
	High *ThumbnailInfo `json:"high"`
	Standart *ThumbnailInfo `json:"standart"`
	MaxRes *ThumbnailInfo `json:"maxres"`
}

type ItemInfo struct {
	PublishedAt time.Time `json:"publishedAt"`
	ChannelId string `json:"channelId"`
	Title string `json:"title"`
	Description string `json:"description"`
}

type Item struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id string `json:"id"`
	Snippet ItemInfo `json:"snippet"`
}

type Body struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id string `json:"id"`
	Items []Item `json:"items"`
}

func main() {
	run([]string{"add", "the office"})
	run([]string{"add", "the office", "idk2"})
	run([]string{"list-playlists"})
	run([]string{""})

	os.Exit(1)
	LoadEnv()

	key := os.Getenv("KEY")
	if key == "" {
		log.Panicln("Invalid KEY")
	}

	const BASE_URL string = "https://www.googleapis.com" +
		"/youtube/v3/playlistItems" +
		"?part=snippet&maxResults=50&playlistId=PL2Fq-K0QdOQiJpufsnhEd1z3xOv2JMHuk"

	url := fmt.Sprintf("%v&key=%v", BASE_URL, key)
	fmt.Println(url)

	res, err := http.Get(url)
	LogError(err)

	var resBody Body
	err = json.NewDecoder(res.Body).Decode(&resBody)

	LogError(err)

	for i, item := range resBody.Items {
		fmt.Printf(
			"[%v] Episode %v: %v (%v)\n",
			item.Id,
			len(resBody.Items) - i,
			item.Snippet.Title,
			item.Snippet.PublishedAt,
		)
	}
}
