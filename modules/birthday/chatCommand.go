// Copyright 2022-2023 Kesuaheli
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

package birthday

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// The Chat (slash) command of the birthday package. Has a few sub commands and options to use all
// features through a single chat command.
type Chat struct {
	birthdayBase

	ID string
}

type subcommand interface {
	handler()
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd Chat) AppCmd() *discordgo.ApplicationCommand {
	options := []*discordgo.ApplicationCommandOption{
		subCommandSet(),
		subCommandRemove(),
		subCommandList(),
		subCommandAnnounce(),
	}

	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "base"),
		NameLocalizations:        util.TranslateLocalization(tp + "base"),
		Description:              lang.GetDefault(tp + "base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "base.description"),
		Options:                  options,
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

	subcommandName := i.ApplicationCommandData().Options[0].Name
	var sub subcommand

	switch subcommandName {
	case lang.GetDefault(tp + "option.set"):
		sub = cmd.subcommandSet()
	case lang.GetDefault(tp + "option.remove"):
		sub = cmd.subcommandRemove()
	case lang.GetDefault(tp + "option.list"):
		sub = cmd.subcommandList()
	case lang.GetDefault(tp + "option.announce"):
		sub = cmd.subcommandAnnounce()
	default:
		return
	}

	sub.handler()
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *Chat) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd Chat) GetID() string {
	return cmd.ID
}
