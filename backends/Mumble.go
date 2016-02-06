package backends

import (
	"crypto/tls"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	_ "github.com/layeh/gumble/opus"

	"../common"
	"../services"
)

type Mumble struct {
	Audio  *gumbleffmpeg.Stream
	Config *gumble.Config
	Client *gumble.Client

	services *services.Services

	connectTime time.Time
	conf_volume float32
}

func NewMumble() *Mumble {
	p := new(Mumble)
	return p
}

func (k *Mumble) Start(fl map[string]string, s *services.Services) {
	k.services = s

	//Config
	k.Config = gumble.NewConfig()
	k.Config.Username = fl["Flag_mumble_username"]
	k.Config.Password = fl["Flag_mumble_password"]
	k.Config.Address = fl["Flag_mumble_server"]

	k.Config.TLSConfig.InsecureSkipVerify = true

	k.conf_volume = 1.0

	//Client creation
	k.Client = gumble.NewClient(k.Config)

	//Attach listeners
	k.Client.Attach(gumbleutil.AutoBitrate)
	k.Client.Attach(k)

	//TLS
	if fl["Flag_mumble_cert_lock"] != "" {
		gumbleutil.CertificateLockFile(k.Client, fl["Flag_mumble_cert_lock"])
	}

	if _, err := os.Stat("./config/Kotik.pem"); err == nil && fl["Flag_dev"] == "false" {
		fl["Flag_mumble_cert"] = "./config/Kotik.pem"
	} else if _, err := os.Stat("./config/Kotik-dev.pem"); err == nil && fl["Flag_dev"] == "true" {
		fl["Flag_mumble_cert"] = "./config/Kotik-dev.pem"
	}

	if fl["Flag_mumble_cert"] != "" {
		if fl["Flag_mumble_cert_key"] == "" {
			fl["Flag_mumble_cert_key"] = fl["Flag_mumble_cert"]
		}
		if certificate, err := tls.LoadX509KeyPair(fl["Flag_mumble_cert"], fl["Flag_mumble_cert_key"]); err != nil {
			panic(err)
		} else {
			k.Config.TLSConfig.Certificates = append(k.Config.TLSConfig.Certificates, certificate)
		}
	}

	//Connect
	if err := k.Client.Connect(); err != nil {
		panic(err)
	}
}

/////
/////Listeners
/////
func (k *Mumble) OnConnect(e *gumble.ConnectEvent) {
	k.connectTime = time.Now()
	fmt.Println(common.Timestamp() + "OnConnect()")
}

func (k *Mumble) OnDisconnect(e *gumble.DisconnectEvent) {
	fmt.Println(common.Timestamp() + "OnDisconnect()")
}

func (k *Mumble) OnUserChange(e *gumble.UserChangeEvent) {
	fmt.Println(common.Timestamp() + "OnUserChange()-> User:" + e.User.Name + " | Type:" + fmt.Sprint(e.Type))
}

func (k *Mumble) OnChannelChange(e *gumble.ChannelChangeEvent) {
	fmt.Println(common.Timestamp() + "OnChannelChange()-> Channel:" + e.Channel.Name + " | Type: " + fmt.Sprint(e.Type))
}

func (k *Mumble) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	fmt.Println(common.Timestamp() + "OnPermissionDenied()-> User:" + e.User.Name + " | Type: " + fmt.Sprint(e.Type))
}

func (k *Mumble) OnUserList(e *gumble.UserListEvent) {
	fmt.Println(common.Timestamp() + "OnUserList()")
}

func (k *Mumble) OnACL(e *gumble.ACLEvent) {
	fmt.Println(common.Timestamp() + "OnACL()")
}

func (k *Mumble) OnBanList(e *gumble.BanListEvent) {
	fmt.Println(common.Timestamp() + "OnBanList()")
}

func (k *Mumble) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
	fmt.Println(common.Timestamp() + "OnContextActionChange()")
}

func (k *Mumble) OnServerConfig(e *gumble.ServerConfigEvent) {
	fmt.Println(common.Timestamp() + "OnServerConfig()")
}

func (k *Mumble) OnTextMessage(e *gumble.TextMessageEvent) {
	fmt.Println(common.Timestamp() + "OnTextMessage()-> Message:" + e.Message)
	if e.Sender != nil {
		fmt.Println(common.Timestamp() + "               -> User: " + e.Sender.Name)
	} else {
		return
	}

	if e.Sender.IsRegistered() == false {
		return
	}

	re_cmd := regexp.MustCompile("^!\\w+")
	switch re_cmd.FindString(e.Message) {
	case "!audio_list":
		k.command_audio_list(e.Sender)
	case "!audio_pause":
		k.command_audio_pause()
	case "!audio_resume":
		k.command_audio_resume()
	case "!audio_stop":
		k.command_audio_stop()
	case "!audio_volume":
		k.command_audio_volume(e.Message)

	case "!channels_list":
		k.command_channels_list(e.Sender)
	case "!channels_moveto":
		k.command_channels_moveto(e.Message)

	case "!help":
		k.command_help(e.Sender)

	case "!update":
		k.command_update()
	case "!disconnect":
		k.command_disconnect(e)
	case "!status":
		k.command_status(e.Sender)
		//case "!twitter":
		//	go k.command_twitter_process(e.Sender)
	}

	re_snd := regexp.MustCompile("#(\\w+)")
	result_snd := re_snd.FindStringSubmatch(e.Message)
	if len(result_snd) == 2 {
		switch result_snd[1] {
		case "ymusic":
			go k.command_play_ymusic(e.Message, e.Sender)
		default:
			go k.command_audio_play_file(e.Message)
		}
	}

	//ivona
	re_ivona := regexp.MustCompile("\\$\\$\\$(.*)")
	result_ivona := re_ivona.FindStringSubmatch(e.Message)
	if len(result_ivona) == 2 {
		go k.command_audio_play_ivona(result_ivona[1])
	}

}

/////
/////Commands
/////

/////Commands/Audio
func (k *Mumble) command_audio_list(e *gumble.User) {
	e.Send(common.Command_Audio_List())
}

func (k *Mumble) command_audio_resume() {
	if k.Audio != nil {
		k.Audio.Play()
	}
}

func (k *Mumble) command_audio_pause() {
	if k.Audio != nil && k.Audio.State() != gumbleffmpeg.StatePaused {
		k.Audio.Pause()
	}
}

func (k *Mumble) command_audio_stop() {
	if k.Audio != nil {
		k.Audio.Stop()
	}
}

func (k *Mumble) command_audio_play_file(text string) {
	if k.Audio != nil && k.Audio.State() == gumbleffmpeg.StatePlaying {
		return
	}
	fmt.Println(common.Timestamp() + "backends/mumble: command_audio_play_file(): " + text)

	filename := strings.Split(text, "#")[1]
	var formats = []string{".ogg", ".mp3", ".wav"}

	for _, format := range formats {
		if _, err := os.Stat("./sounds/" + filename + format); err == nil {
			k.Audio = gumbleffmpeg.New(k.Client, gumbleffmpeg.SourceFile("./sounds/"+filename+format))
			k.Audio.Volume = k.conf_volume
			k.Audio.Play()
			k.Audio.Wait()
		}
	}
}

func (k *Mumble) command_audio_play_ivona(text string) {
	if k.Audio != nil && k.Audio.State() == gumbleffmpeg.StatePlaying {
		return
	}
	fmt.Println(common.Timestamp() + "backends/mumble: command_audio_play_ivona(): " + text)

	rc := k.services.Ivona.GetAudio_ReadCloser(text, "ru")
	k.Audio = gumbleffmpeg.New(k.Client, gumbleffmpeg.SourceReader(rc))
	k.Audio.Volume = k.conf_volume
	k.Audio.Play()
}

func (k *Mumble) command_audio_volume(text string) {
	re_sound := regexp.MustCompile("^!volume[ ]?(\\d+)")
	result_sound := re_sound.FindStringSubmatch(text)
	if len(result_sound) == 2 {
		i, err := strconv.Atoi(result_sound[1])
		if err == nil {
			if i > 100 {
				i = 100
			}
			k.conf_volume = float32(i) / 50.00
			if k.Audio != nil {
				k.Audio.Volume = k.conf_volume
			}
		}
	}
}

/////Commands/Channels

func (k *Mumble) command_channels_list(e *gumble.User) {
	root := k.Client.Channels[0]
	var channels string
	if root != nil {
		channels += "<br/>" + fmt.Sprint(root.ID) + ": " + root.Name + "<br/>"
		if root.Children != nil {
			channels += k.command_channels_list_printchild(root.Children, 1)
		}
	}
	e.Send(channels)
}

func (k *Mumble) command_channels_list_printchild(children gumble.Channels, level int) string {
	var out = ""
	var keys []int
	for k := range children {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, j := range keys {
		for i := 0; i < level; i++ {
			out += "----"
		}
		out += fmt.Sprint(children[uint32(j)].ID) + ": " + children[uint32(j)].Name + "<br/>"
		if children[uint32(j)].Children != nil {
			out += k.command_channels_list_printchild(children[uint32(j)].Children, level+1)
		}
	}
	return out
}

func (k *Mumble) command_channels_moveto(text string) {
	re_id := regexp.MustCompile("^!moveto ([0-9]+)")
	var result_id []string = re_id.FindStringSubmatch(text)
	if len(result_id) == 2 {
		i, err := strconv.Atoi(result_id[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		if k.Client.Channels[uint32(i)] != nil {
			k.Client.Self.Move(k.Client.Channels[uint32(i)])
		}
	} else {
		re_str := regexp.MustCompile("^!moveto (.*)")
		var result_str []string = re_str.FindStringSubmatch(text)

		if len(result_str) == 2 {
			var channel *gumble.Channel

			if result_str[1] == "Root" {
				channel = k.Client.Channels[0]
			} else {
				channel = k.Client.Channels.Find(result_str[1])
			}

			if channel != nil {
				k.Client.Self.Move(channel)
			}
		}
	}
}

/////

func (k *Mumble) command_status(e *gumble.User) {

}

func (k *Mumble) command_update() {
	common.Command_Update()
}

func (k *Mumble) command_help(e *gumble.User) {
	e.Send(common.Command_Help())
}

func (k *Mumble) command_disconnect(e *gumble.TextMessageEvent) {
	//k.Client.Disconnect()
}

func (k *Mumble) command_play_ymusic(text string, e *gumble.User) {
	if k.Audio != nil && k.Audio.State() == gumbleffmpeg.StatePlaying {
		return
	}

	ym := services.YMusic{}
	trackname := strings.Split(text, "#ymusic ")[1]
	file, title := ym.GetTrack(trackname)
	if title != "" {
		e.Send("Найдена композиция: " + title)
	} else {
		e.Send("Композиция не найдена")
	}
	if file != nil {
		k.Audio = gumbleffmpeg.New(k.Client, gumbleffmpeg.SourceReader(file))
		k.Audio.Volume = k.conf_volume
		k.Audio.Play()
	}
}
