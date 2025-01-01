package random

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

var log = logger.New("Random")

type randomBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User

	originalAuthor *discordgo.User
}
