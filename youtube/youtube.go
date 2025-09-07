package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"youtube-rss/utils"
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

type VideoInfo struct {
	PublishedAt time.Time `json:"publishedAt"`
	ChannelId string `json:"channelId"`
	Title string `json:"title"`
	Description string `json:"description"`
}

type Item struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id string `json:"id"`
	Snippet VideoInfo `json:"snippet"`
}

type Body struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	Id string `json:"id"`
	NextPageToken string `json:"nextPageToken"`
	Items []Item `json:"items"`
}

const BASE_URL string = "https://www.googleapis.com" +
	"/youtube/v3/playlistItems" +
	"?part=snippet&maxResults=25"

func GetPlaylist(apiKey string, playlistId string) (Body, error) {
	url := BASE_URL + fmt.Sprintf("&playlistId=%v", playlistId) +
		fmt.Sprintf("&key=%v", apiKey)

	res, err := http.Get(url)
	utils.LogError(err)

	var resBody Body
	err = json.NewDecoder(res.Body).Decode(&resBody)
	utils.LogWarning(err)

	nextPageToken := resBody.NextPageToken
	for len(nextPageToken) != 0 {
		var tmpBody Body

		nextPageParam := fmt.Sprintf("&pageToken=%s", nextPageToken)
		res, err := http.Get(url + nextPageParam)
		utils.LogError(err)

		err = json.NewDecoder(res.Body).Decode(&tmpBody)
		utils.LogWarning(err)

		resBody.Items = append(resBody.Items, tmpBody.Items[:]...)

		nextPageToken = tmpBody.NextPageToken
	}

	return resBody, err
}
