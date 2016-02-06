package backends

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"../common"
	"../services"
)

type Discord struct {
	discord  *discordgo.Session
	self     *discordgo.User
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

	// Get Authenticated User's information
	k.self, err = k.discord.User("@me")
	if err != nil {
		fmt.Println("error fetching self, ", err)
		return
	}

}

func (k *Discord) onMessageCreate(s *discordgo.Session, m *discordgo.Message) {
	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

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

	common.Command_Process(m.Content, m.Author.Username, k)

}

/////
/////Commands
/////

/////Commands/Audio
func (k *Discord) Command_Audio_List(user string) {
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
}

func (k *Discord) Command_Channels_Moveto(text string) {
}

/////Commands/other
func (k *Discord) Command_Help(user string) {
}

func (k *Discord) Command_Update() {
}

func (k *Discord) Command_Disconnect() {
}

func (k *Discord) Command_Status(user string) {
}
