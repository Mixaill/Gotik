package main

/*
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
*/
