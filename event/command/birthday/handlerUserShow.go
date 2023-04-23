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
	"fmt"
	"log"
	"strconv"

	"cake4everybot/data/lang"
	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
)

func (cmd UserShow) handler() {
	targetID, err := strconv.ParseUint(cmd.data.TargetID, 10, 64)
	if err != nil {
		log.Printf("Error on parse target id of birthday user show command: %v\n", err)
		cmd.ReplyError()
		return
	}
	b := birthdayEntry{ID: targetID}

	target := cmd.data.Resolved.Members[cmd.data.TargetID]
	target.User = cmd.data.Resolved.Users[cmd.data.TargetID]

	hasBDay, err := cmd.hasBirthday(b.ID)
	if err != nil {
		log.Printf("Error on show birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	self := cmd.user.ID == cmd.data.TargetID

	if hasBDay {
		err = cmd.getBirthday(&b)
		if err != nil {
			log.Printf("Error on show birthday: %v", err)
			cmd.ReplyError()
			return
		}
		//pretend to have no birthday when its not visible
		hasBDay = self || b.Visible
	}

	embed := util.AuthoredEmbed(cmd.Session, target, tp+"display")

	if !hasBDay {
		if self {
			format := lang.GetDefault(tp + "msg.no_entry")
			mentionCmd := util.MentionCommand(tp+"base", tp+"option.set")
			embed.Description = fmt.Sprintf(format, mentionCmd)
		} else {
			format := lang.GetDefault(tp + "msg.no_entry.user")
			embed.Description = fmt.Sprintf(format, target.Mention())
		}
		cmd.ReplyHiddenEmbed(embed)
		return
	}

	embed.Fields = []*discordgo.MessageEmbedField{{
		Name: b.String(),
	},
	}

	if hasBDay && self && !b.Visible {
		cmd.ReplyHiddenEmbed(embed)
		return
	}

	cmd.ReplyEmbed(embed)
}
