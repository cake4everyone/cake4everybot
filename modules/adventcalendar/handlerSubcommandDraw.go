package adventcalendar

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/util"
)

func (cmd Chat) handleSubcommandDraw() {
	now := time.Now()
	year := now.Year() - 2000
	if now.Month() < time.December || (now.Month() == time.December && now.Day() < 25) {
		// draw from last year's entries instead
		year--
	}
	prefix := fmt.Sprintf("xmas%d", year)
	winner, totalTickets := database.DrawGiveawayWinner(database.GetAllGiveawayEntries(prefix, database.AnnouncementPlatformDiscord, cmd.Interaction.GuildID))
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
