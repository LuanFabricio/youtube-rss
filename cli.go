package main

import (
	"fmt"
)


type Callback func(args []string)

type Option struct {
	name string
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
		description: "Adds a playlist by it name and youtube id",
		callback: func(args []string) {
			if len(args) <= 1 {
				fmt.Println("Please, pass the playlist name and " +
					"youtube id as separeted values.")
				return
			}

			playlist := Playlist{ name: args[0], youtubeId: args[1], }

			playlists = append(playlists, playlist)
		},
	},
}

func registerOptions() map[string]Option {
	optionsMap := make(map[string]Option)

	for _, option := range options {
		optionsMap[option.name] = option
	}

	return optionsMap
}

func run(args []string) {
	optionsMap := registerOptions()

	if option, found := optionsMap[args[0]]; found {
		args := args[1:]
		option.callback(args)
		return
	}

	help := "Please, use one of the following options:\n\n"
	for _, option := range options {
		help += fmt.Sprintf("%v - %v\n", option.name, option.description)
	}
	fmt.Println(help)
}
