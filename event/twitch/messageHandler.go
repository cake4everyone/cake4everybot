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

package twitch

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

const tp string = "twitch.command."

// MessageHandler handles new messages from the twitch chat(s). It will be called on every new
// message.
func MessageHandler(t *twitchgo.Session, channel string, user *twitchgo.IRCUser, message, msgID string, tags twitchgo.IRCMessageTags) {
	//log.Printf("<%s@%s> %s", user.Nickname, channel, message)
}

// HandleGeneralCommand is the handler for any command in a twitch chat.
func HandleGeneralCommand(t *twitchgo.Session, channel string, user *twitchgo.IRCUser, args []string) {
	if len(args) == 0 || args[0] == "" {
		return
	}

	channel, _ = strings.CutPrefix(channel, "#")
	cmd, err := database.GetTwitchCommand(channel, args[0])
	if err != nil {
		log.Printf("Failed to get twitch command for '%s' in channel '%s': %v", args[0], channel, err)
		return
	}
	switch cmd.ResponseType {
	case database.TwitchCommandResponseChat:
		t.SendMessage(channel, cmd.Response)
	case database.TwitchCommandResponseMention:
		t.SendMessagef(channel, "@%s %s", user.Nickname, cmd.Response)
	case database.TwitchCommandResponseFunc:
		var respFunc twitchgo.IRCChannelCommandMessageCallback
		switch cmd.Response {
		case "giveaway_join":
			respFunc = handleGiveawayJoin
		case "giveaway_tickets":
			respFunc = handleGiveawayTickets
		case "giveaway_draw":
			respFunc = handleGiveawayDraw
		case "adv_dc":
			respFunc = handleAdvancedDC
		default:
			return
		}
		respFunc(t, channel, user, args[1:])
	}
}

// handleGiveawayJoin is the handler for a command in a twitch chat. This handler buys a giveaway ticket
// and removes the configured cost amount for a ticket.
func handleGiveawayJoin(t *twitchgo.Session, channel string, user *twitchgo.IRCUser, args []string) {
	channel, _ = strings.CutPrefix(channel, "#")
	const tp = tp + "join."

	p, err := database.NewGiveawayPrize(viper.GetString("event.twitch_giveaway.prizes"))
	if err != nil {
		log.Printf("Error reading prizes file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	if !p.HasPrizeAvailable() {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.no_prizes"), user.Nickname)
		return
	}
	if p.HasPrizeWon(user.Nickname) {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.won"), user.Nickname)
		return
	}
	entry := database.GetGiveawayEntry("tw11", user.Nickname, database.AnnouncementPlatformTwitch, channel) // FIXME: use channel ID instead of name
	if entry.UserID == "" {
		log.Printf("Error getting database giveaway entry: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	if entry.Weight >= 10 {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.max_tickets"), user.Nickname)
		return
	}

	data, err := os.ReadFile(viper.GetString("event.twitch_giveaway.times"))
	if os.IsNotExist(err) {
		data = []byte("{}")
	} else if err != nil {
		log.Printf("Error reading times file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	var times = map[string]time.Time{}
	err = json.Unmarshal(data, &times)
	if err != nil {
		log.Printf("Error parsing times file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	m := viper.GetDuration("event.twitch_giveaway.cooldown")
	next := times[user.Nickname].Add(m * time.Minute)
	cooldown := time.Until(next).Round(time.Second)

	if cooldown > time.Second {
		msgs := lang.GetSlice(tp+"msg.cooldown", lang.FallbackLang())
		var i int
		if len(msgs) >= 2 {
			rand.Shuffle(len(msgs), func(i, j int) {
				msgs[i], msgs[j] = msgs[j], msgs[i]
			})
			i = rand.Intn(len(msgs) - 1)
		}
		t.SendMessagef(channel, msgs[i], user.Nickname, cooldown.String())
		return
	}

	seChannel, err := se.GetChannel(channel)
	if err != nil {
		log.Printf("Error getting streamelements channel '%s': %v", channel, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	sePoints, err := se.GetPoints(seChannel.ID, user.Nickname)
	if err != nil {
		log.Printf("Error getting streamelements points '%s(%s)/%s' : %v", seChannel.ID, channel, user.Nickname, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	joinCost := viper.GetInt("event.twitch_giveaway.ticket_cost")
	if sePoints.Points < joinCost {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.too_few_points"), user.Nickname, sePoints.Points, joinCost-sePoints.Points, joinCost)
		return
	}
	entry = database.AddGiveawayWeight("tw11", user.Nickname, 1, database.AnnouncementPlatformTwitch, channel) // FIXME: use channel ID instead of name
	if entry.UserID == "" {
		log.Printf("Error getting database giveaway entry: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	times[user.Nickname] = time.Now()
	data, err = json.Marshal(times)
	if err != nil {
		log.Printf("Error marshaling times file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	err = os.WriteFile(viper.GetString("event.twitch_giveaway.times"), data, 0644)
	if err != nil {
		log.Printf("Error writing times file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	err = se.AddPoints(seChannel.ID, user.Nickname, -joinCost)
	if err != nil {
		log.Printf("Error adding points for '%s(%s)/%s/-%d': %v", seChannel.ID, channel, user.Nickname, joinCost, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	t.SendMessagef(channel, lang.GetDefault(tp+"msg.success"), user.Nickname, joinCost, entry.Weight, sePoints.Points-joinCost)
}

// handleGiveawayTickets is the handler for the tickets command in a twitch chat. This handler simply
// prints the users amount of tickets
func handleGiveawayTickets(t *twitchgo.Session, channel string, source *twitchgo.IRCUser, args []string) {
	channel, _ = strings.CutPrefix(channel, "#")
	const tp = tp + "tickets."

	var userID string = source.Nickname
	if len(args) >= 1 {
		if s, _ := strings.CutPrefix(args[0], "@"); s != "" {
			userID = strings.ToLower(s)
		}
	}

	p, err := database.NewGiveawayPrize(viper.GetString("event.twitch_giveaway.prizes"))
	if err != nil {
		log.Printf("Error reading prizes file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	if p.HasPrizeWon(userID) {
		if source.Nickname == userID {
			t.SendMessagef(channel, lang.GetDefault(tp+"msg.won"), source.Nickname)
		} else {
			t.SendMessagef(channel, lang.GetDefault(tp+"msg.won.user"), source.Nickname, userID)
		}
		return
	}

	entry := database.GetGiveawayEntry("tw11", userID, database.AnnouncementPlatformTwitch, channel) // FIXME: use channel ID instead of name
	if entry.Weight >= 10 {
		if source.Nickname == userID {
			t.SendMessagef(channel, lang.GetDefault(tp+"msg.max_tickets"), source.Nickname)
		} else {
			t.SendMessagef(channel, lang.GetDefault(tp+"msg.max_tickets.user"), source.Nickname, userID)
		}
		return
	}
	if source.Nickname != userID {
		if entry.Weight == 0 {
			t.SendMessagef(channel, lang.GetDefault(tp+"msg.num.0.user"), source.Nickname, userID)
		} else {
			t.SendMessagef(channel, lang.GetDefault(tp+"msg.num.user"), source.Nickname, userID, entry.Weight)
		}
		return
	}
	var msg string
	if entry.Weight == 0 {
		msg = fmt.Sprintf(lang.GetDefault(tp+"msg.num.0"), source.Nickname)
	} else {
		msg = fmt.Sprintf(lang.GetDefault(tp+"msg.num"), source.Nickname, entry.Weight)
	}

	var curPoints int
	seChannel, err := se.GetChannel(channel)
	if err != nil {
		log.Printf("Error on getting SE channel: %v", err)
		goto skipPoints
	}
	if sePoints, err := se.GetPoints(seChannel.ID, userID); err != nil {
		log.Printf("Error on getting SE points: %v", err)
		goto skipPoints
	} else {
		curPoints = sePoints.Points
	}

	if joinCost := viper.GetInt("event.twitch_giveaway.ticket_cost"); joinCost > curPoints {
		msg += " " + fmt.Sprintf(lang.GetDefault(tp+"msg.extra.need_points"), joinCost-curPoints)
	} else {
		msg += " " + lang.GetDefault(tp+"msg.extra.can_buy")
	}
skipPoints:

	data, err := os.ReadFile(viper.GetString("event.twitch_giveaway.times"))
	if os.IsNotExist(err) {
		data = []byte("{}")
	} else if err != nil {
		log.Printf("Error reading times file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	var times = map[string]time.Time{}
	err = json.Unmarshal(data, &times)
	if err != nil {
		log.Printf("Error parsing times file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	m := viper.GetDuration("event.twitch_giveaway.cooldown")
	next := times[userID].Add(m * time.Minute)
	cooldown := time.Until(next).Round(time.Second)

	if cooldown > 3*time.Second {
		msg += " " + fmt.Sprintf(lang.GetDefault(tp+"msg.extra.cooldown"), cooldown.String())
	}

	t.SendMessage(channel, msg)
}

// handleGiveawayDraw is the handler for the draw command in a twitch chat. This handler selects a random
// winner and removes their tickets.
func handleGiveawayDraw(t *twitchgo.Session, channel string, user *twitchgo.IRCUser, args []string) {
	channel, _ = strings.CutPrefix(channel, "#")
	const tp = tp + "draw."

	//only accept broadcaster
	if channel != user.Nickname {
		return
	}

	p, err := database.NewGiveawayPrize(viper.GetString("event.twitch_giveaway.prizes"))
	if err != nil {
		log.Printf("Error reading prizes file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	prize, ok := p.GetNextPrize()
	if !ok {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.no_prizes"), user.Nickname)
		return
	}

	winner, totalTickets := database.DrawGiveawayWinner(database.GetAllGiveawayEntries("tw11", database.AnnouncementPlatformTwitch, channel)) // FIXME: use channel ID instead of name
	if totalTickets == 0 {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.no_entries"), user.Nickname)
		return
	}

	t.SendMessagef(channel, lang.GetDefault(tp+"msg.winner"), winner.UserID, prize.Name, winner.Weight, float64(winner.Weight*100)/float64(totalTickets))

	err = database.DeleteGiveawayEntry(winner.UserID, database.AnnouncementPlatformTwitch, channel) // FIXME: use channel ID instead of name
	if err != nil {
		log.Printf("Error deleting database giveaway entry: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	prize.Winner = winner.UserID
	err = p.SaveFile()
	if err != nil {
		log.Printf("Error saving prizes file: %v", err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
}

func handleAdvancedDC(t *twitchgo.Session, channel string, user *twitchgo.IRCUser, args []string) {
	channel, _ = strings.CutPrefix(channel, "#")
	if channel == "kesuaheli" && len(args) > 0 && args[0] != "" {
		switch args[0] {
		case "pose", "._.":
			t.SendMessage(channel, "Trete dem Bored-Face Discord Server bei: discord.gg/xGHjcMh2JA")
		case "dmm", "mapmaking":
			t.SendMessage(channel, "Trete dem Deutschen Mapmaking bei: discord.gg/gzSgJZJsdp")
		case "wiki":
			t.SendMessage(channel, "Trete dem Deutschen Minecraft Wiki Discord bei: discord.gg/F75vfpd")
		}
		return
	}
	t.SendMessagef(channel, lang.GetDefault("twitch.command.advanced_dc.default"), viper.GetString("discord.invite"))
}
