package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

/*
type commands interface {
	Command_Audio_List()
	Command_Audio_Stop()
	Command_Audio_Pause()
	Command_Audio_Resume()
	Command_Audio_Volume()

	Channels_List()
	Channels_Moveto()

	Help()
	Update()
	Status()
}*/

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
		"!audio_list                : cписок звуков<br/>" +
		"!audio_pause               : приостановить воспроизведение<br/>" +
		"!audio_resume              : восстановить воспроизведение<br/>" +
		"!audio_stop                : остановить воспроизведение звука<br/>" +
		"!audio_volume [float]      : установить громкость. Максимальная 100, cтандартная 50, минимальная 0, шаг 1<br/>" +

		"!channels_list             : список каналов<br/>" +
		"!channels_moveto [id/name] : перенести бота на другой канал<br/>" +

		"$$$[text]                  : произнести текст<br/>" +
		"#[sound]                   : произнести звук<br/>" +

		"!update                    : делает апдейт<br/>" +
		"!status                    : информация про бота<br/>" +

		"!help                      : эта команда<br/>"

		//"!disconnect       : отключить бота<br/>" +

	return str
}

func Command_Update() {
	if runtime.GOOS == "linux" {
		args := []string{"arg1"}
		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		os.StartProcess("./update_linux.sh", args, procAttr)
		//k.Client.Disconnect()
	} else {
		fmt.Println(Timestamp() + "!update works only in production")
	}
}
