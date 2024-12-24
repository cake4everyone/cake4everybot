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
	"cake4everybot/database"
	"cake4everybot/util"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// The Component of the advent calendar package.
type Component struct {
	adventcalendarBase
	data discordgo.MessageComponentInteractionData
}

// Handle handles the functionality of a component.
func (c Component) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	c.member = i.Member
	c.user = i.User
	if i.Member != nil {
		c.user = i.Member.User
	} else if i.User != nil {
		c.member = &discordgo.Member{User: i.User}
	}
	c.data = i.MessageComponentData()

	ids := strings.Split(c.data.CustomID, ".")
	// pop the first level identifier
	util.ShiftL(ids)

	switch util.ShiftL(ids) {
	case "post":
		c.handlePost(s, ids)
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}
}

// ID returns the custom ID of the modal to identify the module
func (Component) ID() string {
	return "adventcalendar"
}

func (c *Component) handlePost(s *discordgo.Session, ids []string) {
	var (
		buttonYear  = util.ShiftL(ids)
		buttonMonth = util.ShiftL(ids)
		buttonDay   = util.ShiftL(ids)
	)
	timeValue := fmt.Sprintf("%s-%s-%s", buttonYear, buttonMonth, buttonDay)
	postTime, err := time.Parse(time.DateOnly, timeValue)
	if err != nil {
		log.Printf("ERROR: could not parse date: %s: %+v", timeValue, err)
		c.ReplyError()
		return
	}

	if now := time.Now(); now.Year() != postTime.Year() ||
		now.Month() != postTime.Month() ||
		now.Day() != postTime.Day() {
		c.ReplyHiddenSimpleEmbedf(0xFF0000, lang.GetDefault("module.adventcalendar.enter.invalid"))
		return
	}

	entry := database.GetGiveawayEntry("xmas", c.user.ID, database.AnnouncementPlatformDiscord, c.Interaction.GuildID)
	if entry.UserID != c.user.ID {
		log.Printf("ERROR: getEntry() returned with userID '%s' but want '%s'", entry.UserID, c.user.ID)
		c.ReplyError()
		return
	}
	if entry.LastEntry.Equal(postTime) {
		c.ReplyHiddenSimpleEmbedf(0x5865f2, lang.GetDefault("module.adventcalendar.enter.already_entered"), entry.Weight)
		return
	}

	entry = database.AddGiveawayWeight("xmas", c.user.ID, 1, database.AnnouncementPlatformDiscord, c.Interaction.GuildID)

	c.ReplyHiddenSimpleEmbedf(0x00FF00, lang.GetDefault("module.adventcalendar.enter.success"), entry.Weight)
}
