package backends

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"

	"../common"
	"../services"
)

type Discord struct {
	discord  *discordgo.Session
	self     *discordgo.User
	guildID  string
	services *services.Services
}

func NewDiscord() *Discord {
	p := new(Discord)
	return p
}

func (k *Discord) Start(fl map[string]string, s *services.Services) {

	k.services = s

	k.discord = &discordgo.Session{
		OnMessageCreate: k.onMessageCreate,
	}

	// Login to the Discord server and store the authentication token
	err := k.discord.Login(fl["Flag_discord_email"], fl["Flag_discord_password"])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open websocket connection
	err = k.discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	k.internal_findSelf()
	k.internal_findGuild()

}

func (k *Discord) internal_findSelf() {
	// Get Authenticated User's information
	var err error
	k.self, err = k.discord.User("@me")
	if err != nil {
		fmt.Println("backends/discord/start_findSelf(): error fetching self, ", err)
		return
	}
}

func (k *Discord) internal_findGuild() {
	ch, err := k.discord.UserGuilds()

	if err != nil {
		fmt.Println("backends/discord/start_findGuild(): error fetching guild, ", err)
		return
	}
	k.guildID = ch[0].ID
}

func (k *Discord) internal_getChannels_text() []*discordgo.Channel {
	var text []*discordgo.Channel

	channels, err := k.discord.GuildChannels(k.guildID)
	if err != nil {
		fmt.Println("backends/discord/Command_Channels_List(): error fetching guild channels ", err)
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
		fmt.Println("backends/discord/Command_Channels_List(): error fetching guild channels ", err)
		return voice
	}

	for _, val := range channels {
		if val.Type == "voice" {
			voice = append(voice, val)
		}
	}

	return voice
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
}

func (k *Discord) Command_Audio_Pause() {
}

func (k *Discord) Command_Audio_Stop() {
}

func (k *Discord) Command_Audio_Play_File(text string) {
}

func (k *Discord) Command_Audio_Play_Ivona(text string, language string) {
}

func (k *Discord) Command_Audio_Volume(text string) {
}

/////Commands/Channels

func (k *Discord) Command_Channels_List(user string) {

	text := k.internal_getChannels_text()
	voice := k.internal_getChannels_voice()

	var str string
	str = ""

	str = "Текстовые каналы:\n"
	for id, val := range text {
		str = str + "t" + strconv.Itoa(id+1) + ": " + val.Name + "\n"
	}

	str = str + "\nГолосовые каналы:\n"
	for id, val := range voice {
		str = str + "v" + strconv.Itoa(id+1) + ": " + val.Name + "\n"
	}

	k.internal_sendMessage_private(str, user)
}

func (k *Discord) Command_Channels_Moveto(text string) {
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

/////
///// Info
/////

func (k *Discord) Info_Name() string {
	return "discord"
}
