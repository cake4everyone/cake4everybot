package secretsanta

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

func (cmd MsgCmd) handler() {
	joinEmoji, err := util.GetConfigEmoji(cmd.Session, "secretsanta.enter")
	if err != nil {
		log.Printf("ERROR: could not get emoji: %+v", err)
		cmd.ReplyError()
		return
	}
	joinEmojiID := joinEmoji.ID
	if joinEmojiID == "" {
		joinEmojiID = joinEmoji.Name
	}

	msg := cmd.data.Resolved.Messages[cmd.data.TargetID]
	if len(msg.Reactions) == 0 {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.no_reactions"), joinEmojiID)
		return
	}
	var hasReaction bool
	for _, r := range msg.Reactions {
		if !util.CompareEmoji(r.Emoji, joinEmoji) {
			continue
		}
		hasReaction = true
		break
	}

	if !hasReaction {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.no_reactions"), joinEmojiID)
		return
	}

	users, err := cmd.Session.MessageReactions(msg.ChannelID, msg.ID, joinEmojiID, 100, "", "")
	if err != nil {
		log.Printf("Error on get users: %v\n", err)
		cmd.ReplyError()
		return
	}

	if len(users) < 2 {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.not_enough_reactions"), 2)
		return
	}

	e := &discordgo.MessageEmbed{
		Title: lang.GetDefault(tp + "title"),
		Color: 0x690042,
	}

	var (
		names   string
		players = map[string]*player{}
	)
	for _, u := range users {
		member, err := cmd.Session.GuildMember(cmd.Interaction.GuildID, u.ID)
		if member == nil {
			log.Printf("WARN: Could not get member '%s' from guild '%s': %v", u.ID, cmd.Interaction.GuildID, err)
			continue
		}
		players[u.ID] = &player{Member: member}
		names += fmt.Sprintf("%s\n", member.Mention())
	}
	if len(players) < 2 {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.setup.not_enough_reactions"), 2)
		return
	}
	util.AddEmbedField(e, lang.GetDefault(tp+"msg.setup.users"), names, false)

	allPlayers[cmd.Interaction.GuildID] = players
	err = cmd.setPlayers()
	if err != nil {
		log.Printf("Error on set players: %v\n", err)
		cmd.ReplyError()
		return
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			util.CreateButtonComponent("secretsanta.setup.invite", "Invite", discordgo.SuccessButton, nil),
		}},
	}

	util.SetEmbedFooter(cmd.Session, tp+"display", e)
	cmd.ReplyComponentsHiddenEmbed(components, e)
}
