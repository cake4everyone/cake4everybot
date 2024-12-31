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
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// The Chat (slash) command of the advent calendar package.
type Chat struct {
	adventcalendarBase
	ID string
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (Chat) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "base"),
		NameLocalizations:        util.TranslateLocalization(tp + "base"),
		Description:              lang.GetDefault(tp + "base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "base.description"),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "midnight",
				Description: "Midnight trigger",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "morning",
				Description: "Morning trigger",
			},
			{
				Type:                     discordgo.ApplicationCommandOptionSubCommand,
				Name:                     lang.GetDefault(tp + "option.draw"),
				NameLocalizations:        *util.TranslateLocalization(tp + "option.draw"),
				Description:              lang.GetDefault(tp + "option.draw.description"),
				DescriptionLocalizations: *util.TranslateLocalization(tp + "option.draw.description"),
			},
		},
	}
}

// Handle handles the functionality of a command
func (cmd Chat) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	cmd.member = i.Member
	cmd.user = i.User
	if i.Member != nil {
		cmd.user = i.Member.User
	} else if i.User != nil {
		cmd.member = &discordgo.Member{User: i.User}
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "midnight":
		cmd.ReplyDeferedHidden()
		Midnight(s)
		cmd.ReplyHidden("Midnight()")
		return
	case "morning":
		Post(s)
		cmd.ReplyHidden("Post()")
		return
	case lang.GetDefault(tp + "option.draw"):
		cmd.handleSubcommandDraw()
		return
	}

}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *Chat) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd Chat) GetID() string {
	return cmd.ID
}
