package env

import (
	"io"
	"os"
	"strings"

	"youtube-rss/utils"
)

const DOTENV_FILE string = ".env"

func LoadEnv() {
	file, err := os.Open(DOTENV_FILE)
	utils.LogError(err)

	content, err := io.ReadAll(file)
	utils.LogError(err)

	lines := strings.Split(string(content), "\n")
	lines = lines[:len(lines) - 1]

	for _, line := range lines {
		lenLine := len(line)
		i := 0
		key := ""
		for line[i] != '=' && i < lenLine {
			key += string(line[i])
			i++
		}
		i++

		value := ""
		for i < lenLine  {
			value += string(line[i])
			i++
		}
		os.Setenv(key, value)
	}
}
