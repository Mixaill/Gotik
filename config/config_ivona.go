package config

type Config_voice struct {
	Rate           string
	Volume         string
	SentenceBreak  int
	ParagraphBreak int

	Name     string
	Language string
	Gender   string
}

const (
	Ivona_key_1 = "GDNAIKHZ6EJPBXXTKZFA"
	Ivona_key_2 = "akU4WnCw2XozeJeMsnS7pVqlBsLgqF4FQbVRnu05"

	Ivona_Codec      = "MP3"
	Ivona_SampleRate = 22050
	Ivona_Fallback   = "ru"
)

func Ivona_Voices() map[string]Config_voice {

	Ivona_Voices := make(map[string]Config_voice)
	Ivona_Voices["ru"] = Config_voice{Rate: "medium",
		Volume:         "medium",
		SentenceBreak:  100,
		ParagraphBreak: 500,
		Name:           "Maxim",
		Language:       "ru-RU",
		Gender:         "Male"}
	Ivona_Voices["en"] = Config_voice{Rate: "medium",
		Volume:         "medium",
		SentenceBreak:  100,
		ParagraphBreak: 500,
		Name:           "Joey",
		Language:       "en-US",
		Gender:         "Male"}
	return Ivona_Voices
}
