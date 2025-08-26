package cli

import (
	"os"
	db "youtube-rss/database"
)

type FlagCallback func(args []string) []string

type Flag struct {
	name string
	shortname string
	description string
	callback FlagCallback
}

var flags []Flag = []Flag{
	{
		name: "just-update",
		shortname: "ju",
		description: "Just run the update function.",
		callback: func(args []string) []string {
			os.Exit(0)
			return []string{}
		},
	},
	{
		name: "show-all-unwatched",
		shortname: "sau",
		description: "List all of the unwatched episodes from all of the playlists.",
		callback: func(args []string) []string {
			queries := []string{}

			playlists := db.GetPlaylists()

			for _, playlist := range playlists {
				queries = append(queries, "show-unwatched", playlist.Name)
			}

			return queries
		},
	},
}

func registerFlags() map[string]Flag {
	flagsMap := make(map[string]Flag)

	for _, flag := range flags {
		flagsMap["--"+flag.name] = flag
		flagsMap["-"+flag.shortname] = flag
	}

	return flagsMap
}

func applyFlags(args []string) []string {
	flagsMap := registerFlags()

	for i := range args {
		flag, found := flagsMap[args[i]]
		if found {
			args = flag.callback(args)
		}
	}

	return args
}
