package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetDatabase() (*sql.DB, error) {
	return  sql.Open("sqlite3", "./youtube-rss.db?_foreign_keys=on")
}

func createTables() {
	db, err := GetDatabase()
	LogError(err)

	// TODO: move query to another place
	// Create playlists table
	_, err = db.Exec(`
		DROP TABLE if exists playlists;
		CREATE TABLE playlists (
			id integer primary key autoincrement,
			name varchar(64) not null,
			youtube_id varchar(64) not null,
			created_at datetime not null default current_timestamp
		)
	`)
	LogError(err)

	// TODO: move query to another place
	// Create videos table
	_, err = db.Exec(`
		DROP TABLE if exists videos;
		CREATE TABLE videos (
			id integer primary key autoincrement,
			playlist_id interger,
			name varchar(64) not null,
			youtube_id varchar(64) not null,
			watched bool not null default false,
			created_at datetime not null default current_timestamp,
			foreign key (playlist_id) references playlists(id)
		)
	`)
	LogError(err)
}

func isTablesInitialized() bool {
	db, err := GetDatabase()
	LogError(err)

	tables := []string{"playlists", "videos"}
	tablesLen := len(tables)
	tablesString := ""
	for i, tableName := range tables {
		tablesString += fmt.Sprintf("'%v'", tableName)
		if i + 1 < tablesLen {
			tablesString += ", "
		}
	}

	// TODO: move query to another place
	row := db.QueryRow(fmt.Sprintf(`
		SELECT count(*) FROM sqlite_master
		WHERE type='table' AND name in (%s);
	`, tablesString))

	var count int
	row.Scan(&count)

	return count == len(tables)
}

func InitDatabase() {
	if !isTablesInitialized() {
		createTables()
	}
}

func GetPlaylists() []Playlist {
	db, err := GetDatabase()
	LogError(err)

	// TODO: move query to another place
	rows, err := db.Query("select id, name, youtube_id from playlists;")
	playlists := make([]Playlist, 0)

	for rows.Next() {
		var playlist Playlist
		rows.Scan(
			&playlist.id,
			&playlist.name,
			&playlist.youtubeId)
		playlists = append(playlists, playlist)
	}

	return playlists
}

func GetPlaylistByName(name string) Playlist {
	db, err := GetDatabase()
	LogError(err)

	// TODO: move query to another place
	row := db.QueryRow(`
		select
			id,
			name,
			youtube_id
		from playlists
		where name = $1
	`, name)

	var playlist Playlist
	row.Scan(
		&playlist.id,
		&playlist.name,
		&playlist.youtubeId,
	)

	return playlist
}

func AddPlaylist(playlist Playlist) Playlist {
	db, err := GetDatabase()
	LogError(err)

	row := db.QueryRow(`
		INSERT INTO playlists (name, youtube_id)
		values ($1, $2)
		returning id, name, youtube_id
	`, playlist.name, playlist.youtubeId)
	row.Scan(&playlist.id, &playlist.name, &playlist.youtubeId)

	return playlist
}

func RemovePlaylist(playlist Playlist) {
	db, err := GetDatabase()
	LogError(err)

	_, err = db.Exec(`
		DELETE FROM videos
		WHERE playlist_id = $1
	`, *playlist.id)
	LogError(err)

	_, err = db.Exec(`
		DELETE FROM playlists
		WHERE youtube_id = $1
	`, playlist.youtubeId)
	LogError(err)
}

func GetVideosByPlaylist(playlistYoutubeId string) []Video {
	db, err := GetDatabase()
	LogError(err)

	var videos []Video
	rows, err := db.Query(`
		SELECT
			v.playlist_id,
			v.name,
			v.youtube_id,
			v.watched
		FROM playlists p
			JOIN videos v
				on v.playlist_id = p.id
		WHERE p.youtube_id = $1
		ORDER BY v.id
	`, playlistYoutubeId)
	LogError(err)

	for rows.Next() {
		var video Video
		rows.Scan(
			&video.playlistId,
			&video.name,
			&video.youtubeId,
			&video.watched,
		)
		videos = append(videos, video)
	}

	return videos
}

func AddVideoIfNotRegistered(playlistId int, video Item) {
	db, err := GetDatabase()
	LogError(err)

	row := db.QueryRow(`
		SELECT
			count(*)
		FROM videos
		WHERE playlist_id = $1
			and youtube_id = $2
	`, playlistId, video.Id)
	var count int
	row.Scan(&count)

	if count > 0 {
		return
	}

	_, err = db.Exec(`
		INSERT INTO videos (playlist_id, name, youtube_id)
		VALUES ($1, $2, $3)
	`, playlistId, video.Snippet.Title, video.Id)
	LogError(err)
}

func AddPlaylistVideosIfNotRegistered(playlist Playlist) {
	videos, err := GetPlaylist(os.Getenv("KEY"), playlist.youtubeId)
	LogError(err)

	for _, video := range videos.Items {
		AddVideoIfNotRegistered(*playlist.id, video)
	}
}

func AddVideos(apiKey string, playlist Playlist) {
	db, err := GetDatabase()
	LogError(err)

	statement, err := db.Prepare(`
		INSERT INTO videos (playlist_id, name, youtube_id)
		VALUES (?, ?, ?)
	`)

	body, err := GetPlaylist(apiKey, playlist.youtubeId)
	for _, item := range body.Items {
		_, err := statement.Exec(playlist.id, item.Snippet.Title, item.Id)
		LogError(err)
	}
}

func MarkVideoAsWatched(playlistId int, videoTitle string) {
	db, err := GetDatabase()
	LogError(err)

	_, err = db.Exec(`
		UPDATE videos
		SET watched=true
		WHERE name=$1
			and playlist_id = $2
	`, videoTitle, playlistId)
	LogError(err)
}

func MarkAllPlaylistVideosAsWatched(playlistId int) {
	db, err := GetDatabase()
	LogError(err)

	_, err = db.Exec(`
		UPDATE videos
		SET watched=true
		WHERE playlist_id = $2
	`, playlistId)
	LogError(err)
}
