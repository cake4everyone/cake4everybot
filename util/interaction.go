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
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
)

// InteractionUtil is a helper for discords application interactions. It add useful methods for
// simpler and faster coding.
type InteractionUtil struct {
	Session         *discordgo.Session
	Interaction     *discordgo.InteractionCreate
	response        *discordgo.InteractionResponse
	acknowledged    bool
	followupMessage *discordgo.Message
}

func (i *InteractionUtil) respond() {
	if i.response.Type != discordgo.InteractionResponseDeferredChannelMessageWithSource && // deferred responses dont need contents
		i.response.Type != discordgo.InteractionResponseDeferredMessageUpdate &&
		i.response.Type != discordgo.InteractionApplicationCommandAutocompleteResult && // autocomplete responses dont need contents
		i.response.Data.Content == "" &&
		len(i.response.Data.Embeds) == 0 &&
		len(i.response.Data.Components) == 0 {
		log.Printf("ERROR: Reply called without contens! Need at least one of Content, Embed, Component.\n%s", debug.Stack())
		i.ReplyError()
		return
	}

	var err error
	if i.acknowledged {
		data := &discordgo.WebhookParams{
			AllowedMentions: i.response.Data.AllowedMentions,
			Components:      i.response.Data.Components,
			Content:         i.response.Data.Content,
			Embeds:          i.response.Data.Embeds,
			Files:           i.response.Data.Files,
			Flags:           i.response.Data.Flags,
			TTS:             i.response.Data.TTS,
		}
		if i.response.Data.Attachments != nil {
			data.Attachments = *i.response.Data.Attachments
		}
		i.followupMessage, err = i.Session.FollowupMessageCreate(i.Interaction.Interaction, true, data)
		if err != nil {
			log.Printf("ERROR: could not send follow up message: %+v\n%s", err, debug.Stack())
		}
		return
	}
	err = i.Session.InteractionRespond(i.Interaction.Interaction, i.response)
	if err != nil {
		log.Printf("ERROR could not send command response: %+v\n%s", err, debug.Stack())
		return
	}
	i.acknowledged = true
}

func (i *InteractionUtil) respondMessage(update, deferred bool) (sucess bool) {
	rType := discordgo.InteractionResponseChannelMessageWithSource
	if update {
		if i.Interaction.Type != discordgo.InteractionMessageComponent {
			log.Printf("ERROR: Got message update on interaction type '%s', but is only allowed on %s", i.Interaction.Type.String(), discordgo.InteractionMessageComponent.String())
			i.ReplyError()
			return false
		}
		if deferred {
			rType = discordgo.InteractionResponseDeferredMessageUpdate
		} else {
			rType = discordgo.InteractionResponseUpdateMessage
		}
	} else if !update && deferred {
		rType = discordgo.InteractionResponseDeferredChannelMessageWithSource
	}

	i.response = &discordgo.InteractionResponse{
		Type: rType,
		Data: &discordgo.InteractionResponseData{},
	}
	return true
}

// ReplyError sends a simple message to the user to indicate, that something failed or unexpected
// happened during the execution of the interaction.
func (i *InteractionUtil) ReplyError() {
	i.ReplyHidden("Somthing went wrong :(")
}

// ReplyErrorUpdate is like ReplyError but made for an update message for components.
func (i *InteractionUtil) ReplyErrorUpdate() {
	i.ReplyUpdate("Somthing went wrong :(")
}

// ReplyDefered points out to the user that the bots answer could take some time. It also allows the
// bot to extend this interaction for up to 15 minutes.
func (i *InteractionUtil) ReplyDefered() {
	i.respondMessage(false, true)
	i.respond()
}

// ReplyDeferedHidden is like ReplyDefered but also ephemeral.
func (i *InteractionUtil) ReplyDeferedHidden() {
	i.respondMessage(false, true)
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplyDeferedUpdate is like ReplyDefered but make for an update for components.
func (i *InteractionUtil) ReplyDeferedUpdate() {
	i.respondMessage(true, true)
	i.respond()
}

// ReplyDeferedHiddenUpdate is like ReplyDeferedHidden but make for an update for components.
func (i *InteractionUtil) ReplyDeferedHiddenUpdate() {
	i.respondMessage(true, true)
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// Reply prints the given message as reply to the user who executes the command.
func (i *InteractionUtil) Reply(message string) {
	i.respondMessage(false, false)
	i.response.Data.Content = message
	i.respond()
}

// Replyf formats according to a format specifier and prints the result as reply to the user who
// executes the command.
func (i *InteractionUtil) Replyf(format string, a ...any) {
	i.Reply(fmt.Sprintf(format, a...))
}

// ReplyUpdate is like Reply but make for an update for components.
func (i *InteractionUtil) ReplyUpdate(message string) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data.Content = message
	i.respond()
}

// ReplyUpdatef is like Replyf but make for an update for components.
func (i *InteractionUtil) ReplyUpdatef(format string, a ...any) {
	i.ReplyUpdate(fmt.Sprintf(format, a...))
}

// ReplyHidden prints the given message as ephemeral reply to the user who executes the command.
func (i *InteractionUtil) ReplyHidden(message string) {
	i.respondMessage(false, false)
	i.response.Data.Content = message
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplyHiddenf formats according to a format specifier and prints the result as ephemeral reply to
// the user who executes the command.
func (i *InteractionUtil) ReplyHiddenf(format string, a ...any) {
	i.ReplyHidden(fmt.Sprintf(format, a...))
}

// ReplyEmbed prints the given embeds as reply to the user who executes the command.
func (i *InteractionUtil) ReplyEmbed(embeds ...*discordgo.MessageEmbed) {
	i.respondMessage(false, false)
	i.response.Data.Embeds = embeds
	i.respond()
}

// ReplyEmbedUpdate is like ReplyEmbed but make for an update for components.
func (i *InteractionUtil) ReplyEmbedUpdate(embeds ...*discordgo.MessageEmbed) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data.Embeds = embeds
	i.respond()
}

// ReplyHiddenEmbed prints the given embeds as ephemeral reply to the user who executes the command.
func (i *InteractionUtil) ReplyHiddenEmbed(embeds ...*discordgo.MessageEmbed) {
	i.respondMessage(false, false)
	i.response.Data.Embeds = embeds
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplyHiddenEmbedUpdate is like [InteractionUtil.ReplyHiddenEmbed] but made for an update for
// components.
func (i *InteractionUtil) ReplyHiddenEmbedUpdate(embeds ...*discordgo.MessageEmbed) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data.Embeds = embeds
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplySimpleEmbed is a shortcut for replying with a simple embed that only contains a single text
// and has a color.
func (i *InteractionUtil) ReplySimpleEmbed(color int, content string) {
	i.ReplyEmbed(SimpleEmbed(color, content)...)
}

// ReplySimpleEmbedf formats according to a format specifier and is a shortcut for replying with a
// simple embed that only contains a single text and has a color.
func (i *InteractionUtil) ReplySimpleEmbedf(color int, format string, a ...any) {
	i.ReplyEmbed(SimpleEmbedf(color, format, a...)...)
}

// ReplySimpleEmbedUpdate is like ReplySimpleEmbed but make for an update for components.
func (i *InteractionUtil) ReplySimpleEmbedUpdate(color int, content string) {
	i.ReplyEmbedUpdate(SimpleEmbed(color, content)...)
}

// ReplySimpleEmbedUpdatef is like ReplySimpleEmbedf but make for an update for components.
func (i *InteractionUtil) ReplySimpleEmbedUpdatef(color int, format string, a ...any) {
	i.ReplyEmbedUpdate(SimpleEmbedf(color, format, a...)...)
}

// ReplyHiddenSimpleEmbed is like ReplySimpleEmbed but also ephemeral.
func (i *InteractionUtil) ReplyHiddenSimpleEmbed(color int, content string) {
	i.ReplyHiddenEmbed(SimpleEmbed(color, content)...)
}

// ReplyHiddenSimpleEmbedf is like ReplySimpleEmbedf but also ephemeral.
func (i *InteractionUtil) ReplyHiddenSimpleEmbedf(color int, format string, a ...any) {
	i.ReplyHiddenEmbed(SimpleEmbedf(color, format, a...)...)
}

// ReplyHiddenSimpleEmbedUpdate is like [InteractionUtil.ReplyHiddenSimpleEmbed] but made for an
// update for components.
func (i *InteractionUtil) ReplyHiddenSimpleEmbedUpdate(color int, content string) {
	i.ReplyHiddenEmbedUpdate(SimpleEmbed(color, content)...)
}

// ReplyHiddenSimpleEmbedUpdatef is like [InteractionUtil.ReplyHiddenSimpleEmbedf] but made for an
// update for components.
func (i *InteractionUtil) ReplyHiddenSimpleEmbedUpdatef(color int, format string, a ...any) {
	i.ReplyHiddenSimpleEmbedUpdate(color, fmt.Sprintf(format, a...))
}

// ReplyComponents sends a message along with the provied message components.
func (i *InteractionUtil) ReplyComponents(components []discordgo.MessageComponent, message string) {
	i.respondMessage(false, false)
	i.response.Data.Content = message
	i.response.Data.Components = components
	i.respond()
}

// ReplyComponentsf formats according to a format specifier and sends the result along with the
// provied message components.
func (i *InteractionUtil) ReplyComponentsf(components []discordgo.MessageComponent, format string, a ...any) {
	i.ReplyComponents(components, fmt.Sprintf(format, a...))
}

// ReplyComponentsUpdate is like ReplyComponents but make for an update for components.
func (i *InteractionUtil) ReplyComponentsUpdate(components []discordgo.MessageComponent, message string) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data.Content = message
	i.response.Data.Components = components
	i.respond()
}

// ReplyComponentsUpdatef is like ReplyComponentsf but make for an update for components.
func (i *InteractionUtil) ReplyComponentsUpdatef(components []discordgo.MessageComponent, format string, a ...any) {
	i.ReplyComponentsUpdate(components, fmt.Sprintf(format, a...))
}

// ReplyComponentsHidden sends an ephemeral message along with the provided message components.
func (i *InteractionUtil) ReplyComponentsHidden(components []discordgo.MessageComponent, message string) {
	i.respondMessage(false, false)
	i.response.Data.Content = message
	i.response.Data.Components = components
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplyComponentsHiddenf is like ReplyComponentsf but sends an ephemeral message.
func (i *InteractionUtil) ReplyComponentsHiddenf(components []discordgo.MessageComponent, format string, a ...any) {
	i.ReplyComponentsHidden(components, fmt.Sprintf(format, a...))
}

// ReplyComponentsEmbed sends one or more embeds along with the provied message components.
func (i *InteractionUtil) ReplyComponentsEmbed(components []discordgo.MessageComponent, embeds ...*discordgo.MessageEmbed) {
	i.respondMessage(false, false)
	i.response.Data.Embeds = embeds
	i.response.Data.Components = components
	i.respond()
}

// ReplyComponentsEmbedUpdate is like [InteractionUtil.ReplyComponentsEmbed] but made for an update for components.
func (i *InteractionUtil) ReplyComponentsEmbedUpdate(components []discordgo.MessageComponent, embeds ...*discordgo.MessageEmbed) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data.Embeds = embeds
	i.response.Data.Components = components
	i.respond()
}

// ReplyComponentsHiddenEmbed sends the given embeds as ephemeral reply along with the provided message
// components.
func (i *InteractionUtil) ReplyComponentsHiddenEmbed(components []discordgo.MessageComponent, embeds ...*discordgo.MessageEmbed) {
	i.respondMessage(false, false)
	i.response.Data.Embeds = embeds
	i.response.Data.Components = components
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplyComponentsHiddenEmbedUpdate is like [InteractionUtil.ReplyComponentsHiddenEmbed] but made for an update for components.
func (i *InteractionUtil) ReplyComponentsHiddenEmbedUpdate(components []discordgo.MessageComponent, embeds ...*discordgo.MessageEmbed) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data.Embeds = embeds
	i.response.Data.Components = components
	i.response.Data.Flags = discordgo.MessageFlagsEphemeral
	i.respond()
}

// ReplyComponentsSimpleEmbed sends an embed message along with the provied message components.
func (i *InteractionUtil) ReplyComponentsSimpleEmbed(components []discordgo.MessageComponent, color int, content string) {
	i.ReplyComponentsEmbed(components, SimpleEmbed(color, content)...)
}

// ReplyComponentsSimpleEmbedf is like [InteractionUtil.ReplyComponentsSimpleEmbed] but formats the
// embed content according to a format specifier.
func (i *InteractionUtil) ReplyComponentsSimpleEmbedf(components []discordgo.MessageComponent, color int, format string, a ...any) {
	i.ReplyComponentsSimpleEmbed(components, color, fmt.Sprintf(format, a...))
}

// ReplyComponentsSimpleEmbedUpdate is like [InteractionUtil.ReplyComponentsSimpleEmbed] but made for an update for components.
func (i *InteractionUtil) ReplyComponentsSimpleEmbedUpdate(components []discordgo.MessageComponent, color int, content string) {
	i.ReplyComponentsEmbedUpdate(components, SimpleEmbed(color, content)...)
}

// ReplyComponentsSimpleEmbedUpdatef is like [InteractionUtil.ReplyComponentsSimpleEmbedf] but made for an update for components.
func (i *InteractionUtil) ReplyComponentsSimpleEmbedUpdatef(components []discordgo.MessageComponent, color int, format string, a ...any) {
	i.ReplyComponentsSimpleEmbedUpdate(components, color, fmt.Sprintf(format, a...))
}

// ReplyComponentsHiddenSimpleEmbed is like [InteractionUtil.ReplyHiddenSimpleEmbed] but sends the
// embed message along with the provied message components.
func (i *InteractionUtil) ReplyComponentsHiddenSimpleEmbed(components []discordgo.MessageComponent, color int, content string) {
	i.ReplyComponentsHiddenEmbed(components, SimpleEmbed(color, content)...)
}

// ReplyComponentsHiddenSimpleEmbedf is like [InteractionUtil.ReplyHiddenSimpleEmbedf] but sends the
// embed message along with the provied message components.
func (i *InteractionUtil) ReplyComponentsHiddenSimpleEmbedf(components []discordgo.MessageComponent, color int, format string, a ...any) {
	i.ReplyComponentsHiddenSimpleEmbed(components, color, fmt.Sprintf(format, a...))
}

// ReplyComponentsHiddenSimpleEmbedUpdate is like [InteractionUtil.ReplyComponentsHiddenSimpleEmbed] but made for an update for components.
func (i *InteractionUtil) ReplyComponentsHiddenSimpleEmbedUpdate(components []discordgo.MessageComponent, color int, content string) {
	i.ReplyComponentsHiddenEmbedUpdate(components, SimpleEmbed(color, content)...)
}

// ReplyComponentsHiddenSimpleEmbedUpdatef is like [InteractionUtil.ReplyComponentsHiddenSimpleEmbedf] but made for an update for components.
func (i *InteractionUtil) ReplyComponentsHiddenSimpleEmbedUpdatef(components []discordgo.MessageComponent, color int, format string, a ...any) {
	i.ReplyComponentsHiddenSimpleEmbedUpdate(components, color, fmt.Sprintf(format, a...))
}

// ReplyInteractionEdit edits the original interaction message with the given data.
//
// It is intended to be used after a deferred response. If the previous response was not deferred,
// ReplyInteractionEdit will reply with an error message.
//
// If delete is true, the deferred response will be deleted. If there was no follow up message, a
// simple success message will be sent as a reply and instantly deleted.
func (i *InteractionUtil) ReplyInteractionEdit(data *discordgo.InteractionResponseData, delete bool) {
	if i.response == nil || i.response.Type != discordgo.InteractionResponseDeferredMessageUpdate && i.response.Type != discordgo.InteractionResponseDeferredChannelMessageWithSource {
		log.Printf("ERROR: ReplyEdit called without a previous deferred response. Preferably you should use [ReplyComplex()].\n%s", debug.Stack())
		i.ReplyError()
		return
	}

	messageEdit := MessageComplexEdit(data, i.Interaction.ChannelID, i.Interaction.Message.ID)
	_, err := i.Session.ChannelMessageEditComplex(messageEdit)
	if err != nil {
		log.Printf("ERROR: could not edit message: %+v\n%s", err, debug.Stack())
		i.ReplyError()
		return
	}

	if delete {
		if i.followupMessage == nil {
			i.Reply("-# ✅")
		}
		err = i.Session.FollowupMessageDelete(i.Interaction.Interaction, i.followupMessage.ID)
		if err != nil {
			log.Printf("ERROR: could not delete follow up message: %+v\n%s", err, debug.Stack())
			i.ReplyError()
		}
		i.followupMessage = nil
		return
	}
}

// ReplyComplex sends the given interaction response data to the user.
func (i *InteractionUtil) ReplyComplex(data *discordgo.InteractionResponseData) {
	i.respondMessage(false, false)
	i.response.Data = data
	i.respond()
}

// ReplyComplexUpdate is like [InteractionUtil.ReplyComplex] but made for an update for components.
func (i *InteractionUtil) ReplyComplexUpdate(data *discordgo.InteractionResponseData) {
	if !i.respondMessage(true, false) {
		return
	}
	i.response.Data = data
	i.respond()
}

// ReplyAutocomplete returns the given choices to the user. When this is called on an interaction
// type outside form an applicationCommandAutocomplete nothing will happen.
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

// ReplyModal displays a modal (popup) with the specified components to the user.
//
// Params:
//
//	id // To identify the modal when parsing the interaction event
//	title // Displayed title of the modal
//	components // One or more message components to display in this modal
func (i *InteractionUtil) ReplyModal(id, title string, components ...discordgo.MessageComponent) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   id,
			Title:      title,
			Components: components,
		},
	}
	i.respond()
}

// RequireOriginalAuthor checks if the user who executed the current interaction
// is the same as the original author of the interaction.
func (i *InteractionUtil) RequireOriginalAuthor() bool {
	originalAuthor := OriginalAuthor(i.Interaction.Message)
	// the user who executed the current interaction
	var user *discordgo.User = i.Interaction.User
	if user == nil {
		user = i.Interaction.Member.User
	}
	if originalAuthor.ID != user.ID {
		i.ReplyHiddenf(lang.GetDefault("discord.command.generic.msg.error.not_author"), originalAuthor.Mention())
		return false
	}
	return true
}
