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
	"time"

	"github.com/bwmarrin/discordgo"
)

// The announce subcommand. Used when executing the slash-command "/birthday announce".
type subcommandAnnounce struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption
}

// Constructor for subcommandannounce, the struct for the slash-command "/birthday announce".
func (cmd Chat) subcommandAnnounce() subcommandAnnounce {
	subcommand := cmd.Interaction.ApplicationCommandData().Options[0]
	return subcommandAnnounce{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandAnnounce) handler() {
	now := time.Now()
	b, err := getBirthdaysDate(cmd.Interaction.GuildID, now.Day(), int(now.Month()))
	if err != nil {
		log.Printf("Error on announce birthday in guild %s: %v\n", cmd.Interaction.GuildID, err)
		cmd.ReplyError()
		return
	}

	e, n := birthdayAnnounceEmbed(cmd.Session, cmd.Interaction.GuildID, b)

	if n <= 0 {
		cmd.ReplyHiddenEmbed(e)
	} else {
		cmd.ReplyEmbed(e)
	}
}
