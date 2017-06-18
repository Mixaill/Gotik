package main

import (
	"flag"

	"./backends"
	"./config"
	"./services"
)

func main() {

	//Flags
	flag_dev := flag.String("dev", "false", "development mode")

	flag_discord_email := flag.String("discord_email", config.Discord_email, "email of the Discord bot")
	flag_discord_password := flag.String("discord_password", config.Discord_password, "Discord use password")
	flag_discord_username := flag.String("discord_username", config.Discord_username, "name of the Discord bot")

	flag_mumble_username := flag.String("mumble_username", config.Mumble_username, "username of the Mumble bot")
	flag_mumble_password := flag.String("mumble_password", config.Mumble_password, "user Mumble password")
	flag_mumble_server := flag.String("mumble_server", config.Mumble_server, "address of the Mumble server")
	flag_mumble_cert := flag.String("mumble_cert", config.Mumble_cert, "user certificate file (PEM)")
	flag_mumble_cert_key := flag.String("mumble_cert_key", config.Mumble_cert_key, "user certificate key file (PEM)")
	flag_mumble_cert_lock := flag.String("mumble_cert_lock", config.Mumble_cert_lock, "server certificate lock file")

	flag.Parse()

	flags := make(map[string]string)
	flags["Flag_dev"] = *flag_dev
	flags["Flag_discord_email"] = *flag_discord_email
	flags["Flag_discord_password"] = *flag_discord_password
	flags["Flag_discord_username"] = *flag_discord_username
	flags["Flag_mumble_username"] = *flag_mumble_username
	flags["Flag_mumble_password"] = *flag_mumble_password
	flags["Flag_mumble_server"] = *flag_mumble_server
	flags["Flag_mumble_cert"] = *flag_mumble_cert
	flags["Flag_mumble_cert_key"] = *flag_mumble_cert_key
	flags["Flag_mumble_cert_lock"] = *flag_mumble_cert_lock

	//Kotik init
	s := services.Services{}
	b := backends.Backends{}

	//Channels
	ch_keepAlive := make(chan int)

	//Services initialization
	s.Ivona = services.NewIvona()
	//s.Twitter = services.NewTwitter()
	//s.Twitter.Bootstrap()
	s.YTranslate = services.NewYTranslate()

	//Backends initialization
	//b.Discord = backends.NewDiscord()
	b.Mumble = backends.NewMumble()

	//go b.Discord.Start(flags, &s)
	go b.Mumble.Start(flags, &s)

	//Workers

	//TwitterWorker(b, s)

	<-ch_keepAlive
}
