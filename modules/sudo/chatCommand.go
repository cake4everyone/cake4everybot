package sudo

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/util"
)

// The User command of the sudo package.
type UserDisconnect struct {
	sudoBase

	data discordgo.ApplicationCommandInteractionData
	ID   string
}

// AppCmd (ApplicationCommand) returns the definition of the user command
func (cmd UserDisconnect) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:              discordgo.UserApplicationCommand,
		Name:              lang.GetDefault(tp + "user.disconnect.base"),
		NameLocalizations: util.TranslateLocalization(tp + "user.disconnect.base"),
	}
}

// Handle handles the functionality of a command
func (cmd UserDisconnect) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	cmd.member = i.Member
	cmd.user = i.User
	if i.Member != nil {
		cmd.user = i.Member.User
	} else if i.User != nil {
		cmd.member = &discordgo.Member{User: i.User}
	}

	cmd.data = cmd.Interaction.ApplicationCommandData()

	if !cmd.isDisconnectable(cmd.data.TargetID) {
		cmd.ReplyModal("noop", lang.GetDefault(tp+"msg.disconnect.modal.no_permission.title"),
			discordgo.TextDisplay{
				Content: lang.GetDefault(tp + "msg.disconnect.modal.no_permission.content"),
			},
		)
		return
	}

	log.Printf("guild id: %s", cmd.Interaction.GuildID)
	err := s.GuildMemberMove(cmd.Interaction.GuildID, cmd.data.TargetID, nil)
	if restErr, ok := err.(*discordgo.RESTError); ok {
		log.Printf("is rest error: %+v", restErr)
		if restErr.Response.StatusCode == http.StatusNotFound {
			cmd.ReplyModal("noop", lang.GetDefault(tp+"msg.disconnect.modal.no_voice.title"),
				discordgo.TextDisplay{
					Content: lang.GetDefault(tp + "msg.disconnect.modal.no_voice.content"),
				},
			)
			return
		}
	}
	if err != nil {
		log.Printf("ERROR: could not disconnect user %s in guild %s: %+v", cmd.data.TargetID, cmd.Interaction.GuildID, err)
		cmd.ReplyError()
		return
	}
	cmd.QuitSilently()
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *UserDisconnect) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd UserDisconnect) GetID() string {
	return cmd.ID
}

func (cmd UserDisconnect) isDisconnectable(targetID string) bool {
	log.Printf("//TODO: implement sudo permission check for %s disconnecting %s", cmd.user.ID, targetID)
	isSelf := cmd.user.ID == targetID
	return isSelf || database.HasSudoCommandPermission(cmd.user.ID, targetID, database.PermissionDisconnectUsers)
}
