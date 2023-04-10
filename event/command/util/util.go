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

package util

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type InteractionUtil struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	response    *discordgo.InteractionResponse
}

// Replyf formats according to a format specifier
// and prints the result as reply to the user who
// executes the command.
func (i *InteractionUtil) Replyf(format string, a ...any) {
	i.Reply(fmt.Sprintf(format, a...))
}

// Prints the given message as reply to the
// user who executes the command.
func (i *InteractionUtil) Reply(message string) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
	i.respond()
}

// Prints the given embeds as reply to the
// user who executes the command.
func (i *InteractionUtil) ReplyEmbed(embeds ...*discordgo.MessageEmbed) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	}
	i.respond()
}

// Replyf formats according to a format specifier
// and prints the result as emphemral reply to
// the user who executes the command.
func (i *InteractionUtil) ReplyHiddenf(format string, a ...any) {
	i.ReplyHidden(fmt.Sprintf(format, a...))
}

// Prints the given message as emphemral reply
// to the user who executes the command.
func (i *InteractionUtil) ReplyHidden(message string) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	i.respond()
}

// Prints the given embeds as emphemral reply
// to the user who executes the command.
func (i *InteractionUtil) ReplyHiddenEmbed(embeds ...*discordgo.MessageEmbed) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}
	i.respond()
}

// ReplyAutocomplete returns the given choices to
// the user. When this is called on an interaction
// type outside form an applicationCommandAutocomplete
// nothing will happen.
func (i *InteractionUtil) ReplyAutocomplete(choices []*discordgo.ApplicationCommandOptionChoice) {
	if i.Interaction.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}
	i.respond()
}

func (i *InteractionUtil) ReplyError() {
	i.ReplyHidden("Somthing went wrong :(")
}

func (i *InteractionUtil) respond() {
	err := i.Session.InteractionRespond(i.Interaction.Interaction, i.response)
	if err != nil {
		fmt.Printf("Error while sending command response: %v", err)
	}
}
