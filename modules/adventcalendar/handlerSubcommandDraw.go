package adventcalendar

import (
	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (cmd Chat) handleSubcommandDraw() {
	winner, totalTickets := database.DrawGiveawayWinner(database.GetAllGiveawayEntries("xmas", database.AnnouncementPlatformDiscord, cmd.Interaction.GuildID))
	if totalTickets == 0 {
		cmd.ReplyHidden(lang.GetDefault(tp + "msg.no_entries.draw"))
		return
	}

	member, err := cmd.Session.GuildMember(cmd.Interaction.GuildID, winner.UserID)
	if err != nil {
		log.Printf("WARN: Could not get winner as member '%s' from guild '%s': %v", cmd.Interaction.GuildID, winner.UserID, err)
		log.Print("Trying to get user instead...")

		user, err := cmd.Session.User(winner.UserID)
		if err != nil {
			log.Printf("ERROR: Could not get winner user '%s': %v", winner.UserID, err)
			cmd.ReplyError()
			return
		}
		member = &discordgo.Member{User: user}
	}

	name := member.Nick
	if name == "" {
		name = member.User.Username
	}

	e := &discordgo.MessageEmbed{
		Title: lang.GetDefault(tp + "msg.winner.title"),
		Description: fmt.Sprintf(
			lang.GetDefault(tp+"msg.winner.details"),
			member.Mention(),
			winner.Weight,
			float64(100*winner.Weight)/float64(totalTickets),
		),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.AvatarURL(""),
		},
		Color: 0x00A000,
		Fields: []*discordgo.MessageEmbedField{{
			Value: fmt.Sprintf(lang.GetDefault(tp+"msg.winner.congratulation"), name),
		}},
	}
	util.SetEmbedFooter(cmd.Session, "module.adventcalendar.embed_footer", e)
	cmd.ReplyEmbed(e)
}
