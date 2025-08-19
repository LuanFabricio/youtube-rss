package main

import (
	"fmt"
	"os"
)


type Callback func(args []string)

type Option struct {
	name string
	args []string
	description string
	callback Callback
}

type Playlist struct {
	id *int
	name string
	youtubeId string
}

var playlists []Playlist = []Playlist{
	{ name: "The Standup", youtubeId: "PL2Fq-K0QdOQiJpufsnhEd1z3xOv2JMHuk"},
}

var options []Option = []Option {
	{
		name: "list-playlists",
		args: []string{},
		description: "List the playlists.",
		callback: func(args []string) {
			playlists = GetPlaylists()
			display := ""
			for _, playlist := range playlists {
				display += fmt.Sprintf(
					"- [%v]: %v\n",
					playlist.youtubeId,
					playlist.name,
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

			playlist := Playlist{ name: args[0], youtubeId: args[1], }

			_ = AddPlaylist(playlist)
			// playlists = append(playlists, playlist)
		},
	},
	{
		name: "rm",
		args: []string{"playlist id"},
		description: "Remove a playlist by youtube id.",
		callback: func(args []string) {
			if len(args) < 1 {
				fmt.Println("Please, pass the youtube id")
				return
			}

			playlistId := args[0]
			if !removePlaylistById(playlistId) {
				fmt.Printf("Cannot found playlist \"%v\".\n", playlistId)
			}
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
			index := -1
			playlists = GetPlaylists()
			for i := range playlists {
				if playlists[i].name == playlistName {
					index = i
					break
				}
			}

			if index == -1 {
				fmt.Println("Playlist name not found." +
					"You can list the playlists with" +
					" `list-playlists` command.")
				return
			}
			playlist := playlists[index]

			body, err := GetPlaylist(os.Getenv("KEY"), playlist.youtubeId)
			LogError(err)

			lenBody := len(body.Items)
			fmt.Printf("Showing %v (%v):\n", playlist.name, playlist.youtubeId)
			for i, item := range body.Items {
				fmt.Printf(
					"\tEpisode %03d: %v\n",
					lenBody - i,
					item.Snippet.Title,
				)
			}
		},
	},
}

func removePlaylistById(youtubeId string) bool {
	index := -1

	for i, playlist := range playlists {
		if playlist.youtubeId == youtubeId {
			index = i
			break
		}
	}

	if index > -1 {
		playlists = append(playlists[:index], playlists[index+1:]...)
		return true
	}

	return false
}

func registerOptions() map[string]Option {
	optionsMap := make(map[string]Option)

	for _, option := range options {
		optionsMap[option.name] = option
	}

	return optionsMap
}

func help() {
	help := "Please, use one of the following options:\n\n"
	for _, option := range options {
		args := ""
		fmt.Println(option.args)
		for _, arg := range option.args {
			args += fmt.Sprintf(" <%v>", arg)
		}
		help += fmt.Sprintf("%v%v - %v\n", option.name, args, option.description)
	}
	fmt.Println(help)

}

func updateVideos() {
	playlists := GetPlaylists()
	for i, playlist := range playlists {
		fmt.Printf("[%v] Updating playlist %v...\n", i + 1, playlist.name)
		AddPlaylistVideosIfNotRegistered(playlist)
	}
}

func run(args []string) {
	updateVideos()

	args = args[1:]
	if len(args) == 0 {
		help()
		return
	}

	optionsMap := registerOptions()
	option, found := optionsMap[args[0]]
	if found {
		args := args[1:]
		option.callback(args)
		return
	}

	help()
}
