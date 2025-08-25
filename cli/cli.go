package cli

import (
	"fmt"

	db "youtube-rss/database"
	"youtube-rss/models"
)

type OptionCallback func(args []string)

type Option struct {
	name string
	args []string
	description string
	callback OptionCallback
}

func help() {
	help := "Please, use one of the following options:\n\n"
	for _, option := range options {
		args := ""
		for _, arg := range option.args {
			args += fmt.Sprintf(" <%v>", arg)
		}
		help += fmt.Sprintf("%v%v - %v\n", option.name, args, option.description)
	}
	fmt.Println(help)
}

func updatePlaylist(playlist models.Playlist, channel chan int) {
	videosNames := db.AddPlaylistVideosIfNotRegistered(playlist)

	for _, videoName := range videosNames {
		fmt.Printf("[%v] %v\n", playlist.Name, videoName)
	}

	channel <- 0
}

func updateVideos() {
	fmt.Println("Updating playlists...")
	playlists := db.GetPlaylists()
	channel := make(chan int)
	for _, playlist := range playlists {
		go updatePlaylist(playlist, channel)
	}

	for range playlists {
		<- channel
	}

	fmt.Println()
}

func Run(args []string) {
	updateVideos()

	args = args[1:]

	args = applyFlags(args)

	if len(args) == 0 {
		help()
		return
	}

	optionsMap := registerOptions()
	option, found := optionsMap[args[0]]

	for found && len(args) > 0 {
		args = args[1:]
		option.callback(args)

		args = args[len(option.args):]
		if len(args) == 0 {
			break
		}
		option, found = optionsMap[args[0]]
	}

	if !found {
		help()
	}
}
