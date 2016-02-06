package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type Backend_interface interface {
	Command_Audio_List(user string)
	Command_Audio_Stop()
	Command_Audio_Pause()
	Command_Audio_Play_File(text string)
	Command_Audio_Play_Ivona(text string, language string)
	Command_Audio_Resume()
	Command_Audio_Volume(text string)

	Command_Channels_List(user string)
	Command_Channels_Moveto(text string)

	Command_Help(user string)

	Command_Disconnect()
	Command_Status(user string)
	Command_Update()
}

func Command_Audio_List() string {
	files, _ := ioutil.ReadDir("./sounds/")
	var sounds string = ""
	for _, f := range files {
		if f.IsDir() == false {
			var filename = f.Name()
			var extension = filepath.Ext(filename)
			var name = filename[0 : len(filename)-len(extension)]
			sounds += "<br/>" + name
		}
	}
	return sounds
}

func Command_Help() string {
	str := "<br/>" +
		"$$$[text]                  : произнести текст<br/>" +
		"#[sound]                   : произнести звук<br/><br/>" +

		"!audio_list                : cписок звуков<br/>" +
		"!audio_pause               : приостановить воспроизведение<br/>" +
		"!audio_play_file [sound]   : произнести звук<br/>" +
		"!audio_play_ivona [text]   : произнести текст<br/>" +
		"!audio_resume              : восстановить воспроизведение<br/>" +
		"!audio_stop                : остановить воспроизведение звука<br/>" +
		"!audio_volume [float]      : установить громкость. Максимальная 100, cтандартная 50, минимальная 0, шаг 1<br/><br/>" +

		"!channels_list             : список каналов<br/>" +
		"!channels_moveto [id/name] : перенести бота на другой канал<br/><br/>" +

		"!help                      : эта команда<br/><br/>" +

		"!disconnect                : отключить бота<br/>" +
		"!status                    : информация про бота<br/>" +
		"!update                    : делает апдейт<br/>"

	return str
}

func Command_Update() {
	if runtime.GOOS == "linux" {
		args := []string{"arg1"}
		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		os.StartProcess("./update_linux.sh", args, procAttr)
	} else {
		fmt.Println(Timestamp() + "!update works only in production")
	}
}

func Command_Process(message string, user string, backend Backend_interface) {

	message = strings.TrimSuffix(message, "</p>")
	message = strings.TrimPrefix(message, "<p>")

	re_cmd := regexp.MustCompile("^!\\w+")
	switch re_cmd.FindString(message) {
	case "!audio_list":
		go backend.Command_Audio_List(user)
	case "!audio_pause":
		go backend.Command_Audio_Pause()
	case "!audio_play_file":
		go backend.Command_Audio_Play_File(message)
	case "!audio_play_ivona":
		go backend.Command_Audio_Play_Ivona(message, "ru")
	case "!audio_resume":
		go backend.Command_Audio_Resume()
	case "!audio_stop":
		go backend.Command_Audio_Stop()
	case "!audio_volume":
		go backend.Command_Audio_Volume(message)

	case "!channels_list":
		go backend.Command_Channels_List(user)
	case "!channels_moveto":
		go backend.Command_Channels_Moveto(message)

	case "!help":
		go backend.Command_Help(user)

	case "!disconnect":
		go backend.Command_Disconnect()
	case "!status":
		go backend.Command_Status(user)
	case "!update":
		go backend.Command_Update()

		//case "!twitter":
		//	go k.command_twitter_process(e.Sender)
	}

	// alias: #
	re_snd := regexp.MustCompile("#(\\w+)")
	result_snd := re_snd.FindStringSubmatch(message)
	if len(result_snd) == 2 {
		switch result_snd[1] {
		case "ymusic":
		//	go k.command_play_ymusic(e.Message, e.Sender)
		default:
			go backend.Command_Audio_Play_File(message)
		}
	}

	// alias: $$$
	re_ivona := regexp.MustCompile("\\$\\$\\$(.*)")
	result_ivona := re_ivona.FindStringSubmatch(message)
	if len(result_ivona) == 2 {
		go backend.Command_Audio_Play_Ivona(message, "ru")
	}

}
