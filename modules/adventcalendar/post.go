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

package adventcalendar

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// Post is a scheduled function to run everyday at 8:00
func Post(s *discordgo.Session) {
	t := time.Now()
	if t.Month() != 12 || t.Day() > 24 {
		return
	}
	log.Printf("New Post for %s", t.Format("_2. Jan"))

	channels, err := util.GetChannelsFromDatabase(s, "adventcalendar_channel")
	if err != nil {
		log.Printf("ERROR: Could not get advent calendar channel: %+v", err)
		return
	}

	for _, channelID := range channels {
		data := postData(t)
		_, err = s.ChannelMessageSendComplex(channelID, data)
		if err != nil {
			log.Printf("Failed to send new post for advent calendar in channel '%s': %+v", channelID, err)
			return
		}
	}
}

func postData(t time.Time) *discordgo.MessageSend {
	var line1 string
	if t.Day() == 23 && t.Month() == 12 {
		line1 = lang.GetDefault("module.adventcalendar.post.message.day_23")
	} else if t.Day() == 24 && t.Month() == 12 {
		line1 = lang.GetDefault("module.adventcalendar.post.message.day_24")
	} else {
		format := lang.GetDefault("module.adventcalendar.post.message")
		line1 = fmt.Sprintf(format, 24-t.Day(), t.Day())
	}
	line2 := lang.GetDefault("module.adventcalendar.post.message2")
	message := fmt.Sprintf("%s\n%s", line1, line2)

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				fmt.Sprintf("%s.post.%s", Component.ID(Component{}), t.Format("2006.01.02")),
				lang.GetDefault("module.adventcalendar.post.button"),
				discordgo.PrimaryButton,
				util.GetConfigComponentEmoji("adventcalendar"),
			),
		}},
	}

	filepath := fmt.Sprintf("%s/%d.png", viper.GetString("event.adventcalendar.images"), t.Day())
	log.Printf("image path: %s", filepath)
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Printf("Failed to open advent calendar image: %+v", err)
		return nil
	}

	return &discordgo.MessageSend{
		Content:    message,
		Components: components,
		Files: []*discordgo.File{
			{
				Name:        fmt.Sprintf("c4e_advent_calendar_%d.png", t.Day()),
				ContentType: "image/png",
				Reader:      file,
			},
		},
	}
}
