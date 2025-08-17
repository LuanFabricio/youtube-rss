package main

import (
	"fmt"
)


type Callback func(args []string)

type Option struct {
	name string
	args []string
	description string
	callback Callback
}

type Playlist struct {
	name string
	youtubeId string

}

var playlists []Playlist = []Playlist{
	{ name: "The Standup", youtubeId: "idk"},
}

var options []Option = []Option {
	{
		name: "list-playlists",
		args: []string{},
		description: "List the playlists.",
		callback: func(args []string) {
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

			playlists = append(playlists, playlist)
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

func run(args []string) {
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
