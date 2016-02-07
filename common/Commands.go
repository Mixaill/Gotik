package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"

	"../services"
	"../utils"
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

	Command_Twitter_ReadTwits(twits []anaconda.Tweet)
	Command_Twitter_Status(user string)

	Get_ConnectTime() time.Time
	Get_Services() *services.Services
	Get_Volume() float32

	Info_Name() string
}

func Command_Audio_List(backend Backend_interface) string {
	files, _ := ioutil.ReadDir("./sounds/")
	var sounds string = "Доступные звуки:<br/>"
	for _, f := range files {
		if f.IsDir() == false {
			var filename = f.Name()
			var extension = filepath.Ext(filename)
			var name = filename[0 : len(filename)-len(extension)]
			sounds += "<br/>" + name
		}
	}

	if backend.Info_Name() == "discord" {
		sounds = strings.Replace(sounds, "<br/>", "\n", -1)
	}

	return sounds
}

func Command_Help(backend Backend_interface) string {
	//audio
	str := "<br/>" +
		"$$$[text]                  : произнести текст<br/>" +
		"#[sound]                   : произнести звук<br/><br/>" +

		"!audio_list                : cписок звуков<br/>" +
		"!audio_play_file [sound]   : произнести звук<br/>" +
		"!audio_play_ivona [text]   : произнести текст<br/>" +
		"!audio_stop                : остановить воспроизведение звука<br/>"

	if backend.Info_Name() == "mumble" {
		str = str +
			"!audio_pause               : приостановить воспроизведение<br/>" +
			"!audio_resume              : восстановить воспроизведение<br/>" +
			"!audio_volume [float]      : установить громкость. Максимальная 100, cтандартная 50, минимальная 0, шаг 1<br/>"
	}

	//channels
	str = str +
		"<br/>!channels_list             : список каналов<br/>" +
		"!channels_moveto [id/name] : перенести бота на другой канал<br/>"

	//twitter
	str = str + "<br/>!twitter_status : информация о твиттере<br/><br/>"

	//other
	str = str + "<br/>!help                      : эта команда<br/><br/>"

	if false == true {
		str = str +
			"!disconnect                : отключить бота<br/>" +
			"!status                    : информация про бота<br/>" +
			"!update                    : делает апдейт<br/>"
	}

	if backend.Info_Name() == "discord" {
		str = strings.Replace(str, "<br/>", "\n", -1)
	}

	return str
}

func Command_Update() {
	if runtime.GOOS == "linux" {
		args := []string{"arg1"}
		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		os.StartProcess("./update_linux.sh", args, procAttr)
	} else {
		fmt.Println(utils.Timestamp() + "!update works only in production")
	}
}

func Command_Status(backend Backend_interface) string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	var str string = ""
	str = "<br/>" +
		"Backend:              : " + backend.Info_Name() + " <br/>" +
		"Uptime                : " + strconv.FormatFloat(time.Since(backend.Get_ConnectTime()).Hours(), 'f', 2, 64) + " hours <br/>"
	if backend.Info_Name() == "mumble" {

	}

	if backend.Info_Name() == "discord" {
		str = strings.Replace(str, "<br/>", "\n", -1)
	}

	return str
}

func Command_Process(message string, user string, backend Backend_interface) {

	message = strings.Replace(message, "<br/>", "", -1)
	message = strings.Replace(message, "</p>", "", -1)
	message = strings.Replace(message, "<p>", "", -1)

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

	case "!twitter_status":
		go backend.Command_Twitter_Status(user)

	case "!help":
		go backend.Command_Help(user)

	case "!disconnect":
		go backend.Command_Disconnect()
	case "!status":
		go backend.Command_Status(user)
	case "!update":
		go backend.Command_Update()

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
