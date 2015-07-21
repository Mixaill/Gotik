package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type Twitter struct {
	api        *anaconda.TwitterApi
	users      []anaconda.User
	twitTurn   []anaconda.Tweet
	lastUpdate time.Time
}

func (tw *Twitter) Init() {
	anaconda.SetConsumerKey("ememheL6J1S5r1RwBzoaKc4KR")
	anaconda.SetConsumerSecret("HKyDHzrirIqVSNJu414ajTLcGXoWIpnkrpNERAY1AEgKP3an6C")
	tw.api = anaconda.NewTwitterApi("50292089-521odSFl4uKfV0oLcIryalgetYuP0nxo7i6bQCjTE", "djU5cDWW9rsR2kWmYExR9CxLpJdZdaNTWe7YOubV9VB5b")
	//tw.lastUpdate = time.Now()
}

func (tw *Twitter) UsersAdd(name string) {
	//check if user exists in twitter
	u, err := tw.api.GetUsersLookup(name, nil)
	if err == nil && len(u) > 0 {
		//check if user already exists in tw.Users
		for _, user := range tw.users {
			if user.ScreenName == name {
				fmt.Println("user @" + name + " already existed")
				return
			}
		}
		//add user
		tw.users = append(tw.users, u[0])
	} else {
		fmt.Println(err)
	}
}

func (tw *Twitter) UsersDel(name string) {
	//delete user from tw.users
	var id int = -1
	for i, user := range tw.users {
		if user.ScreenName == name {
			id = i
		}
	}
	if id > -1 {
		tw.users = tw.users[:id+copy(tw.users[id:], tw.users[id+1:])]

		//delete user tweets from tw.twitTurn
		twids := []int{}
		for i, twit := range tw.twitTurn {
			if twit.User.ScreenName == name {
				twids = append([]int{i}, twids...)
			}
		}
		for _, j := range twids {
			tw.twitTurn = tw.twitTurn[:j+copy(tw.twitTurn[j:], tw.twitTurn[j+1:])]
		}
	}
}

func (tw *Twitter) UsersGet() string {
	var users string = ""
	for _, user := range tw.users {
		users += user.ScreenName + ", "
	}
	return users
}

func (tw *Twitter) TurnFill() {
	for _, user := range tw.users {
		v := url.Values{}
		v.Set("screen_name", user.ScreenName)
		v.Set("count", "30")
		v.Set("exclude_replies", "true")
		timeline, err := tw.api.GetUserTimeline(v)
		if err != nil {
			fmt.Println(err)
		}
		for _, twit := range timeline {
			if time, _ := twit.CreatedAtTime(); time.After(tw.lastUpdate) {
				tw.twitTurn = append(tw.twitTurn, twit)
			}
		}
	}
	tw.lastUpdate = time.Now()
}

func (tw *Twitter) TurnRelease() []string {
	var twits []string
	for _, twit := range tw.twitTurn {
		t, _ := twit.CreatedAtTime()
		str := "@ " + twit.User.ScreenName + " " + strconv.FormatInt(int64(time.Since(t).Minutes()), 10) + " минуты назад. " + strings.Replace(twit.Text, "\n", "\\n", -1)
		twits = append(twits, str)
	}
	tw.twitTurn = tw.twitTurn[:0]
	return twits
}
