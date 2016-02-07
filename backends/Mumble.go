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
}

func (k *Mumble) OnDisconnect(e *gumble.DisconnectEvent) {
}

func (k *Mumble) OnUserChange(e *gumble.UserChangeEvent) {
}

func (k *Mumble) OnChannelChange(e *gumble.ChannelChangeEvent) {
}

func (k *Mumble) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
}

func (k *Mumble) OnUserList(e *gumble.UserListEvent) {
}

func (k *Mumble) OnACL(e *gumble.ACLEvent) {
}

func (k *Mumble) OnBanList(e *gumble.BanListEvent) {
}

func (k *Mumble) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (k *Mumble) OnServerConfig(e *gumble.ServerConfigEvent) {
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

	common.Command_Process(e.Message, e.Sender.Name, k)
}

/////
/////Commands
/////

/////Commands/Audio
func (k *Mumble) Command_Audio_List(user string) {
	k.Client.Users.Find(user).Send(common.Command_Audio_List(k))
}

func (k *Mumble) Command_Audio_Resume() {
	if k.Audio != nil {
		k.Audio.Play()
	}
}

func (k *Mumble) Command_Audio_Pause() {
	if k.Audio != nil && k.Audio.State() != gumbleffmpeg.StatePaused {
		k.Audio.Pause()
	}
}

func (k *Mumble) Command_Audio_Stop() {
	if k.Audio != nil {
		k.Audio.Stop()
	}
}

func (k *Mumble) Command_Audio_Play_File(text string) {
	if k.Audio != nil && k.Audio.State() == gumbleffmpeg.StatePlaying {
		return
	}

	var filename string
	if text[0] == '#' {
		filename = strings.SplitN(text, "#", 2)[1]
	} else {
		filename = strings.SplitN(text, " ", 2)[1]
	}
	fmt.Println(common.Timestamp() + "backends/mumble: command_audio_play_file(): " + filename)

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

func (k *Mumble) Command_Audio_Play_Ivona(text string, language string) {
	if k.Audio != nil && k.Audio.State() == gumbleffmpeg.StatePlaying {
		return
	}

	var tt string
	if strings.HasPrefix(text, "$$$") {
		tt = strings.Split(text, "$$$")[1]
	} else {
		tt = strings.SplitN(text, " ", 2)[1]
	}
	fmt.Println(common.Timestamp() + "backends/mumble: command_audio_play_ivona(): " + tt)

	rc := k.services.Ivona.GetAudio_ReadCloser(tt, language)
	k.Audio = gumbleffmpeg.New(k.Client, gumbleffmpeg.SourceReader(rc))
	k.Audio.Volume = k.conf_volume
	k.Audio.Play()
}

func (k *Mumble) Command_Audio_Volume(text string) {
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

func (k *Mumble) Command_Channels_List(user string) {

	root := k.Client.Channels[0]
	var channels string
	if root != nil {
		channels += "<br/>" + fmt.Sprint(root.ID) + ": " + root.Name + "<br/>"
		if root.Children != nil {
			channels += k.command_channels_list_printchild(root.Children, 1)
		}
	}
	k.Client.Users.Find(user).Send(channels)
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

func (k *Mumble) Command_Channels_Moveto(text string) {
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

/////Commands/other
func (k *Mumble) Command_Help(user string) {
	k.Client.Users.Find(user).Send(common.Command_Help(k))
}

func (k *Mumble) Command_Update() {
	common.Command_Update()
}

func (k *Mumble) Command_Disconnect() {
}

func (k *Mumble) Command_Status(user string) {

}

/////
///// Info
/////

func (k *Mumble) Info_Name() string {
	return "mumble"
}
