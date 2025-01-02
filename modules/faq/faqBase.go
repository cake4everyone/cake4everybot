package faq

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

var log = logger.New("FAQ")

type faqBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}
