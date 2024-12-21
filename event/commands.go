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

package event

import (
	"cake4everybot/event/command"
	"cake4everybot/event/component"
	"cake4everybot/event/modal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()
		if cmd, ok := command.CommandMap[data.Name]; ok {
			cmd.Handle(s, i)
		}
	// TODO: Add seperate handler for autocomplete
	//case discordgo.InteractionApplicationCommandAutocomplete:
	//	data := i.ApplicationCommandData()
	//	if cmd, ok := command.CommandMap[data.Name]; ok {
	//		cmd.HandleAutocomplete(s, i)
	//	}
	case discordgo.InteractionMessageComponent:
		data := i.MessageComponentData()
		if c, ok := component.ComponentMap[strings.Split(data.CustomID, ".")[0]]; ok {
			c.Handle(s, i)
		} else {
			log.Printf("got component interaction from unknown module '%s' (full id '%s')", strings.Split(data.CustomID, ".")[0], data.CustomID)
		}

	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()
		if m, ok := modal.ModalMap[strings.Split(data.CustomID, ".")[0]]; ok {
			m.HandleModal(s, i)
		}
	}
}
