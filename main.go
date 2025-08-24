package main

import (
	"log"
	"os"

	db "youtube-rss/database"
	"youtube-rss/cli"
	"youtube-rss/env"
)

func main() {
	db.InitDatabase()

	env.LoadEnv()

	key := os.Getenv("KEY")
	if key == "" {
		log.Panicln("Invalid KEY")
	}

	cli.Run(os.Args)
}
