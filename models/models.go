package models

import "time"

type Playlist struct {
	Id *int
	Name string
	YoutubeId string
}

type Video struct {
	Id *int
	PlaylistId int
	Name string
	YoutubeId string
	Watched bool
	PublishedAt *time.Time
	CreatedAt *time.Time
	Episode *int
}
