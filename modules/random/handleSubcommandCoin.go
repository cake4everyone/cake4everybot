package random

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// The set subcommand. Used when executing the slash-command "/random coin".
type subcommandCoin struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption
}

// Constructor for subcommandCoin, the struct for the slash-command "/random coin".
func (cmd Chat) subcommandCoin() subcommandCoin {
	var subcommand *discordgo.ApplicationCommandInteractionDataOption
	if cmd.Interaction != nil {
		subcommand = cmd.Interaction.ApplicationCommandData().Options[0]
	}
	return subcommandCoin{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandCoin) appCmd() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.coin"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.coin"),
		Description:              lang.GetDefault(tp + "option.coin.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.coin.description"),
	}
}

func (cmd subcommandCoin) handle() {
}
