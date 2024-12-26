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
	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/util"
	webYT "cake4everybot/webserver/youtube"

	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	videoBaseURL   string = "https://youtu.be/%s"
	channelBaseURL string = "https://youtube.com/channel/%s"
)

type guild struct {
	guild   *discordgo.Guild
	channel *discordgo.Channel
	ping    string
}

// Announce takes a youtube video and announces it in discord channels
func Announce(s *discordgo.Session, event *webYT.Video) {
	guilds, err := getGuilds(s)
	if err != nil {
		log.Printf("Error on getting channels: %v\n", err)
		return
	}
	if len(guilds) == 0 {
		log.Printf("No channels to announce video. Dropping announcement for 'www.youtu.be/%s'", event.ID)
		return
	}

	var (
		videoURL   = fmt.Sprintf(videoBaseURL, event.ID)
		channelURL = fmt.Sprintf(channelBaseURL, event.ChannelID)
		title      = fmt.Sprintf(lang.GetDefault("youtube.msg.new_vid"), event.Channel)
		thumb      = event.Thumbnails["high"]
	)

	embed := &discordgo.MessageEmbed{
		Title:       event.Title,
		Description: saveTrimText(event.Description, 100),
		URL:         videoURL,
		Color:       0xFF0000,
		Author:      &discordgo.MessageEmbedAuthor{URL: channelURL, Name: title},
		Image:       &discordgo.MessageEmbedImage{URL: thumb.URL, Width: thumb.Width, Height: thumb.Height},
	}
	util.SetEmbedFooter(s, "youtube.embed_footer", embed)

	// send the embed to the channels
	for _, g := range guilds {
		var err error
		if g.ping == "<@&0>" {
			// send without a ping
			_, err = s.ChannelMessageSendEmbed(g.channel.ID, embed)
		} else {
			// send with a ping
			data := &discordgo.MessageSend{
				Content: g.ping,
				Embed:   embed,
			}
			_, err = s.ChannelMessageSendComplex(g.channel.ID, data)
		}

		if err != nil {
			log.Printf("Error on sending video announcement to channel %s (#%s) in guild %s (%s): %v", g.channel.ID, g.channel.Name, g.guild.ID, g.guild.Name, err)
		}
	}
}

// getGuilds returns a list of guild object containing all guilds
// (that specified an youtube announcement channel) as well as the
// announcement channel an the role as pingable string.
func getGuilds(s *discordgo.Session) (guilds []guild, err error) {
	rows, err := database.Query("SELECT id,youtube_channel,youtube_role FROM guilds")
	if err != nil {
		return guilds, err
	}
	defer rows.Close()

	var guildID, channelID, roleID uint64
	for rows.Next() {
		err = rows.Scan(&guildID, &channelID, &roleID)
		if err != nil {
			log.Printf("Error on scanning row (channel/%d/%d) from database: %v\n", guildID, channelID, err)
			continue
		}

		if guildID == 0 || channelID == 0 {
			continue
		}

		g, err := s.Guild(fmt.Sprint(guildID))
		if err != nil {
			log.Printf("Error on getting guild for id %d: %v\n", guildID, err)
			continue
		}
		c, err := s.Channel(fmt.Sprint(channelID))
		if err != nil {
			log.Printf("Error on getting youtube channel for id %d: %v\n", channelID, err)
			continue
		}
		if c.GuildID != g.ID {
			log.Printf("Warning: tried to announce video in channel/%s/%s, but this channel (#%s) is from guild %s ('%s')\n", g.ID, c.ID, c.Name, c.GuildID, g.Name)
			continue
		}

		guilds = append(guilds, guild{
			guild:   g,
			channel: c,
			ping:    fmt.Sprintf("<@&%d>", roleID),
		})
	}
	return guilds, nil
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
