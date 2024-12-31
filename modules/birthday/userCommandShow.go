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

package birthday

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// UserShow represents a user command of the birthday package. It adds the ability to directly show
// a users birthday through a simple context click.
type UserShow struct {
	birthdayBase

	data discordgo.ApplicationCommandInteractionData
	ID   string
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd UserShow) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:              discordgo.UserApplicationCommand,
		Name:              lang.GetDefault(tp + "user.show.base"),
		NameLocalizations: util.TranslateLocalization(tp + "user.show.base"),
	}
}

// Handle handles the functionality of a command
func (cmd UserShow) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	cmd.member = i.Member
	cmd.user = i.User
	if i.Member != nil {
		cmd.user = i.Member.User
	}

	cmd.data = cmd.Interaction.ApplicationCommandData()
	cmd.handler()
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *UserShow) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd UserShow) GetID() string {
	return cmd.ID
}
