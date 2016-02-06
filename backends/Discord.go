package backends

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"../services"
)

type Discord struct {
	discord  *discordgo.Session
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
}

func (k *Discord) onMessageCreate(s *discordgo.Session, m *discordgo.Message) {
	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
