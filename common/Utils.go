package common

import (
	"os"
	"strings"
	"time"
)

func Timestamp() string {
	return time.Now().UTC().Format("2006.01.02 15:04:05") + ": "
}

func GetAudioFilePath(text string) string {
	var filename string
	if text[0] == '#' {
		filename = strings.SplitN(text, "#", 2)[1]
	} else {
		filename = strings.SplitN(text, " ", 2)[1]
	}

	var formats = []string{".ogg", ".mp3", ".wav"}

	for _, format := range formats {
		if _, err := os.Stat("./sounds/" + filename + format); err == nil {
			return "./sounds/" + filename + format
		}
	}
	return ""
}
