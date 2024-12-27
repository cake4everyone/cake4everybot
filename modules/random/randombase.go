package random

import (
	"cake4everybot/logger"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

var log = logger.New("Random")

type randomBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}
