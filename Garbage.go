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

/*
func command_twitter_fetch() {
	go k.command_twitter_bootstrap()

	//twitter fetching
	ticker := time.NewTicker(k.Twitter_i.UpdateRateGet())
	go func() {
		for {
			select {
			case <-ticker.C:
				k.command_twitter_process(nil)
			case <-k.ch_TwitterFetching:
				ticker.Stop()
				return
			}
		}
	}()
}*/

/*
func (k *Mumble) command_twitter_process(sender *gumble.User) {
	k.Twitter_i.TurnFill()
	if sender != nil {
		sender.Send("Сейчас будет зачитано " + strconv.FormatInt(int64(k.Twitter_i.TurnSize()), 10) + " твитов")
	}

	twits := k.Twitter_i.TurnRelease()
	ch_CallBack := make(chan int)
	for _, twit := range twits {
		go k.command_play_ivona(twit, ch_CallBack)
		<-ch_CallBack
	}
}
*/

/*
func (k *Kotik) command_status(e *gumble.User) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	var str string = ""
	str = "<br/>" +
		"Uptime                : " + strconv.FormatFloat(time.Since(k.connectTime).Hours(), 'f', 2, 64) + " hours <br/>" +
		"Memory alloc          : " + strconv.FormatFloat(float64(mem.Alloc)/1024.0/1024.0, 'f', 2, 64) + " MB <br/>" +
		"Volume                : " + strconv.FormatInt(int64(k.conf_volume*50.00), 10) + "% <br/>" +
		"Twitter subscriptions : " + k.Twitter_i.UsersGet() + "<br/>" +
		"Twitter update rate   : " + strconv.FormatFloat(k.Twitter_i.UpdateRateGet().Minutes(), 'f', 2, 64) + " minutes <br/>"

	e.Send(str)
}
*/
