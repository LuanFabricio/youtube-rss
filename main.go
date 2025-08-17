package main

import (
	"log"
	"os"
)

func main() {

	// run([]string{"add", "the office"})
	// run([]string{"add", "the office", "idk2"})
	// run([]string{"list-playlists"})
	// run([]string{""})
	// run([]string{"rm", "idk2"})
	// run([]string{"rm", "1"})
	// run([]string{"rm"})

	LoadEnv()

	key := os.Getenv("KEY")
	if key == "" {
		log.Panicln("Invalid KEY")
	}

	run(os.Args)
}
