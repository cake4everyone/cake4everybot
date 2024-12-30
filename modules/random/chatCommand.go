package random

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => random
	tp = "discord.command.random."
)

// The Chat (slash) command of the random package. Has a few sub commands and options to use all
// features through a single chat command.
type Chat struct {
	randomBase
	ID string
}

type subcommand interface {
	appCmd() *discordgo.ApplicationCommandOption
	handle()
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd Chat) AppCmd() *discordgo.ApplicationCommand {
	options := []*discordgo.ApplicationCommandOption{
		cmd.subcommandCoin().appCmd(),
		cmd.subcommandDice().appCmd(),
		cmd.subcommandTeams().appCmd(),
	}

	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "base"),
		NameLocalizations:        util.TranslateLocalization(tp + "base"),
		Description:              lang.GetDefault(tp + "base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "base.description"),
		Options:                  options,
	}
}

// Handle handles the functionality of a command
func (cmd Chat) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	cmd.member = i.Member
	cmd.user = i.User
	if i.Member != nil {
		cmd.user = i.Member.User
	} else if i.User != nil {
		cmd.member = &discordgo.Member{User: i.User}
	}

	subcommandName := i.ApplicationCommandData().Options[0].Name
	var sub subcommand
	switch subcommandName {
	case lang.GetDefault(tp + "option.dice"):
		sub = cmd.subcommandDice()
	case lang.GetDefault(tp + "option.coin"):
		sub = cmd.subcommandCoin()
	case lang.GetDefault(tp + "option.teams"):
		sub = cmd.subcommandTeams()
	default:
		return
	}

	sub.handle()
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *Chat) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd Chat) GetID() string {
	return cmd.ID
}
