package announcement

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/util"
)

// The Component of the announcement package.
type Component struct {
	announcementBase
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
	//lint:ignore SA4005 currently not used but will be when implementing the component // TODO: remove when implementing the component
	c.data = i.MessageComponentData()

	ids := strings.Split(c.data.CustomID, ".")
	// pop the first level identifier
	util.ShiftL(ids)

	switch util.ShiftL(ids) {
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}

}

// ID returns the custom ID of the modal to identify the module
func (c Component) ID() string {
	return "announcement"
}
