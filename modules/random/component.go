package random

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/util"
)

// The Component of the random package.
type Component struct {
	randomBase
	data discordgo.MessageComponentInteractionData
}

// Handle handles the functionality of a component.
func (c Component) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	c.member = i.Member
	c.user = i.User
	if i.Member != nil {
		c.user = i.Member.User
	} else if i.User != nil {
		c.member = &discordgo.Member{User: i.User}
	}
	c.data = i.MessageComponentData()

	if c.Interaction.Message.Type == discordgo.MessageTypeChatInputCommand {
		c.originalAuthor = c.Interaction.Message.Interaction.User
	} else {
		c.originalAuthor = c.Interaction.Message.Author
	}

	ids := strings.Split(c.data.CustomID, ".")
	// pop the first level identifier
	util.ShiftL(ids)

	switch util.ShiftL(ids) {
	case "coin":
		c.subcommandCoin().handleComponent(ids)
		return
	case "dice":
		c.subcommandDice().handleComponent(ids)
		return
	case "teams":
		c.subcommandTeams().handleComponent(ids)
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}

}

// ID returns the custom ID of the modal to identify the module
func (c Component) ID() string {
	return "random"
}
