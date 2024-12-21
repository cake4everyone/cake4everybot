package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleInvite(ids []string) {
	switch util.ShiftL(ids) {
	case "show_match":
		c.handleInviteShowMatch(ids)
		return
	case "set_address":
		c.handleInviteSetAddress(ids)
		return
	case "nudge_match":
		c.handleInviteNudgeMatch(ids)
		return
	case "confirm_nudge":
		c.handleInviteConfirmNudge(ids)
		return
	case "send_package":
		c.handleInviteSendPackage(ids)
		return
	case "add_package_tracking":
		c.handleAddPackageTracking(ids)
		return
	case "show_package_tracking":
		c.handleShowPackageTracking(ids)
		return
	case "confirm_send_package":
		c.handleInviteConfirmSendPackage(ids)
		return
	case "received_package":
		c.handleInviteReceivedPackage(ids)
		return
	case "confirm_received_package":
		c.handleInviteConfirmReceivedPackage(ids)
		return
	case "delete":
		err := c.Session.ChannelMessageDelete(c.Interaction.ChannelID, c.Interaction.Message.ID)
		if err != nil {
			log.Printf("ERROR: could not delete message %s/%s: %+v", c.Interaction.ChannelID, c.Interaction.Message.ID, err)
			c.ReplyError()
		}
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}
}

func (c Component) handleInviteShowMatch(ids []string) {
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

	e := util.AuthoredEmbed(c.Session, player.Match.Member, tp+"display")
	e.Title = fmt.Sprintf(lang.GetDefault(tp+"msg.invite.show_match.title"), player.Match.Member.DisplayName())
	e.Description = lang.GetDefault(tp + "msg.invite.show_match.description")
	e.Color = 0x690042
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:  lang.GetDefault(tp + "msg.invite.show_match.address"),
		Value: fmt.Sprintf("```\n%s\n```\n%s", player.Match.Address, lang.GetDefault(tp+"msg.invite.show_match.nudge_description")),
	})
	if player.Match.Address == "" {
		e.Fields[0].Value = lang.GetDefault(tp + "msg.invite.show_match.address_not_set")
	}

	util.SetEmbedFooter(c.Session, tp+"display", e)

	var components []discordgo.MessageComponent
	if player.SendPackage == 0 {
		components = append(components, util.CreateButtonComponent(
			fmt.Sprintf("secretsanta.invite.nudge_match.%s", c.Interaction.GuildID),
			lang.GetDefault(tp+"msg.invite.button.nudge_match"),
			discordgo.SecondaryButton,
			util.GetConfigComponentEmoji("secretsanta.invite.nudge_match"),
		))
		if player.Match.Address != "" {
			components = append(components, util.CreateButtonComponent(
				fmt.Sprintf("secretsanta.invite.send_package.%s", c.Interaction.GuildID),
				lang.GetDefault(tp+"msg.invite.button.send_package"),
				discordgo.SuccessButton,
				util.GetConfigComponentEmoji("secretsanta.invite.send_package"),
			))
		}
	} else {
		components = append(components, util.CreateButtonComponent(
			fmt.Sprintf("secretsanta.invite.add_package_tracking.%s", c.Interaction.GuildID),
			lang.GetDefault(tp+"msg.invite.button.add_package_tracking"),
			discordgo.SecondaryButton,
			util.GetConfigComponentEmoji("secretsanta.invite.add_package_tracking"),
		))
	}
	if len(components) == 0 {
		c.ReplyHiddenEmbed(e)
	} else {
		c.ReplyComponentsHiddenEmbed([]discordgo.MessageComponent{discordgo.ActionsRow{Components: components}}, e)
	}
}

func (c Component) handleInviteSetAddress(ids []string) {
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

	c.ReplyModal("secretsanta.set_address."+c.Interaction.GuildID, lang.GetDefault(tp+"msg.invite.modal.set_address.title"), discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "address",
			Label:       lang.GetDefault(tp + "msg.invite.modal.set_address.label"),
			Style:       discordgo.TextInputParagraph,
			Placeholder: lang.GetDefault(tp + "msg.invite.modal.set_address.placeholder"),
			Value:       player.Address,
			Required:    true,
		},
	}})
}

func (c Component) handleInviteNudgeMatch(ids []string) {
	c.ReplyComponentsHiddenSimpleEmbedUpdate(
		[]discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.confirm_nudge."+strings.Join(ids, "."),
				lang.GetDefault(tp+"msg.invite.button.nudge_match"),
				discordgo.PrimaryButton,
				util.GetConfigComponentEmoji("secretsanta.invite.nudge_match"),
			),
		}}},
		0x690042,
		lang.GetDefault(tp+"msg.invite.nudge_match.confirm"))
}

func (c Component) handleInviteConfirmNudge(ids []string) {
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
	player.Match.PendingNudge = true
	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	matchChannel, err := c.Session.UserChannelCreate(player.Match.User.ID)
	if err != nil {
		log.Printf("ERROR: could not create DM channel with user %s: %+v", player.Match.User.ID, err)
		c.ReplyError()
		return
	}
	_, err = c.Session.ChannelMessageEditEmbed(matchChannel.ID, player.Match.MessageID, player.Match.InviteEmbed(c.Session))
	if err != nil {
		log.Printf("ERROR: could not edit match message embed: %+v", err)
		c.ReplyError()
		return
	}

	data := &discordgo.MessageSend{
		Content:   lang.GetDefault(tp + "msg.invite.nudge_received"),
		Reference: &discordgo.MessageReference{MessageID: player.Match.MessageID},
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.delete",
				lang.GetDefault(tp+"msg.invite.button.delete"),
				discordgo.DangerButton,
				util.GetConfigComponentEmoji("secretsanta.invite.delete"),
			),
		}}},
	}
	_, err = c.Session.ChannelMessageSendComplex(matchChannel.ID, data)
	if err != nil {
		log.Printf("ERROR: could not send nudge message: %+v", err)
		c.ReplyError()
		return
	}

	_, err = c.Session.ChannelMessageEditEmbed(c.Interaction.ChannelID, player.MessageID, player.InviteEmbed(c.Session))
	if err != nil {
		log.Printf("ERROR: could not edit invite message embed: %+v", err)
		c.ReplyError()
		return
	}
	c.ReplyHiddenSimpleEmbedUpdate(0x690042, lang.GetDefault(tp+"msg.invite.nudge_match.success"))
}

func (c Component) handleInviteSendPackage(ids []string) {
	c.ReplyComponentsHiddenSimpleEmbedUpdate(
		[]discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.confirm_send_package."+strings.Join(ids, "."),
				lang.GetDefault(tp+"msg.invite.button.send_package"),
				discordgo.SuccessButton,
				util.GetConfigComponentEmoji("secretsanta.invite.send_package"),
			),
		}}},
		0x690042,
		lang.GetDefault(tp+"msg.invite.send_package.confirm"))
}

func (c Component) handleAddPackageTracking(ids []string) {
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

	c.ReplyModal("secretsanta.add_package_tracking."+c.Interaction.GuildID, lang.GetDefault(tp+"msg.invite.modal.add_package_tracking.title"), discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "package_tracking",
			Label:       lang.GetDefault(tp + "msg.invite.modal.add_package_tracking.label"),
			Style:       discordgo.TextInputParagraph,
			Placeholder: lang.GetDefault(tp + "msg.invite.modal.add_package_tracking.placeholder"),
			Value:       player.PackageTracking,
			Required:    false,
		},
	}})
}

func (c Component) handleShowPackageTracking(ids []string) {
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

	c.ReplyHiddenSimpleEmbedf(0x690042, "## %s\n%s", lang.GetDefault(tp+"msg.invite.package_tracking.title"), c.getSantaForPlayer(player.User.ID).PackageTracking)
}

func (c Component) handleInviteConfirmSendPackage(ids []string) {
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
	player.SendPackage = 1
	player.Match.PendingNudge = false
	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	if _, _, ok = c.updateInviteMessage(player); !ok {
		c.ReplyError()
		return
	}
	var matchChannel *discordgo.Channel
	if matchChannel, _, ok = c.updateInviteMessage(player.Match); !ok {
		c.ReplyError()
		return
	}

	data := &discordgo.MessageSend{
		Content:   lang.GetDefault(tp + "msg.invite.send_package"),
		Reference: &discordgo.MessageReference{MessageID: player.Match.MessageID},
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.delete",
				lang.GetDefault(tp+"msg.invite.button.delete"),
				discordgo.DangerButton,
				util.GetConfigComponentEmoji("secretsanta.invite.delete"),
			),
		}}},
	}
	_, err = c.Session.ChannelMessageSendComplex(matchChannel.ID, data)
	if err != nil {
		log.Printf("ERROR: could not send nudge message: %+v", err)
		c.ReplyError()
		return
	}
	c.ReplyHiddenSimpleEmbedUpdate(0x690042, lang.GetDefault(tp+"msg.invite.send_package.success"))
}

func (c Component) handleInviteReceivedPackage(ids []string) {
	c.ReplyComponentsHiddenSimpleEmbed(
		[]discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.confirm_received_package."+strings.Join(ids, "."),
				lang.GetDefault(tp+"msg.invite.button.received_package"),
				discordgo.SuccessButton,
				util.GetConfigComponentEmoji("secretsanta.invite.received_package"),
			),
		}}},
		0x690042,
		lang.GetDefault(tp+"msg.invite.received_package.confirm"))
}

func (c Component) handleInviteConfirmReceivedPackage(ids []string) {
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
	santaPlayer := c.getSantaForPlayer(player.User.ID)
	santaPlayer.SendPackage = 2
	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	if _, _, ok = c.updateInviteMessage(player); !ok {
		c.ReplyError()
		return
	}
	var santaChannel *discordgo.Channel
	if santaChannel, _, ok = c.updateInviteMessage(santaPlayer); !ok {
		c.ReplyError()
		return
	}

	data := &discordgo.MessageSend{
		Content:   lang.GetDefault(tp + "msg.invite.received_package"),
		Reference: &discordgo.MessageReference{MessageID: santaPlayer.MessageID},
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent(
				"secretsanta.invite.delete",
				lang.GetDefault(tp+"msg.invite.button.delete"),
				discordgo.DangerButton,
				util.GetConfigComponentEmoji("secretsanta.invite.delete"),
			),
		}}},
	}
	_, err = c.Session.ChannelMessageSendComplex(santaChannel.ID, data)
	if err != nil {
		log.Printf("ERROR: could not send nudge message: %+v", err)
		c.ReplyError()
		return
	}
	c.ReplyHiddenSimpleEmbedUpdate(0x690042, lang.GetDefault(tp+"msg.invite.received_package.success"))
}
