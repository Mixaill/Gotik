package gotik

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/gumble/gumbleutil"
	_ "github.com/layeh/gumble/opus"

	"github.com/jpadilla/ivona-go"
)

//Kotik struct
type Kotik struct {
	Audio  *gumble_ffmpeg.Stream
	Config *gumble.Config
	Client *gumble.Client
	Ivona  *ivona.Ivona

	keepAlive   chan bool
	connectTime time.Time
}

//Helpers
func (k *Kotik) Timestamp() string {
	return time.Now().UTC().Format("2006.01.02 15:04:05") + ": "
}

//Listeners
func (k *Kotik) OnConnect(e *gumble.ConnectEvent) {
	k.connectTime = time.Now()
	fmt.Println(k.Timestamp() + "OnConnect()")
}

func (k *Kotik) OnDisconnect(e *gumble.DisconnectEvent) {
	fmt.Println(k.Timestamp() + "OnDisconnect()")
	os.Exit(0)
}

func (k *Kotik) OnTextMessage(e *gumble.TextMessageEvent) {
	fmt.Println(k.Timestamp() + "OnTextMessage()-> User: " + e.Sender.Name + "| Message:" + e.Message)

	re_cmd := regexp.MustCompile("^!\\w+")
	switch re_cmd.FindString(e.Message) {
	case "!help":
		k.command_help(e.Sender)
	case "!list_channels":
		k.command_list_channels(e.Sender)
	case "!list_sounds":
		k.command_list_sounds(e.Sender)
	case "!moveto":
		k.command_moveto(e.Message)
	case "!stop":
		k.command_stop()
	case "!volume":
		k.command_volume(e.Message)
	case "!update":
		k.command_update(e.Sender)
	case "!disconnect":
		k.command_disconnect(e)
	case "!pause":
		k.command_pause()
	case "!resume":
		k.command_resume()
	case "!status":
		k.command_status(e.Sender)
	}

	re_snd := regexp.MustCompile("#(\\w+)")
	result_snd := re_snd.FindStringSubmatch(e.Message)
	if len(result_snd) == 2 {
		switch result_snd[1] {
		case "ymusic":
			k.command_play_ymusic(e.Message, e.Sender)
		default:
			k.command_play_simple(e.Message)
		}
	}

	re_ivona := regexp.MustCompile("\\$\\$\\$(.*)")
	result_ivona := re_ivona.FindStringSubmatch(e.Message)
	if len(result_ivona) == 2 {
		k.command_play_ivona(result_ivona[1])
	}

}

func (k *Kotik) OnUserChange(e *gumble.UserChangeEvent) {
	fmt.Println(k.Timestamp() + "OnUserChange()-> User:" + e.User.Name + " | Type:" + fmt.Sprint(e.Type))
}

func (k *Kotik) OnChannelChange(e *gumble.ChannelChangeEvent) {
	fmt.Println(k.Timestamp() + "OnChannelChange()-> Channel:" + e.Channel.Name + " | Type: " + fmt.Sprint(e.Type))
}

func (k *Kotik) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	fmt.Println(k.Timestamp() + "OnPermissionDenied()-> User:" + e.User.Name + " | Type: " + fmt.Sprint(e.Type))
}

func (k *Kotik) OnUserList(e *gumble.UserListEvent) {
	fmt.Println(k.Timestamp() + "OnUserList()")
}

func (k *Kotik) OnACL(e *gumble.ACLEvent) {
	fmt.Println(k.Timestamp() + "OnACL()")
}

func (k *Kotik) OnBanList(e *gumble.BanListEvent) {
	fmt.Println(k.Timestamp() + "OnBanList()")
}

func (k *Kotik) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
	fmt.Println(k.Timestamp() + "OnContextActionChange()")
}

func (k *Kotik) OnServerConfig(e *gumble.ServerConfigEvent) {
	fmt.Println(k.Timestamp() + "OnServerConfig()")
}

//Commands
func (k *Kotik) command_disconnect(e *gumble.TextMessageEvent) {
	//k.Client.Disconnect()
}

func (k *Kotik) command_help(e *gumble.User) {
	str := "<br/>" +
		"#[soundfile]      : проиграть звук<br/>" +
		"$$$[text]         : произнести текст<br/>" +
		"!disconnect       : отключить бота<br/>" +
		"!list_channels    : список каналов<br/>" +
		"!list_sounds      : список звуков<br/>" +
		"!moveto [id/name] : перенести бота на другой канал<br/>" +
		"!help             : эта команда<br/>" +
		"!stop             : остановить воспроизведение звука<br/>" +
		"!pause            : приостановить воспроизведение<br/>" +
		"!resume           : восстановить воспроизведение<br/>" +
		"!volume [float]   : установить громкость. Максимальная 100, cтандартная 50, минимальная 0, шаг 1<br/>" +
		"!update           : делает апдейт<br/>" +
		"!status           : информация про бота<br/>"
	e.Send(str)
}

func (k *Kotik) command_list_channels(e *gumble.User) {
	root := k.Client.Channels[0]
	var channels string
	if root != nil {
		channels += "<br/>" + fmt.Sprint(root.ID) + ": " + root.Name + "<br/>"
		if root.Children != nil {
			channels += k.command_list_channels_printchild(root.Children, 1)
		}
	}
	e.Send(channels)
}

func (k *Kotik) command_list_channels_printchild(children gumble.Channels, level int) string {
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
			out += k.command_list_channels_printchild(children[uint32(j)].Children, level+1)
		}
	}
	return out
}

func (k *Kotik) command_list_sounds(e *gumble.User) {
	files, _ := ioutil.ReadDir("./sounds/")
	var sounds string = ""
	for _, f := range files {
		sounds += "<br/>" + f.Name()
	}
	e.Send(sounds)
}

func (k *Kotik) command_moveto(text string) {
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

func (k *Kotik) command_pause() {
	k.Audio.Pause()
}

func (k *Kotik) command_play_ivona(text string) {
	if k.Audio.IsPlaying() {
		return
	}

	opts_input := ivona.Input{Data: text, Type: "text/plain"}
	opts_outputformat := ivona.OutputFormat{Codec: "MP3", SampleRate: 22050}
	opts_parameters := ivona.Parameters{Rate: "medium", Volume: "medium", SentenceBreak: 100, ParagraphBreak: 500}
	opts_voice := ivona.Voice{Name: "Maxim", Language: "ru-RU", Gender: "Male"}
	speechopts := ivona.SpeechOptions{Input: &opts_input, OutputFormat: &opts_outputformat, Parameters: &opts_parameters, Voice: &opts_voice}

	respond, err := k.Ivona.CreateSpeech(speechopts)
	if err == nil {
		cb := &ClosingBuffer{bytes.NewBuffer(respond.Audio)}
		var rc io.ReadCloser
		rc = cb
		k.Audio.Source = gumble_ffmpeg.SourceReader(rc)
		k.Audio.Play()
	}
}

func (k *Kotik) command_play_simple(text string) {
	if k.Audio.IsPlaying() {
		return
	}

	filename := strings.Split(text, "#")[1]
	if _, err := os.Stat("./sounds/" + filename); err == nil {
		k.Audio.Source = gumble_ffmpeg.SourceFile("./sounds/" + filename)
		k.Audio.Play()
	}
}

func (k *Kotik) command_play_ymusic(text string, e *gumble.User) {
	if k.Audio.IsPlaying() {
		return
	}

	ym := YMusic{}
	trackname := strings.Split(text, "#ymusic ")[1]
	file, title := ym.getTrack(trackname)
	if title != "" {
		e.Send("Найдена композиция: " + title)
	} else {
		e.Send("Композиция не найдена")
	}
	if file != nil {
		k.Audio.Source = gumble_ffmpeg.SourceReader(file)
		k.Audio.Play()
	}
}

func (k *Kotik) command_resume() {
	k.Audio.Play()
}

func (k *Kotik) command_status(e *gumble.User) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	var str string = ""
	str = "<br/>" +
		"Uptime       : " + strconv.FormatFloat(time.Since(k.connectTime).Hours(), 'f', 2, 64) + "hours <br/>" +
		"Memory alloc : " + strconv.FormatFloat(float64(mem.Alloc)/1024.0/1024.0, 'f', 2, 64) + "MB <br/>"+
		"Volume       : " + strconv.FormatInt(int64(k.Audio.Volume*50.00),10) + "% <br/>"
		
	e.Send(str)
}

func (k *Kotik) command_stop() {
	k.Audio.Stop()
}

func (k *Kotik) command_update(e *gumble.User) {
	if runtime.GOOS == "linux" {
		args := []string{"arg1"}
		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		os.StartProcess("./update_linux.sh", args, procAttr)
		k.Client.Disconnect()
	} else {
		fmt.Println(k.Timestamp() + "!update works only in production")
	}

}

func (k *Kotik) command_volume(text string) {
	re_sound := regexp.MustCompile("^!volume[ ]?(\\d+)")
	result_sound := re_sound.FindStringSubmatch(text)
	if len(result_sound) == 2 {
		i, err := strconv.Atoi(result_sound[1])
		if err == nil {
			if i > 100 {
				i = 100
			}
			k.Audio.Volume = float32(i)/50.00
		}
	}
}

func main() {
	//Flags
	flag_username := flag.String("username", "Kotik-dev", "username of the bot")
	flag_password := flag.String("password", "", "user password")
	flag_dev := flag.Bool("dev",true,"development mode")
	flag_server := flag.String("server", "direct.galyonkin.com:64738", "address of the server")
	flag_certificateFile := flag.String("certificate", "", "user certificate file (PEM)")
	flag_keyFile := flag.String("key", "", "user certificate key file (PEM)")
	flag_lock := flag.String("lock", "", "server certificate lock file")
	flag.Parse()

	//Config
	k := Kotik{}
	k.Config = gumble.NewConfig()
	k.Config.Username = *flag_username
	k.Config.Password = *flag_password
	k.Config.Address = *flag_server
	k.Config.TLSConfig.InsecureSkipVerify = true

	//Client creation
	k.Client = gumble.NewClient(k.Config)
	k.Audio = gumble_ffmpeg.New(k.Client)
	k.Audio.Command = "ffmpeg"

	//Ivona creation
	k.Ivona = ivona.New("GDNAIKHZ6EJPBXXTKZFA", "akU4WnCw2XozeJeMsnS7pVqlBsLgqF4FQbVRnu05")

	//Attach listeners
	k.Client.Attach(gumbleutil.AutoBitrate)
	k.Client.Attach(&k)

	//TLS
	if *flag_lock != "" {
		gumbleutil.CertificateLockFile(k.Client, *flag_lock)
	}
	
	if _, err := os.Stat("./Kotik.pem"); err == nil && *flag_dev==false{ 
		*flag_certificateFile = "./Kotik.pem"
	}else if _, err := os.Stat("./Kotik-dev.pem"); err == nil && *flag_dev==true{
		*flag_certificateFile = "./Kotik-dev.pem"
	}
	
	if *flag_certificateFile != "" {
		if *flag_keyFile == "" {
			flag_keyFile = flag_certificateFile
		}
		if certificate, err := tls.LoadX509KeyPair(*flag_certificateFile, *flag_keyFile); err != nil {
			panic(err)
		} else {
			k.Config.TLSConfig.Certificates = append(k.Config.TLSConfig.Certificates, certificate)
		}
	}

	//Connect
	if err := k.Client.Connect(); err != nil {
		panic(err)
	}

	<-k.keepAlive
}
