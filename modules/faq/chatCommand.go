package faq

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// The Chat (slash) command of the faq package.
type Chat struct {
	faqBase
	ID string
}

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => faq
	tp = "discord.command.faq."
)

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd Chat) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "base"),
		NameLocalizations:        util.TranslateLocalization(tp + "base"),
		Description:              lang.GetDefault(tp + "base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "base.description"),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         lang.GetDefault(tp + "option.question"),
				Description:  lang.GetDefault(tp + "option.question.description"),
				Required:     false,
				Autocomplete: true,
			},
		},
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

	data := i.ApplicationCommandData()
	var question string
	for _, option := range data.Options {
		switch option.Name {
		case lang.GetDefault(tp + "option.question"):
			question = option.StringValue()
		}
	}

	if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
		cmd.handleAutocomplete(question)
	} else if question == "" {
		cmd.ReplyComplex(cmd.getAllFAQsMessage())
		return
	} else {
		cmd.ReplyComplex(cmd.getFAQMessage(question))
	}
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *Chat) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd Chat) GetID() string {
	return cmd.ID
}
