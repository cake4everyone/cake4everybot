package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// The Chat (slash) command of the secret santa package.
type Chat struct {
	secretSantaBase
	ID string
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (Chat) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "cmd.base"),
		NameLocalizations:        util.TranslateLocalization(tp + "cmd.base"),
		Description:              lang.GetDefault(tp + "cmd.base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "cmd.base.description"),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:                     discordgo.ApplicationCommandOptionSubCommand,
				Name:                     lang.GetDefault(tp + "cmd.option.show"),
				NameLocalizations:        *util.TranslateLocalization(tp + "cmd.option.show"),
				Description:              lang.GetDefault(tp + "cmd.option.show.description"),
				DescriptionLocalizations: *util.TranslateLocalization(tp + "cmd.option.show.description"),
			},
			{
				Type:                     discordgo.ApplicationCommandOptionSubCommand,
				Name:                     lang.GetDefault(tp + "cmd.option.update"),
				NameLocalizations:        *util.TranslateLocalization(tp + "cmd.option.update"),
				Description:              lang.GetDefault(tp + "cmd.option.update.description"),
				DescriptionLocalizations: *util.TranslateLocalization(tp + "cmd.option.update.description"),
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

	switch i.ApplicationCommandData().Options[0].Name {
	case lang.GetDefault(tp + "cmd.option.show"):
		cmd.handleSubcommandShow()
		return
	case lang.GetDefault(tp + "cmd.option.update"):
		cmd.handleSubcommandUpdate()
		return
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
