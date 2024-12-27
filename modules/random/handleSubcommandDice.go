package random

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// The set subcommand. Used when executing the slash-command "/random dice".
type subcommandDice struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption
}

// Constructor for subcommandDice, the struct for the slash-command "/random dice".
func (cmd Chat) subcommandDice() subcommandDice {
	var subcommand *discordgo.ApplicationCommandInteractionDataOption
	if cmd.Interaction != nil {
		subcommand = cmd.Interaction.ApplicationCommandData().Options[0]
	}
	return subcommandDice{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandDice) appCmd() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		cmd.optionRange(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.dice"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.dice"),
		Description:              lang.GetDefault(tp + "option.dice.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.dice.description"),
		Options:                  options,
	}
}

func (cmd subcommandDice) optionRange() *discordgo.ApplicationCommandOption {
	minValueTwo := float64(2)
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.dice.option.range"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.dice.option.range"),
		Description:              lang.GetDefault(tp + "option.dice.option.range.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.dice.option.range.description"),
		Required:                 false,
		MinValue:                 &minValueTwo,
	}
}

func (cmd subcommandDice) handle() {
}
