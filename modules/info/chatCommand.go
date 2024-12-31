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

package info

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/status"
	"github.com/cake4everyone/cake4everybot/util"
	"github.com/spf13/viper"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => info
	tp = "discord.command.info."
)

// The Chat (slash) command of the info package. Simply prints a
// little infomation about the bot.
type Chat struct {
	util.InteractionUtil

	ID string
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd Chat) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "base"),
		NameLocalizations:        util.TranslateLocalization(tp + "base"),
		Description:              lang.GetDefault(tp + "base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "base.description"),
	}
}

// Handle handles the functionality of a command
func (cmd Chat) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}

	e := &discordgo.MessageEmbed{
		Title: lang.GetDefault(tp + "title"),
		Color: 0x00FF00,
	}
	util.AddEmbedField(e,
		lang.Get(tp+"start_time", lang.FallbackLang()),
		fmt.Sprintf("<t:%d:R>", status.GetStartTime().Unix()),
		true,
	)
	util.AddEmbedField(e,
		lang.Get(tp+"latency", lang.FallbackLang()),
		fmt.Sprintf("%dms", s.LastHeartbeatAck.Sub(s.LastHeartbeatSent).Milliseconds()),
		true,
	)
	version := fmt.Sprintf("v%s", viper.GetString("version"))
	versionURL := fmt.Sprintf("https://github.com/Kesuaheli/cake4everybotgo/releases/tag/%s", version)
	util.AddEmbedField(e,
		lang.Get(tp+"version", lang.FallbackLang()),
		fmt.Sprintf("[%s](%s)", version, versionURL),
		false,
	)
	util.SetEmbedFooter(s, tp+"display", e)

	cmd.ReplyEmbed(e)
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *Chat) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd Chat) GetID() string {
	return cmd.ID
}
