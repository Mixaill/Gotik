package backends

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"

	"../common"
	"../services"
	"../utils"
)

/////
/////Ctor
/////

type Discord struct {
	discord  *discordgo.Session
	self     *discordgo.User
	guildID  string
	services *services.Services

	conf_connectTime    time.Time
	conf_twitterChannel int
	conf_twitterEnable bool
}

func NewDiscord() *Discord {
	p := new(Discord)
	return p
}

/////
/////Helpers
/////

func (k *Discord) Start(fl map[string]string, s *services.Services) {

	k.services = s

	k.discord = &discordgo.Session{
		OnMessageCreate: k.onMessageCreate,
	}

	// Login to the Discord server and store the authentication token
	err := k.discord.Login(fl["Flag_discord_email"], fl["Flag_discord_password"])
	if err != nil {
		fmt.Println("backends/discord/Start()-> error: ", err)
		return
	}

	// Open websocket connection
	err = k.discord.Open()
	if err != nil {
		fmt.Println("backends/discord/Start()->", err)
		return
	}

	//find self object and server(guild)ID
	k.internal_findSelf()
	k.internal_findGuild()

	//connect to voice channel #1
	k.Command_Channels_Moveto("!channels_moveto 1")

	//connect twitter broadcasting to text channel #1
	k.conf_twitterChannel = 1
	k.conf_twitterEnable = true
}

func (k *Discord) internal_findSelf() {
	// Get Authenticated User's information
	var err error
	k.self, err = k.discord.User("@me")
	if err != nil {
		fmt.Println("backends/discord/internal_findSelf(): error fetching self, ", err)
		return
	}
}

func (k *Discord) internal_findGuild() {
	ch, err := k.discord.UserGuilds()

	if err != nil {
		fmt.Println("backends/discord/internal_findGuild(): error fetching guild, ", err)
		return
	}
	k.guildID = ch[0].ID
}

func (k *Discord) internal_getChannels_text() []*discordgo.Channel {
	var text []*discordgo.Channel

	channels, err := k.discord.GuildChannels(k.guildID)
	if err != nil {
		fmt.Println("backends/discord/internal_getChannels_text(): error fetching guild channels ", err)
		return text
	}

	for _, val := range channels {
		if val.Type == "text" {
			text = append(text, val)
		}
	}

	return text
}

func (k *Discord) internal_getChannels_voice() []*discordgo.Channel {
	var voice []*discordgo.Channel

	channels, err := k.discord.GuildChannels(k.guildID)
	if err != nil {
		fmt.Println("backends/discord/internal_getChannels_voice(): error fetching guild channels ", err)
		return voice
	}

	for _, val := range channels {
		if val.Type == "voice" {
			voice = append(voice, val)
		}
	}

	return voice
}

func (k *Discord) internal_sendMessage_group(message string, channelID string) {

	channelExists := false

	//fetching channel
	channels, err := k.discord.GuildChannels(k.guildID)
	if err != nil {
		fmt.Println("backends/discord/internal_sendMessage_group(): error fetching group channels, ", err)
		return
	}

	//checking if channel exists
	for _, element := range channels {
		if element.ID == channelID {
			channelExists = true
		}
	}
	if channelExists == false {
		fmt.Println("backends/discord/internal_sendMessage_group(): error finding channel, ")
		return
	}

	//sending message
	_, err = k.discord.ChannelMessageSend(channelID, message)
	if err != nil {
		fmt.Println("backends/discord/internal_sendMessage_group(): error sending group message, ", err)
		return
	}
}

func (k *Discord) internal_sendMessage_private(message string, user string) {
	var chanid string
	chanid = ""

	//find in existing private channels
	if chanid == "" {
		ch, err := k.discord.UserChannels()

		if err != nil {
			fmt.Println("backends/discord/internal_sendMessage_private(): error fetching user channels, ", err)
			return
		}

		for _, element := range ch {
			if element.Recipient.Username == user {
				chanid = element.ID
			}
		}
	}

	//create channel if not found
	if chanid == "" {
		guild, err := k.discord.Guild(k.guildID)
		if err != nil {
			fmt.Println("backends/discord/internal_sendMessage_private(): error fetching guild, ", err)
			return
		}

		members := guild.Members
		for _, element := range members {
			if element.User.Username == user {
				chn, err := k.discord.UserChannelCreate(element.User.ID)

				if err != nil {
					fmt.Println("backends/discord/internal_sendMessage_private(): error creating private channel, ", err)
					return
				}
				chanid = chn.ID
			}
		}
	}

	if chanid == "" {
		fmt.Println("backends/discord/internal_sendMessage_private(): error finding user ", user)
		return
	}

	_, err := k.discord.ChannelMessageSend(chanid, message)
	if err != nil {
		fmt.Println("backends/discord/internal_sendMessage_private(): error sending private message, ", err)
		return
	}
}

/////
/////Events
/////

func (k *Discord) onMessageCreate(s *discordgo.Session, m *discordgo.Message) {
	var err error

	// if msg is from self, ignore it entirely.
	if k.self == nil || k.self.ID == "" {
		k.self, err = k.discord.User("@me")
		if err != nil {
			fmt.Println("backends/discord/onMessageCreate: Error fetching self, ", err)
		}
	}
	if m.Author.ID == k.self.ID {
		return
	}

	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	common.Command_Process(m.Content, m.Author.Username, k)

}

/////
/////Commands
/////

/////Commands/Audio
func (k *Discord) Command_Audio_List(user string) {
	k.internal_sendMessage_private(common.Command_Audio_List(k), user)
}

func (k *Discord) Command_Audio_Resume() {
	//unimplemented
}

func (k *Discord) Command_Audio_Pause() {
	//unimplemented
}

func (k *Discord) Command_Audio_Stop() {
	dgvoice.KillPlayer()
}

func (k *Discord) Command_Audio_Play_File(text string) {
	if !k.discord.Voice.Ready {
		return
	}

	filepath := utils.GetAudioFilePath(text)
	if filepath != "" {
		dgvoice.PlayAudioFile(k.discord, filepath)
	}

}

func (k *Discord) Command_Audio_Play_Ivona(text string, language string) {
	if !k.discord.Voice.Ready {
		return
	}

	var tt string
	if strings.HasPrefix(text, "$$$") {
		tt = strings.Split(text, "$$$")[1]
	} else {
		tt = strings.SplitN(text, " ", 2)[1]
	}

	if language == "" {
		language = k.services.YTranslate.Detect(tt)
	}

	filename := k.services.Ivona.GetAudio_File(tt, language)

	if filename != "" {
		dgvoice.PlayAudioFile(k.discord, filename)
	}
}

func (k *Discord) Command_Audio_Volume(text string) {
	//unimplemented
}

/////Commands/Channels

func (k *Discord) Command_Channels_List(user string) {

	voice := k.internal_getChannels_voice()

	var str string
	str = ""

	str = str + "\nКаналы:\n"
	for id, val := range voice {
		str = str + strconv.Itoa(id+1) + ": " + val.Name + "\n"
	}

	k.internal_sendMessage_private(str, user)
}

func (k *Discord) Command_Channels_Moveto(text string) {
	re_id := regexp.MustCompile("^!channels_moveto ([0-9]+)")
	var result_id []string = re_id.FindStringSubmatch(text)
	if len(result_id) == 2 {
		i, err := strconv.Atoi(result_id[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		channels := k.internal_getChannels_voice()
		if i <= len(channels) {
			err := k.discord.ChannelVoiceJoin(channels[i-1].GuildID, channels[i-1].ID, false, true)
			if err != nil {
				fmt.Println("backends/discord/Command_Channels_Moveto(): error connecting voice channel", err)
				return
			}
		}
	}
}

/////Commands/other
func (k *Discord) Command_Help(user string) {
	k.internal_sendMessage_private(common.Command_Help(k), user)
}

func (k *Discord) Command_Update() {
	common.Command_Update()
}

func (k *Discord) Command_Disconnect() {
}

func (k *Discord) Command_Status(user string) {
}

/////Commands/twitter
func (k *Discord) Command_Twitter_ReadTwits(twits []anaconda.Tweet) {
	if k.conf_twitterEnable == true {
		for _, val := range twits {
			k.internal_sendMessage_group(utils.TwitterFormatForText(val), k.internal_getChannels_text()[k.conf_twitterChannel-1].ID)
			k.Command_Audio_Play_Ivona(utils.TwitterFormatForAudio(val), val.Lang)
		}
	}
}

func (k *Discord) Command_Twitter_Status(user string) {
	var str string
	str = "\n" +
		"Twitter subscriptions : " + k.services.Twitter.UsersGet() + "\n" +
		"Twitter update rate   : " + strconv.FormatFloat(k.services.Twitter.UpdateRateGet().Minutes(), 'f', 2, 64) + " minutes \n" +
		"Twitter channel       : " + strconv.Itoa(k.conf_twitterChannel) + "." + k.internal_getChannels_text()[k.conf_twitterChannel-1].Name + "\n"
	k.internal_sendMessage_private(str, user)
}

/////
///// Getters
/////
func (k *Discord) Get_ConnectTime() time.Time {
	return k.conf_connectTime
}

func (k *Discord) Get_Services() *services.Services {
	return k.services
}

func (k *Discord) Get_Volume() float32 {
	//unimplemented
	return 0.0
}

/////
///// Info
/////

func (k *Discord) Info_Name() string {
	return "discord"
}
