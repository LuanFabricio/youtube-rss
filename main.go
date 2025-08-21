package main

import (
	"log"
	"os"

	"youtube-rss/cli"
	db "youtube-rss/database"
)

func main() {

	// run([]string{"add", "the office"})
	// run([]string{"add", "the office", "idk2"})
	// run([]string{"list-playlists"})
	// run([]string{""})
	// run([]string{"rm", "idk2"})
	// run([]string{"rm", "1"})
	// run([]string{"rm"})

	db.InitDatabase()

	LoadEnv()

	key := os.Getenv("KEY")
	if key == "" {
		log.Panicln("Invalid KEY")
	}

	cli.Run(os.Args)
}
