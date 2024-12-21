package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// MsgCmd represents the mesaage command of the secretsanta package. It adds the ability to start a
// new secret santa game.
type MsgCmd struct {
	secretSantaBase

	data discordgo.ApplicationCommandInteractionData
	ID   string
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd *MsgCmd) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:              discordgo.MessageApplicationCommand,
		Name:              lang.GetDefault(tp + "base"),
		NameLocalizations: util.TranslateLocalization(tp + "base"),
	}
}

// Handle handles the functionality of a command
func (cmd *MsgCmd) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	cmd.member = i.Member
	cmd.user = i.User
	if i.Member != nil {
		cmd.user = i.Member.User
	}

	cmd.data = cmd.Interaction.ApplicationCommandData()
	cmd.handler()
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *MsgCmd) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd *MsgCmd) GetID() string {
	return cmd.ID
}
