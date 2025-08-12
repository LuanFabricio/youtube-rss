package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const DOTENV_FILE string = ".env"

func LoadEnv() {
	file, err := os.Open(DOTENV_FILE)
	LogError(err)

	content, err := io.ReadAll(file)
	LogError(err)

	lines := strings.Split(string(content), "\n")
	lines = lines[:len(lines) - 1]
	fmt.Println(lines)

	for _, line := range lines {
		lenLine := len(line)
		i := 0
		key := ""
		for line[i] != '=' && i < lenLine {
			fmt.Printf("[%v]: %v\n", i, string(line[i]))
			key += string(line[i])
			i++
		}
		i++

		value := ""
		for i < lenLine  {
			fmt.Printf("[%v]: %v\n", i, string(line[i]))
			value += string(line[i])
			i++
		}
		os.Setenv(key, value)
	}
}
