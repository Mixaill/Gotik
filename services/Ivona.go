package services

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"../common"
	"../config"

	"github.com/jpadilla/ivona-go"
)

type Ivona_voice struct {
	Rate           string
	Volume         string
	SentenceBreak  int
	ParagraphBreak int

	Name     string
	Language string
	Gender   string
}

type Ivona struct {
	iv     *ivona.Ivona
	voices map[string]config.Config_voice
}

func NewIvona() *Ivona {
	p := new(Ivona)
	p.iv = ivona.New(config.Ivona_key_1, config.Ivona_key_2)
	p.voices = config.Ivona_Voices()
	return p
}

func (i *Ivona) GetAudio_Response(text string, language string) (*ivona.SpeechResponse, error) {
	var val config.Config_voice
	var ok bool

	if val, ok = i.voices[language]; ok {
	} else {
		val = i.voices[config.Ivona_Fallback]
	}

	opts_input := ivona.Input{Data: text, Type: "text/plain"}
	opts_outputformat := ivona.OutputFormat{Codec: string(config.Ivona_Codec), SampleRate: int(config.Ivona_SampleRate)}
	opts_parameters := ivona.Parameters{Rate: val.Rate, Volume: val.Volume, SentenceBreak: val.SentenceBreak, ParagraphBreak: val.ParagraphBreak}
	opts_voice := ivona.Voice{Name: val.Name, Language: val.Language, Gender: val.Gender}

	speechopts := ivona.SpeechOptions{Input: &opts_input, OutputFormat: &opts_outputformat, Parameters: &opts_parameters, Voice: &opts_voice}

	respond, err := i.iv.CreateSpeech(speechopts)

	if err != nil {
		fmt.Println("services/ivona/GetAudio_Response(): error ", err)
	}

	return respond, err
}

func (i *Ivona) GetAudio_ReadCloser(text string, language string) io.ReadCloser {
	respond, err := i.GetAudio_Response(text, language)

	var rc io.ReadCloser
	if err == nil {
		rc = &common.ClosingBuffer{bytes.NewBuffer(respond.Audio)}
	}
	return rc
}

func (i *Ivona) GetAudio_File(text string, language string) string {
	respond, err := i.GetAudio_Response(text, language)

	if err == nil {
		os.MkdirAll("./temp", os.ModePerm)
		ioutil.WriteFile("./temp/ivona_discord.mp3", respond.Audio, os.ModePerm)
		return "./temp/ivona_discord.mp3"
	}
	return ""
}
