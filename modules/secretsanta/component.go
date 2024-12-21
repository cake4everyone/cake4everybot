package secretsanta

import (
	"cake4everybot/util"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// The Component of the secret santa package.
type Component struct {
	secretSantaBase
	data  discordgo.MessageComponentInteractionData
	modal discordgo.ModalSubmitInteractionData
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

	ids := strings.Split(c.data.CustomID, ".")
	// pop the first level identifier
	util.ShiftL(ids)

	switch util.ShiftL(ids) {
	case "setup":
		c.handleSetup(ids)
		return
	case "invite":
		c.handleInvite(ids)
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}

}

// HandleModal handles the functionality of a modal.
func (c Component) HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	c.member = i.Member
	c.user = i.User
	if i.Member != nil {
		c.user = i.Member.User
	} else if i.User != nil {
		c.member = &discordgo.Member{User: i.User}
	}
	c.modal = i.ModalSubmitData()

	ids := strings.Split(c.modal.CustomID, ".")
	// pop the first level identifier
	util.ShiftL(ids)

	switch util.ShiftL(ids) {
	case "set_address":
		c.handleModalSetAddress(ids)
		return
	case "add_package_tracking":
		c.handleModalAddPackageTracking(ids)
		return
	default:
		log.Printf("Unknown modal submit ID: %s", c.modal.CustomID)
	}
}

// ID returns the custom ID of the modal to identify the module
func (c Component) ID() string {
	return "secretsanta"
}
