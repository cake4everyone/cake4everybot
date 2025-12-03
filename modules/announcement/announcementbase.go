package announcement

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

var log = logger.New("Announcement")

type announcementBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}
