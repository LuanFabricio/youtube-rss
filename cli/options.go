package cli

import (
	"fmt"

	db "youtube-rss/database"
	"youtube-rss/models"
)

var options []Option = []Option {
	{
		name: "list",
		args: []string{},
		description: "List the playlists.",
		callback: func(args []string) {
			playlists := db.GetPlaylists()
			display := ""
			for _, playlist := range playlists {
				display += fmt.Sprintf(
					"- [%v]: %v\n",
					playlist.YoutubeId,
					playlist.Name,
				)
			}
			fmt.Println(display)
		},
	},
	{
		name: "add",
		args: []string{"name", "playlist id"},
		description: "Adds a playlist with name and youtube id",
		callback: func(args []string) {
			if len(args) < 2 {
				fmt.Println("Please, pass the playlist name and " +
					"youtube id as separeted values.")
				return
			}

			playlist := models.Playlist{ Name: args[0], YoutubeId: args[1], }

			_ = db.AddPlaylist(playlist)
			// playlists = append(playlists, playlist)
		},
	},
	{
		name: "rm",
		args: []string{"playlist name"},
		description: "Remove a playlist by it name.",
		callback: func(args []string) {
			if len(args) < 1 {
				fmt.Println("Please, pass the youtube id")
				return
			}

			playlistName := args[0]
			playlist := db.GetPlaylistByName(playlistName)
			db.RemovePlaylist(playlist)
		},
	},
	{
		name: "show",
		args: []string{"name"},
		description: "Show the playlist content given a name.",
		callback: func(args []string) {
			if len(args) < 1 {
				fmt.Println("Please, pass the playlist name." +
					"You can see it with `list-playlists` command.")
				return
			}

			playlistName := args[0]
			playlist := db.GetPlaylistByName(playlistName)

			videos := db.GetVideosByPlaylist(playlist.YoutubeId)

			lenBody := len(videos)
			fmt.Printf("Showing %v (%v):\n", playlist.Name, playlist.YoutubeId)
			for i, item := range videos {
				var watchedMark string
				if item.Watched {
					watchedMark += "*"
				} else {
					watchedMark += " "
				}

				fmt.Printf(
					"\t %vEpisode %03d: %v\n",
					watchedMark,
					lenBody - i,
					item.Name,
				)
			}
		},
	},
	{
		name: "watch",
		args: []string{"playlist", "video title"},
		description: "Mark a video as watched by it id.",
		callback: func(args []string) {
			if len(args) < 2 {
				fmt.Println("Please, provide the playlist and video name.")
				return
			}

			playlistName := args[0]
			playlist := db.GetPlaylistByName(playlistName)
			videoName := args[1]
			db.MarkVideoAsWatched(*playlist.Id, videoName)
		},
	},
	{
		name: "watch-all",
		args: []string{"playlist"},
		description: "Mark all of the playlist videos as watched.",
		callback: func(args []string) {
			if len(args) < 1 {
				fmt.Println("Please, provide the playlist name.")
			}

			playlistName := args[0]
			playlist := db.GetPlaylistByName(playlistName)
			db.MarkAllPlaylistVideosAsWatched(*playlist.Id)
		},
	},
}
