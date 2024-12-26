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
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// The remove subcommand. Used when executing the slash-command "/birthday remove".
type subcommandRemove struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption
}

// Constructor for subcommandremove, the struct for the slash-command "/birthday remove".
func (cmd Chat) subcommandRemove() subcommandRemove {
	subcommand := cmd.Interaction.ApplicationCommandData().Options[0]
	return subcommandRemove{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandRemove) handler() {
	authorID, err := strconv.ParseUint(cmd.user.ID, 10, 64)
	if err != nil {
		log.Printf("Error on parse author id of birthday command: %v\n", err)
		cmd.ReplyError()
		return
	}

	hasBDay, err := cmd.hasBirthday(authorID)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	e := util.AuthoredEmbed(cmd.Session, cmd.member, tp+"display")

	if !hasBDay {
		e.Description = lang.Get(tp+"msg.remove.not_found", lang.FallbackLang())
		e.Color = 0xFF0000
		cmd.ReplyHiddenEmbed(e)
		return
	}

	b, err := cmd.removeBirthday(authorID)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	e.Description = lang.Get(tp+"msg.remove", lang.FallbackLang())
	e.Color = 0x00FF00
	wasBefore := lang.Get(tp+"msg.remove.was", lang.FallbackLang())
	wasBefore = fmt.Sprintf(wasBefore, b)
	e.Fields = []*discordgo.MessageEmbedField{{
		Name:   wasBefore,
		Inline: true,
	}}

	if b.Visible {
		cmd.ReplyEmbed(e)
	} else {
		util.AddReplyHiddenField(e)
		cmd.ReplyHiddenEmbed(e)
	}
}
