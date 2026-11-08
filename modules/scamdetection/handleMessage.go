package scamdetection

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/util"
)

// HandleMessage is the handler for messages sent in a channel.
func HandleMessage(s *discordgo.Session, message *discordgo.MessageCreate) {
	if !util.ContainsString(message.Member.Roles, RoleID) {
		return
	}
	if len(message.Attachments) == 0 {
		err := s.GuildMemberRoleRemove(message.GuildID, message.Author.ID, RoleID)
		if err != nil {
			log.Printf("ERROR remove role from member: %s", err)
		}
		return
	}

	// New member sends attachments in their first message
	author := message.Member
	author.User = message.Author
	possibleScammer(s, author, message.Message)
}
