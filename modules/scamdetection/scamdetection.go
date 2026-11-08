// Package scamdetection is a module to detect scams and timouts them
package scamdetection

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

var log = logger.New("ScamDetection")

const (
	RoleID     = "" // FIXME: read from db
	LogChannel = "" // FIXME: read from db
)

type scamDetectionBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

func possibleScammer(s *discordgo.Session, author *discordgo.Member, message *discordgo.Message) {
	timoutUntil := time.Now().Add(15 * time.Second)
	err := s.GuildMemberTimeout(message.GuildID, author.User.ID, &timoutUntil)
	if err != nil {
		log.Printf("ERROR timout scammer: %s", err)
	}

	files := make([]*discordgo.File, 0, len(message.Attachments))
	for _, attachment := range message.Attachments {
		resp, err := http.Get(attachment.URL)
		if err != nil {
			log.Printf("ERROR download attachment: %s", err)
			continue
		}
		attachmentBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("ERROR read attachment: %s", err)
			continue
		}
		buf := bytes.NewBuffer(attachmentBytes)

		files = append(files, &discordgo.File{
			Name:        attachment.Filename,
			ContentType: attachment.ContentType,
			Reader:      buf,
		})
	}

	err = s.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		log.Printf("ERROR delete scammer message: %s", err)
	}

	accountCreation, _ := discordgo.SnowflakeTimestamp(author.User.ID)
	const ComponentIDFormat = "scamdetection.%s.%s"
	data := &discordgo.MessageSend{
		Content: fmt.Sprintf("Possible scammer '%s'", author.Mention()),
		Embeds: []*discordgo.MessageEmbed{{
			Title:       "Possible Scammer Detection",
			Description: fmt.Sprintf("%s ID: `%s`", author.Mention(), author.User.ID),
			Color:       0xFF8800,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Account created at",
					Value: fmt.Sprintf("<t:%d:R>", accountCreation.UTC().Unix()),
				},
				{
					Name:  "Joined at",
					Value: fmt.Sprintf("<t:%d:R>", author.JoinedAt.UTC().Unix()),
				},
				{
					Name:  "Timed out until",
					Value: fmt.Sprintf("<t:%d:R>", timoutUntil.UTC().Unix()),
				},
				{
					Name:  "Number of attachments",
					Value: fmt.Sprint(len(message.Attachments)),
				},
			},
		}},
		Files: files,
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Ban User",
					Style:    discordgo.DangerButton,
					CustomID: fmt.Sprintf(ComponentIDFormat, "ban", author.User.ID),
				},
				discordgo.Button{
					Label:    "Remove Timeout",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf(ComponentIDFormat, "remove_timeout", author.User.ID),
				},
			},
		}},
	}
	_, err = s.ChannelMessageSendComplex(LogChannel, data)
	if err != nil {
		log.Printf("ERROR mod log message: %s", err)
	}
}

// HandleGuildJoin is the handler for new members joined.
func HandleGuildJoin(s *discordgo.Session, member *discordgo.GuildMemberAdd) {
	log.Printf("New member: %s", member.User.Username)
	err := s.GuildMemberRoleAdd(member.GuildID, member.User.ID, RoleID)
	if err != nil {
		log.Printf("ERROR adding role to new member: %s", err)
	}
}
