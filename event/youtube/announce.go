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

package youtube

import (
	"cake4everybot/database"
	"cake4everybot/util"
	webYT "cake4everybot/webserver/youtube"
	"database/sql"

	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	videoBaseURL   string = "https://youtu.be/%s"
	channelBaseURL string = "https://youtube.com/channel/%s"
)

// Announce takes a youtube video and announces it in discord channels
func Announce(s *discordgo.Session, event *webYT.Video) {
	announcements, err := database.GetAnnouncement(database.AnnouncementPlatformYoutube, event.ChannelID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	var (
		videoURL   = fmt.Sprintf(videoBaseURL, event.ID)
		channelURL = fmt.Sprintf(channelBaseURL, event.ChannelID)
		thumb      = event.Thumbnails["high"]
	)

	embed := &discordgo.MessageEmbed{
		Title:       event.Title,
		Description: saveTrimText(event.Description, 100),
		URL:         videoURL,
		Color:       0xFF0000,
		Author:      &discordgo.MessageEmbedAuthor{URL: channelURL, Name: event.Channel},
		Image:       &discordgo.MessageEmbedImage{URL: thumb.URL, Width: thumb.Width, Height: thumb.Height},
	}
	util.SetEmbedFooter(s, "youtube.embed_footer", embed)
	embeds := []*discordgo.MessageEmbed{embed}

	// send the embed to the channels
	for _, announcement := range announcements {
		if announcement.Notification != "" && strings.Contains(announcement.Notification, "%s") {
			embed.Author.Name = fmt.Sprintf(announcement.Notification, event.Channel)
		} else if announcement.Notification != "" {
			embed.Author.Name = announcement.Notification
		}

		var err error
		if announcement.RoleID == "" {
			// send without a ping
			_, err = s.ChannelMessageSendEmbeds(announcement.ChannelID, embeds)
		} else {
			// send with a ping
			data := &discordgo.MessageSend{
				Content: fmt.Sprintf("<@&%s>", announcement.RoleID),
				Embeds:  embeds,
			}
			var msg *discordgo.Message
			msg, err = s.ChannelMessageSendComplex(announcement.ChannelID, data)
			if err == nil {
				_, err = s.ChannelMessageEditEmbeds(announcement.ChannelID, msg.ID, embeds)
			}
		}

		if err != nil {
			log.Printf("Error on sending video announcement to channel %s/%s: %v", announcement.GuildID, announcement.ChannelID, err)
		}
	}
}

// saveTrimText returns a trimmed version of the given string. It
// will be trimmed to n characters but then continues to the next
// space character. If s is shorter or equal to n, then s is
// returned. When words get cut of a "..." gets appended at the end.
func saveTrimText(s string, n int) string {
	s = strings.ReplaceAll(s, "\n\t", " ")
	if n <= 0 || s == " " {
		return ""
	}
	if len(s) <= n {
		return s
	}

	// offset
	o := strings.Index(s[n-3:], " ")
	if o == -1 || len(s) <= n+o+1 {
		return s
	}

	return s[:n+o-2] + "..."
}
