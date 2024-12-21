package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleModalSetAddress() {
	addressFiled := c.modal.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
	if addressFiled.Value == c.player.Address {
		c.ReplyHidden(lang.GetDefault(tp + "msg.invite.set_address.not_changed"))
		return
	}

	c.player.Address = addressFiled.Value
	c.player.PendingNudge = false
	err := c.setPlayers()
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	_, err = c.Session.ChannelMessageEditEmbed(c.Interaction.ChannelID, c.player.MessageID, c.player.InviteEmbed(c.Session))
	if err != nil {
		log.Printf("ERROR: could not update bot message for %s '%s/%s': %+v", c.player.DisplayName(), c.Interaction.ChannelID, c.player.MessageID, err)
		c.ReplyError()
		return
	}

	santaPlayer := c.getSantaForPlayer(c.player.User.ID)
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
			Value: fmt.Sprintf("```\n%s\n```", c.player.Address),
		}},
	}

	util.SetEmbedFooter(c.Session, tp+"display", e)
	c.ReplyHiddenEmbed(e)
}

func (c Component) handleModalAddPackageTracking() {
	packageTrackingField := c.modal.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
	if packageTrackingField.Value == c.player.PackageTracking {
		c.ReplyHidden(lang.GetDefault(tp + "msg.invite.add_package_tracking.not_changed"))
		return
	}

	c.player.PackageTracking = packageTrackingField.Value
	err := c.setPlayers()
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	var matchChannel *discordgo.Channel
	var ok bool
	if matchChannel, _, ok = c.updateInviteMessage(c.player.Match); !ok {
		c.ReplyError()
		return
	}
	if c.player.PackageTracking != "" {
		_, err = c.Session.ChannelMessageSendComplex(matchChannel.ID, &discordgo.MessageSend{
			Content:   lang.GetDefault(tp + "msg.invite.add_package_tracking.santa_updated"),
			Reference: &discordgo.MessageReference{MessageID: c.player.Match.MessageID},
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
			log.Printf("ERROR: could not send package tracking update message for %s '%s/%s': %+v", c.player.Match.DisplayName(), matchChannel.ID, c.player.Match.MessageID, err)
			c.ReplyError()
			return
		}
	}

	c.ReplyHiddenSimpleEmbed(0x690042, lang.GetDefault(tp+"msg.invite.add_package_tracking.success"))
}
