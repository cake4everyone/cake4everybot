package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleModalSetAddress(ids []string) {
	c.Interaction.GuildID = util.ShiftL(ids)
	players, err := c.getPlayers()
	if err != nil {
		log.Printf("ERROR: could not get players: %+v", err)
		c.ReplyError()
		return
	}
	if len(players) == 0 {
		log.Printf("ERROR: no players in guild %s", c.Interaction.GuildID)
		c.ReplyError()
		return
	}
	player, ok := players[c.Interaction.User.ID]
	if !ok {
		log.Printf("ERROR: could not find player %s in guild %s: %+v", c.Interaction.User.ID, c.Interaction.GuildID, c.Interaction.User.ID)
		c.ReplyError()
		return
	}

	addressFiled := c.modal.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
	if addressFiled.Value == player.Address {
		c.ReplyHidden(lang.GetDefault(tp + "msg.invite.set_address.not_changed"))
		return
	}

	player.Address = addressFiled.Value
	player.PendingNudge = false
	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	_, err = c.Session.ChannelMessageEditEmbed(c.Interaction.ChannelID, player.MessageID, player.InviteEmbed(c.Session))
	if err != nil {
		log.Printf("ERROR: could not update bot message for %s '%s/%s': %+v", player.DisplayName(), c.Interaction.ChannelID, player.MessageID, err)
		c.ReplyError()
		return
	}

	santaPlayer := c.getSantaForPlayer(player.User.ID)
	santaChannel, err := c.Session.UserChannelCreate(santaPlayer.User.ID)
	if err != nil {
		log.Printf("ERROR: could not get user channel for %s: %+v", santaPlayer.DisplayName(), err)
		c.ReplyError()
		return
	}
	_, err = c.Session.ChannelMessageEditEmbed(santaChannel.ID, santaPlayer.MessageID, santaPlayer.InviteEmbed(c.Session))
	if err != nil {
		log.Printf("ERROR: could not update bot message for %s '%s/%s': %+v", santaPlayer.DisplayName(), santaChannel.ID, santaPlayer.MessageID, err)
		c.ReplyError()
		return
	}
	_, err = c.Session.ChannelMessageSendComplex(santaChannel.ID, &discordgo.MessageSend{
		Content:   lang.GetDefault(tp + "msg.invite.set_address.match_updated"),
		Reference: &discordgo.MessageReference{MessageID: santaPlayer.MessageID},
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.delete",
				lang.GetDefault(tp+"msg.invite.button.delete"),
				discordgo.DangerButton,
				util.GetConfigComponentEmoji("secretsanta.invite.delete"),
			),
		}}},
	})
	if err != nil {
		log.Printf("ERROR: could not send address update message for %s '%s/%s': %+v", santaPlayer.DisplayName(), santaChannel.ID, santaPlayer.MessageID, err)
		c.ReplyError()
		return
	}

	e := &discordgo.MessageEmbed{
		Color: 0x00FF00,
		Fields: []*discordgo.MessageEmbedField{{
			Name:  lang.GetDefault(tp + "msg.invite.set_address.changed"),
			Value: fmt.Sprintf("```\n%s\n```", player.Address),
		}},
	}

	util.SetEmbedFooter(c.Session, tp+"display", e)
	c.ReplyHiddenEmbed(e)
}

func (c Component) handleModalAddPackageTracking(ids []string) {
	c.Interaction.GuildID = util.ShiftL(ids)
	players, err := c.getPlayers()
	if err != nil {
		log.Printf("ERROR: could not get players: %+v", err)
		c.ReplyError()
		return
	}
	if len(players) == 0 {
		log.Printf("ERROR: no players in guild %s", c.Interaction.GuildID)
		c.ReplyError()
		return
	}
	player, ok := players[c.Interaction.User.ID]
	if !ok {
		log.Printf("ERROR: could not find player %s in guild %s: %+v", c.Interaction.User.ID, c.Interaction.GuildID, c.Interaction.User.ID)
		c.ReplyError()
		return
	}

	packageTrackingField := c.modal.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
	if packageTrackingField.Value == player.PackageTracking {
		c.ReplyHidden(lang.GetDefault(tp + "msg.invite.add_package_tracking.not_changed"))
		return
	}

	player.PackageTracking = packageTrackingField.Value
	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	var matchChannel *discordgo.Channel
	if matchChannel, _, ok = c.updateInviteMessage(player.Match); !ok {
		c.ReplyError()
		return
	}
	if player.PackageTracking != "" {
		_, err = c.Session.ChannelMessageSendComplex(matchChannel.ID, &discordgo.MessageSend{
			Content:   lang.GetDefault(tp + "msg.invite.add_package_tracking.santa_updated"),
			Reference: &discordgo.MessageReference{MessageID: player.Match.MessageID},
			Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				util.CreateButtonComponent(
					"secretsanta.invite.delete",
					lang.GetDefault(tp+"msg.invite.button.delete"),
					discordgo.DangerButton,
					util.GetConfigComponentEmoji("secretsanta.invite.delete"),
				),
			}}},
		})
		if err != nil {
			log.Printf("ERROR: could not send package tracking update message for %s '%s/%s': %+v", player.Match.DisplayName(), matchChannel.ID, player.Match.MessageID, err)
			c.ReplyError()
			return
		}
	}

	c.ReplyHiddenSimpleEmbed(0x690042, lang.GetDefault(tp+"msg.invite.add_package_tracking.success"))
}
