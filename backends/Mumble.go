package backends

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"layeh.com/gumble/gumbleutil"
	_ "layeh.com/gumble/opus"

	"../common"
	"../services"
	"../utils"
)

type Mumble struct {
	Audio     *gumbleffmpeg.Stream
	Config    *gumble.Config
	TLSConfig *tls.Config
	Client    *gumble.Client

	services *services.Services

	conf_connectTime   time.Time
	conf_volume        float32
	conf_twitterEnable bool
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
	k.conf_volume = 1.0

	//Attach listeners
	k.Config.Attach(gumbleutil.AutoBitrate)
	k.Config.Attach(k)

	//TLS
	k.TLSConfig = new(tls.Config)
	k.TLSConfig.InsecureSkipVerify = true
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
			k.TLSConfig.Certificates = append(k.TLSConfig.Certificates, certificate)
		}
	}

	k.Client, _ = gumble.DialWithDialer(new(net.Dialer), fl["Flag_mumble_server"], k.Config, k.TLSConfig)
}

/////Listeners
/////
func (k *Mumble) OnConnect(e *gumble.ConnectEvent) {
	k.conf_connectTime = time.Now()
	k.conf_twitterEnable = true
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
	fmt.Println(utils.Timestamp() + "OnTextMessage()-> Message:" + e.Message)
	if e.Sender != nil {
		fmt.Println(utils.Timestamp() + "               -> User: " + e.Sender.Name)
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

	filepath := utils.GetAudioFilePath(text)
	if filepath != "" {
		k.Audio = gumbleffmpeg.New(k.Client, gumbleffmpeg.SourceFile(filepath))
		k.Audio.Volume = k.conf_volume
		k.Audio.Play()
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
	fmt.Println(utils.Timestamp() + "backends/mumble: command_audio_play_ivona(): " + language + ":" + tt)

	if language == "" {
		language = k.services.YTranslate.Detect(tt)
	}

	rc := k.services.Ivona.GetAudio_ReadCloser(tt, language)
	k.Audio = gumbleffmpeg.New(k.Client, gumbleffmpeg.SourceReader(rc))
	k.Audio.Volume = k.conf_volume
	k.Audio.Play()
}

func (k *Mumble) Command_Audio_Volume(text string) {
	re_sound := regexp.MustCompile("^!audio_volume[ ]?(\\d+)")
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
	re_id := regexp.MustCompile("^!channels_moveto ([0-9]+)")
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
		re_str := regexp.MustCompile("^!channels_moveto (.*)")
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
	var str string
	str = str +
		"Volume                : " + strconv.FormatInt(int64(k.Get_Volume()*50.00), 10) + "% <br/>"
	k.Client.Users.Find(user).Send(str)
}

/////Commands/twitter
func (k *Mumble) Command_Twitter_ReadTwits(twits []anaconda.Tweet) {
	if k.conf_twitterEnable == true {
		for _, val := range twits {
			k.Client.Self.Channel.Send(utils.TwitterFormatForText(val), false)
			k.Command_Audio_Play_Ivona(utils.TwitterFormatForAudio(val), val.Lang)
			k.Audio.Wait()
		}
	}
}

func (k *Mumble) Command_Twitter_Status(user string) {
	var str string
	str = str +
		"<br/>Twitter subscriptions : " + k.services.Twitter.UsersGet() + "<br/>" +
		"Twitter update rate   : " + strconv.FormatFloat(k.services.Twitter.UpdateRateGet().Minutes(), 'f', 2, 64) + " minutes <br/>"
	k.Client.Users.Find(user).Send(str)
}

func (k *Mumble) Command_Twitter_Switch() {
	k.conf_twitterEnable = !k.conf_twitterEnable
}

/////
///// Getters
/////
func (k *Mumble) Get_ConnectTime() time.Time {
	return k.conf_connectTime
}

func (k *Mumble) Get_Services() *services.Services {
	return k.services
}

func (k *Mumble) Get_Volume() float32 {
	return k.conf_volume
}

/////
///// Info
/////

func (k *Mumble) Info_Name() string {
	return "mumble"
}
