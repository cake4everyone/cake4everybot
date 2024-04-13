// Copyright 2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event

import (
	"cake4everybot/modules/adventcalendar"
	"cake4everybot/modules/birthday"

	webYT "cake4everybot/webserver/youtube"

	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

func addScheduledTriggers(dc *discordgo.Session, t *twitchgo.Session, webChan chan struct{}) {
	go scheduleFunction(dc, t, 0, 0,
		adventcalendar.Midnight,
	)

	go scheduleFunction(dc, t, viper.GetInt("event.morning_hour"), viper.GetInt("event.morning_minute"),
		birthday.Check,
		adventcalendar.Post,
	)

	go refreshYoutube(webChan)
}

func scheduleFunction(dc *discordgo.Session, t *twitchgo.Session, hour, min int, callbacks ...interface{}) {
	if len(callbacks) == 0 {
		return
	}
	log.Printf("scheduled %d function(s) for %2d:%02d!", len(callbacks), hour, min)
	time.Sleep(time.Second * 5)
	for {
		now := time.Now()

		nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
		if nextRun.Before(now) {
			nextRun = nextRun.Add(time.Hour * 24)
		}
		time.Sleep(nextRun.Sub(now))

		for _, c := range callbacks {
			switch f := c.(type) {
			case func(*discordgo.Session):
				f(dc)
			case func(*twitchgo.Session):
				f(t)
			}
		}
	}
}

func refreshYoutube(webChan chan struct{}) {
	<-webChan
	for {
		webYT.RefreshSubscriptions()

		// loop every 4 days
		time.Sleep(4 * 24 * time.Hour)
	}
}
