package random

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// The set subcommand. Used when executing the slash-command "/random teams".
type subcommandTeams struct {
	randomBase
	*Chat
	data *discordgo.ApplicationCommandInteractionDataOption
}

func (rb randomBase) subcommandTeams() subcommandTeams {
	return subcommandTeams{randomBase: rb}
}

// Constructor for subcommandTeams, the struct for the slash-command "/random teams".
func (cmd *Chat) subcommandTeams() subcommandTeams {
	var subcommand *discordgo.ApplicationCommandInteractionDataOption
	if cmd.Interaction != nil {
		subcommand = cmd.Interaction.ApplicationCommandData().Options[0]
	}
	return subcommandTeams{
		randomBase: cmd.randomBase,
		Chat:       cmd,
		data:       subcommand,
	}
}

func (cmd subcommandTeams) appCmd() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.teams"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.teams"),
		Description:              lang.GetDefault(tp + "option.teams.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.teams.description"),
	}
}

func (cmd subcommandTeams) handle() {
}
