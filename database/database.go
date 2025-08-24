package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"youtube-rss/models"
	"youtube-rss/utils"
	"youtube-rss/youtube"
)

var BASE_FOLDER string = utils.GetBaseFolder()

func GetDatabase() (*sql.DB, error) {
	return  sql.Open(
		"sqlite3",
		filepath.Join(
			BASE_FOLDER,
			"youtube-rss.db?_foreign_keys=on",
		),
	)
}

func createTables() {
	db, err := GetDatabase()
	utils.LogError(err)

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
	utils.LogError(err)

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
			published_at datetime not null,
			created_at datetime not null default current_timestamp,
			foreign key (playlist_id) references playlists(id)
		)
	`)
	utils.LogError(err)
}

func isTablesInitialized() bool {
	db, err := GetDatabase()
	utils.LogError(err)

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

func GetPlaylists() []models.Playlist {
	db, err := GetDatabase()
	utils.LogError(err)

	// TODO: move query to another place
	rows, err := db.Query("select id, name, youtube_id from playlists;")
	playlists := make([]models.Playlist, 0)

	for rows.Next() {
		var playlist models.Playlist
		rows.Scan(
			&playlist.Id,
			&playlist.Name,
			&playlist.YoutubeId)
		playlists = append(playlists, playlist)
	}

	return playlists
}

func GetPlaylistByName(name string) models.Playlist {
	db, err := GetDatabase()
	utils.LogError(err)

	// TODO: move query to another place
	row := db.QueryRow(`
		select
			id,
			name,
			youtube_id
		from playlists
		where name = $1
	`, name)

	var playlist models.Playlist
	row.Scan(
		&playlist.Id,
		&playlist.Name,
		&playlist.YoutubeId,
	)

	return playlist
}

func AddPlaylist(playlist models.Playlist) models.Playlist {
	db, err := GetDatabase()
	utils.LogError(err)

	row := db.QueryRow(`
		INSERT INTO playlists (name, youtube_id)
		values ($1, $2)
		returning id, name, youtube_id
	`, playlist.Name, playlist.YoutubeId)
	row.Scan(&playlist.Id, &playlist.Name, &playlist.YoutubeId)

	return playlist
}

func RemovePlaylist(playlist models.Playlist) {
	db, err := GetDatabase()
	utils.LogError(err)

	_, err = db.Exec(`
		DELETE FROM videos
		WHERE playlist_id = $1
	`, *playlist.Id)
	utils.LogError(err)

	_, err = db.Exec(`
		DELETE FROM playlists
		WHERE youtube_id = $1
	`, playlist.YoutubeId)
	utils.LogError(err)
}

func GetVideosByPlaylist(playlistYoutubeId string) []models.Video {
	db, err := GetDatabase()
	utils.LogError(err)

	var videos []models.Video
	rows, err := db.Query(`
		SELECT
			v.playlist_id,
			v.name,
			v.youtube_id,
			v.watched,
			v.published_at
		FROM playlists p
			JOIN videos v
				on v.playlist_id = p.id
		WHERE p.youtube_id = $1
		ORDER BY v.published_at desc
	`, playlistYoutubeId)
	utils.LogError(err)

	for rows.Next() {
		var video models.Video
		rows.Scan(
			&video.PlaylistId,
			&video.Name,
			&video.YoutubeId,
			&video.Watched,
			&video.PublishedAt,
		)
		videos = append(videos, video)
	}

	return videos
}

func AddVideoIfNotRegistered(playlistId int, video youtube.Item) bool {
	db, err := GetDatabase()
	utils.LogError(err)

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
		return false
	}

	_, err = db.Exec(`
		INSERT INTO videos (playlist_id, name, youtube_id, published_at)
		VALUES ($1, $2, $3, $4)
	`, playlistId, video.Snippet.Title, video.Id, video.Snippet.PublishedAt)
	utils.LogError(err)

	return true
}

func AddPlaylistVideosIfNotRegistered(playlist models.Playlist) []string {
	videos, err := youtube.GetPlaylist(os.Getenv("KEY"), playlist.YoutubeId)
	utils.LogError(err)

	newVideosNames := []string{}
	for _, video := range videos.Items {
		if AddVideoIfNotRegistered(*playlist.Id, video) {
			newVideosNames = append(newVideosNames, video.Snippet.Title)
		}
	}
	return newVideosNames
}

func MarkVideoAsWatched(playlistId int, videoTitle string) {
	db, err := GetDatabase()
	utils.LogError(err)

	_, err = db.Exec(`
		UPDATE videos
		SET watched=true
		WHERE name=$1
			and playlist_id = $2
	`, videoTitle, playlistId)
	utils.LogError(err)
}

func MarkAllPlaylistVideosAsWatched(playlistId int) {
	db, err := GetDatabase()
	utils.LogError(err)

	_, err = db.Exec(`
		UPDATE videos
		SET watched=true
		WHERE playlist_id = $2
	`, playlistId)
	utils.LogError(err)
}
