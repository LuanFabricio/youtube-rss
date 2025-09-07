package utils

import (
	"log"
	"os"
	"path/filepath"
)

func LogError(err error) {
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
}

func LogWarning(err error) {
	if err != nil {
		log.Printf("WARNING: %v\n", err)
	}
}

func GetBaseFolder() string {
	user_folder, err := os.UserHomeDir()
	LogError(err)

	return filepath.Join(user_folder, ".youtube-rss")
}

func CreateBaseFolderIfNotExists() {
	modes := os.ModeDir | os.ModePerm
	err := os.Mkdir(GetBaseFolder(), modes)
	if !os.IsExist(err) {
		LogError(err)
	}
}
