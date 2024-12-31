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
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

var minValueOne = 1.0
var minValueZero = 0.0

func subCommandSet() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		commandOptionSetDay(),
		commandOptionSetMonth(),
		commandOptionSetYear(),
		commandOptionSetVisible(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.set"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set"),
		Description:              lang.GetDefault(tp + "option.set.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.description"),
		Options:                  options,
	}
}

func commandOptionSetDay() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.day"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.day"),
		Description:              lang.GetDefault(tp + "option.set.option.day.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.day.description"),
		Required:                 true,
		Autocomplete:             true,
		MinValue:                 &minValueOne,
		MaxValue:                 31,
	}
}

func commandOptionSetMonth() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.month"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.month"),
		Description:              lang.GetDefault(tp + "option.set.option.month.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.month.description"),
		Required:                 true,
		Autocomplete:             true,
		MinValue:                 &minValueOne,
		MaxValue:                 12,
	}
}

func commandOptionSetYear() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.year"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.year"),
		Description:              lang.GetDefault(tp + "option.set.option.year.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.year.description"),
		Autocomplete:             true,
	}
}

func commandOptionSetVisible() *discordgo.ApplicationCommandOption {
	choices := []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:              lang.GetDefault("discord.command.generic.yes"),
			NameLocalizations: *util.TranslateLocalization("discord.command.generic.yes"),
			Value:             1,
		},
		{
			Name:              lang.GetDefault("discord.command.generic.no"),
			NameLocalizations: *util.TranslateLocalization("discord.command.generic.no"),
			Value:             0,
		},
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.visible"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.visible"),
		Description:              lang.GetDefault(tp + "option.set.option.visible.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.visible.description"),
		Choices:                  choices,
		MinValue:                 &minValueZero,
		MaxValue:                 1,
	}
}

func subCommandRemove() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.remove"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.remove"),
		Description:              lang.GetDefault(tp + "option.remove.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.remove.description"),
	}
}

func subCommandList() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		commandOptionListMonth(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.list"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.list"),
		Description:              lang.GetDefault(tp + "option.list.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.list.description"),
		Options:                  options,
	}
}

func commandOptionListMonth() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.list.option.month"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.list.option.month"),
		Description:              lang.GetDefault(tp + "option.list.option.month.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.list.option.month.description"),
		Required:                 true,
		Choices:                  monthChoices("", 0, false),
		MinValue:                 &minValueOne,
		MaxValue:                 12,
	}
}

func subCommandAnnounce() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.announce"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.announce"),
		Description:              lang.GetDefault(tp + "option.announce.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.announce.description"),
	}
}
