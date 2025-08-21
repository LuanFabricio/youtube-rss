package utils

import "log"

func LogError(err error) {
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
}
