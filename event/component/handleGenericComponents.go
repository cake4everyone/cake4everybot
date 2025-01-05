package component

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/util"
)

// GenericComponents is the Component handler for generic components
type GenericComponents struct {
	util.InteractionUtil
	member         *discordgo.Member
	user           *discordgo.User
	originalAuthor *discordgo.User
	data           discordgo.MessageComponentInteractionData
}

// Handle handles the functionality of a component interaction
func (gc GenericComponents) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gc.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	gc.member = i.Member
	gc.user = i.User
	if i.Member != nil {
		gc.user = i.Member.User
	} else if i.User != nil {
		gc.member = &discordgo.Member{User: i.User}
	}
	gc.data = i.MessageComponentData()

	if gc.Interaction.Message.Type == discordgo.MessageTypeChatInputCommand {
		gc.originalAuthor = gc.Interaction.Message.Interaction.User
	} else {
		gc.originalAuthor = gc.Interaction.Message.Author
	}

	ids := strings.Split(gc.data.CustomID, ".")
	// pop the first level identifier
	util.ShiftL(ids)

	switch util.ShiftL(ids) {
	case "delete":
		if gc.RequireOriginalAuthor() {
			gc.handleDelete()
		}
		return
	default:
		log.Printf("Unknown component interaction ID: %s", gc.data.CustomID)
		gc.ReplyError()
	}
}

// ID returns the component ID to identify the module
func (gc GenericComponents) ID() string {
	return "generic"
}

func (gc GenericComponents) handleDelete() {
	err := gc.Session.ChannelMessageDelete(gc.Interaction.ChannelID, gc.Interaction.Message.ID)
	if err != nil {
		log.Printf("ERROR: could not delete message %s/%s: %+v", gc.Interaction.ChannelID, gc.Interaction.Message.ID, err)
		gc.ReplyError()
	}
}
