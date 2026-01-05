package sudo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => sudo
	tp = "discord.command.sudo."
)

var log = logger.New("Sudo")

type sudoBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}
