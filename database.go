package main

import (
	"database/sql"
	"fmt"

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
