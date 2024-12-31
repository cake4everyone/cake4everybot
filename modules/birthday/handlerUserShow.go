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
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
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
			log.Printf("Error on show birthday: %v\n", err)
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
		embed.Color = 0xFF0000
		cmd.ReplyHiddenEmbed(embed)
		return
	}

	var age string
	if b.Year > 0 {
		age = fmt.Sprintf(" (%d)", b.Age()+1)
	}
	embed.Fields = []*discordgo.MessageEmbedField{{
		Name: fmt.Sprintf("%s <t:%d:R>%s", b.String(), b.NextUnix(), age),
	}}
	embed.Color = 0x00FF00

	if b.Visible {
		cmd.ReplyEmbed(embed)
	} else {
		util.AddReplyHiddenField(embed)
		cmd.ReplyHiddenEmbed(embed)
	}
}
