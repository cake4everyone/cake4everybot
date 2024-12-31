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

package event

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/database"
)

func addVoiceStateListeners(s *discordgo.Session) {
	handler := func(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {

		isAfk := isInAfkVoiceChannel(s, e)

		if isVoiceChannelUpdate(e) && (e.ChannelID == "" || isAfk) {
			setNoMicPermission(s, e, false)
		} else if isVoiceChannelUpdate(e) {
			setNoMicPermission(s, e, true)
		}
	}

	s.AddHandler(handler)
}

func isVoiceChannelUpdate(e *discordgo.VoiceStateUpdate) bool {
	return e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID != e.ChannelID
}

func isInAfkVoiceChannel(s *discordgo.Session, e *discordgo.VoiceStateUpdate) (isAfk bool) {
	guild, err := s.Guild(e.GuildID)
	if err != nil {
		log.Printf("ERROR: on join afk vc: %v\n", err)
		return false
	}

	return e.ChannelID == guild.AfkChannelID
}

func setNoMicPermission(s *discordgo.Session, e *discordgo.VoiceStateUpdate, state bool) {
	var NoMicChannelID uint64
	err := database.QueryRow("SELECT no_mic_id FROM guilds WHERE id = ?", e.GuildID).Scan(&NoMicChannelID)
	if err != nil {
		log.Printf("Error on no mic permission database call: %v\n", err)
	}
	if state {
		err = s.ChannelPermissionSet(fmt.Sprint(NoMicChannelID),
			e.Member.User.ID,
			discordgo.PermissionOverwriteTypeMember,
			discordgo.PermissionViewChannel,
			0)
	} else {
		err = s.ChannelPermissionDelete(fmt.Sprint(NoMicChannelID), e.Member.User.ID)
	}

	if err != nil {
		log.Printf("Error on no mic permission: %v\n", err)
		return
	}
}
