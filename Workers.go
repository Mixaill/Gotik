package main

import (
	"time"

	"./backends"
	"./services"
)

func TwitterWorker(backnds backends.Backends, srvcs services.Services) {
	ticker := time.NewTicker(srvcs.Twitter.UpdateRateGet())
	go func() {
		for {
			select {
			case <-ticker.C:
				go TwitterWorkerProcess(backnds, srvcs)
				ticker.Stop()
				go TwitterWorker(backnds, srvcs)
				return
			}
		}
	}()
}

func TwitterWorkerProcess(backnds backends.Backends, srvcs services.Services) {
	srvcs.Twitter.TurnFill()
	twits := srvcs.Twitter.TurnRelease()

	backnds.Discord.Command_Twitter_ReadTwits(twits)
	backnds.Mumble.Command_Twitter_ReadTwits(twits)
}
