package utils

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
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

///Twitter
func TwitterFormatForAudio(twit anaconda.Tweet) string {
	var str string

	if twit.Lang == "en" {
		str = "kitten "
	} else {
		str = "котик "
	}
	str = str + twit.User.ScreenName + ". " + strings.Replace(twit.Text, "\n", "\\n", -1)

	re := regexp.MustCompile("http[s]?:\\/\\/t\\.co\\/.*?([ ]|$)")
	str = re.ReplaceAllString(str, "")
	str = strings.Replace(str, "/", "", -1)

	return str
}

func TwitterFormatForText(twit anaconda.Tweet) string {
	return "@" + twit.User.ScreenName + ": " + twit.Text
}

///ClosingBuffer
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}
