package main

import (
	"log"
	"os"

	"youtube-rss/cli"
	db "youtube-rss/database"
	"youtube-rss/env"
	"youtube-rss/utils"
)

func main() {
	utils.CreateBaseFolderIfNotExists()
	db.InitDatabase()

	env.LoadEnv()

	key := os.Getenv("KEY")
	if key == "" {
		log.Panicln("Invalid KEY")
	}

	cli.Run(os.Args)
}
