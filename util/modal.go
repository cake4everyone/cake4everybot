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

package util

import "github.com/bwmarrin/discordgo"

// CreateButtonComponent returns a simple button component with the specified configurations.
// Params:
//
//	tp                 // The translation prefix used to generate ids and labels
//	component_group_id // A group id to generate labels
//	id                 // Custom id to identify the button when pressed (automatically prefixed)
//	style              // Style of the button (see https://discord.com/developers/docs/interactions/message-components#button-object-button-styles)
//	Optional: emoji    // An emoji to put in the label, can be empty
func CreateButtonComponent(id, label string, style discordgo.ButtonStyle, emoji *discordgo.ComponentEmoji) discordgo.Button {
	return discordgo.Button{
		CustomID: id,
		Label:    label,
		Style:    style,
		Emoji:    emoji,
	}
}

// CreateURLButtonComponent returns a URL button component with the specified configurations.
// Params:
//
//	tp                 // The translation prefix used to generate ids and labels
//	component_group_id // A group id to generate labels
//	id                 // Custom id to generate labels
//	url                // The link to open when clicked
//	Optional: emoji    // An emoji to put in the label, can be empty
func CreateURLButtonComponent(id, label, url string, emoji *discordgo.ComponentEmoji) discordgo.Button {
	return discordgo.Button{
		CustomID: id,
		Label:    label,
		Emoji:    emoji,
		URL:      url,
	}
}

// CreateTextInputComponent returns a text input form for modals with the specified configurations.
// Params:
//
//	tp                 // The translation prefix used to generate ids and labels
//	component_group_id // A group id to generate labels
//	id                 // Custom id to identify the input field after submitting
//	style              // Single or multi line
//	required           // If this has to be not empty
//	minLength          // Minimum number of characters that has to be entered
//	maxLength          // Maximum number of characters that are able to be entered
func CreateTextInputComponent(id, label, placeholder, value string, style discordgo.TextInputStyle, requred bool, minLength, maxLength int) discordgo.TextInput {
	return discordgo.TextInput{
		CustomID:    id,
		Label:       label,
		Style:       style,
		Placeholder: placeholder,
		Value:       value,
		Required:    requred,
		MinLength:   minLength,
		MaxLength:   maxLength,
	}
}

// CreateChannelSelectMenuComponent returns a channel select menu form for modals with the specified
// configurations. Params:
//
//	id                 // Custom id to identify the input field after submitting
//	placeholder        // Placeholder text if nothing is selected
//	minValues          // Minimum number of items that must be choosen
//	maxValues          // Maximum number of items that can be choosen
//	channelTypes       // Channel types to include in the select menu
func CreateChannelSelectMenuComponent(id, placeholder string, minValues, maxValues int, channelTypes ...discordgo.ChannelType) discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:     id,
		MenuType:     discordgo.ChannelSelectMenu,
		Placeholder:  placeholder,
		MinValues:    &minValues,
		MaxValues:    maxValues,
		ChannelTypes: channelTypes,
	}
}

// CreateMentionableSelectMenuComponent returns a user and roles select menu form for modals with
// the specified configurations. Params:
//
//	id                 // Custom id to identify the input field after submitting
//	placeholder        // Placeholder text if nothing is selected
//	minValues          // Minimum number of items that must be choosen
//	maxValues          // Maximum number of items that can be choosen
func CreateMentionableSelectMenuComponent(id, placeholder string, minValues, maxValues int) discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    id,
		MenuType:    discordgo.MentionableSelectMenu,
		Placeholder: placeholder,
		MinValues:   &minValues,
		MaxValues:   maxValues,
	}
}

// CreateRoleSelectMenuComponent returns a roles select menu form for modals with the specified
// configurations. Params:
//
//	id                 // Custom id to identify the input field after submitting
//	placeholder        // Placeholder text if nothing is selected
//	minValues          // Minimum number of items that must be choosen
//	maxValues          // Maximum number of items that can be choosen
func CreateRoleSelectMenuComponent(id, placeholder string, minValues, maxValues int) discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    id,
		MenuType:    discordgo.RoleSelectMenu,
		Placeholder: placeholder,
		MinValues:   &minValues,
		MaxValues:   maxValues,
	}
}

// CreateStringSelectMenuComponent returns a string select menu form for modals with the specified
// configurations. Params:
//
//	id                 // Custom id to identify the input field after submitting
//	placeholder        // Placeholder text if nothing is selected
//	minValues          // Minimum number of items that must be choosen
//	maxValues          // Maximum number of items that can be choosen
//	options            // Choices for the select menu
func CreateStringSelectMenuComponent(id, placeholder string, minValues, maxValues int, options ...discordgo.SelectMenuOption) discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    id,
		MenuType:    discordgo.StringSelectMenu,
		Placeholder: placeholder,
		MinValues:   &minValues,
		MaxValues:   maxValues,
		Options:     options,
	}
}

// CreateUserSelectMenuComponent returns a user select menu form for modals with the specified
// configurations. Params:
//
//	id                 // Custom id to identify the input field after submitting
//	placeholder        // Placeholder text if nothing is selected
//	minValues          // Minimum number of items that must be choosen
//	maxValues          // Maximum number of items that can be choosen
func CreateUserSelectMenuComponent(id, placeholder string, minValues, maxValues int) discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    id,
		MenuType:    discordgo.UserSelectMenu,
		Placeholder: placeholder,
		MinValues:   &minValues,
		MaxValues:   maxValues,
	}
}
