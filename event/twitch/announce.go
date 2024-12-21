package twitch

import (
	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/twitch"
	"cake4everybot/util"
	webTwitch "cake4everybot/webserver/twitch"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HandleChannelUpdate is the event handler for the "channel.update" event from twitch.
func HandleChannelUpdate(s *discordgo.Session, e *webTwitch.ChannelUpdateEvent) {
	HandleStreamAnnouncementChange(s, e.BroadcasterUserID, e.Title, "")
}

// HandleStreamOnline is the event handler for the "stream.online" event from twitch.
func HandleStreamOnline(s *discordgo.Session, e *webTwitch.StreamOnlineEvent) {
	HandleStreamAnnouncementChange(s, e.BroadcasterUserID, "", lang.GetDefault("module.twitch.msg.nofification"))
}

// HandleStreamOffline is the event handler for the "stream.offline" event from twitch.
func HandleStreamOffline(s *discordgo.Session, e *webTwitch.StreamOfflineEvent) {
	HandleStreamAnnouncementChange(s, e.BroadcasterUserID, "", "")
}

// HandleStreamAnnouncementChange is a general event handler for twitch events, that should update
// the discord announcement embed.
func HandleStreamAnnouncementChange(s *discordgo.Session, platformID, title, notification string) {
	announcements, err := database.GetAnnouncement(database.AnnouncementPlatformTwitch, platformID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	for _, announcement := range announcements {
		err = updateAnnouncementMessage(s, announcement, title, notification)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func getAnnouncementMessage(s *discordgo.Session, announcement *database.Announcement) (msg *discordgo.Message, err error) {
	if announcement.MessageID == "" {
		return nil, nil
	}

	msg, err = s.ChannelMessage(announcement.ChannelID, announcement.MessageID)
	if restErr, ok := err.(*discordgo.RESTError); ok {
		// if the lastMessageID returns a 404, i.e. it was deleted, create a new one
		if restErr.Response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	return msg, err
}

func newAnnouncementMessage(s *discordgo.Session, announcement *database.Announcement, embed *discordgo.MessageEmbed) (msg *discordgo.Message, err error) {
	msg, err = s.ChannelMessageSendEmbed(announcement.ChannelID, embed)
	if err != nil {
		return
	}
	return msg, announcement.UpdateAnnouncementMessage(msg.ID)
}

func updateAnnouncementMessage(s *discordgo.Session, announcement *database.Announcement, title, notification string) error {
	msg, err := getAnnouncementMessage(s, announcement)
	if err != nil {
		return fmt.Errorf("get announcement in channel '%s': %v", announcement, err)
	}

	var (
		embed  *discordgo.MessageEmbed
		user   *twitch.User
		stream *twitch.Stream
	)

	if msg == nil || len(msg.Embeds) == 0 {
		embed = &discordgo.MessageEmbed{}
		util.SetEmbedFooter(s, "module.twitch.embed_footer", embed)
	} else {
		embed = msg.Embeds[0]
	}
	users, err := twitch.GetUsersByID(announcement.PlatformID)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("get users: found no user with ID '%s'", announcement.PlatformID)
	}
	user = users[0]
	streams, err := twitch.GetStreamsByID(announcement.PlatformID)
	if err != nil {
		return err
	}
	if len(streams) == 0 {
		stream = nil
	} else {
		stream = streams[0]
	}

	if stream != nil {
		setOnlineEmbed(embed, title, user, stream)
	} else {
		setOfflineEmbed(embed, user)
	}

	if notification != "" {
		if announcement.RoleID != "" {
			notification += fmt.Sprintf("\n<@&%s>", announcement.RoleID)
		}
		msgNotification, err := s.ChannelMessageSend(announcement.ChannelID, fmt.Sprintf(notification, user.DisplayName))
		if err != nil {
			return fmt.Errorf("send notification: %v", err)
		}
		go s.ChannelMessageDelete(announcement.ChannelID, msgNotification.ID)
	}

	if msg == nil {
		_, err = newAnnouncementMessage(s, announcement, embed)
	} else {
		m := discordgo.NewMessageEdit(announcement.ChannelID, msg.ID).SetEmbed(embed)
		m.Flags = msg.Flags & (math.MaxInt - discordgo.MessageFlagsSuppressEmbeds)
		m.Flags |= 1 << 12 // setting SUPPRESS_NOTIFICATIONS bit just to prevent Flags to be '0' and thus get removed by the json omitempty
		_, err = s.ChannelMessageEditComplex(m)
	}
	if err != nil {
		return fmt.Errorf("update announcement in channel '%s': %v", announcement, err)
	}
	return nil
}

func setDefaultEmbed(embed *discordgo.MessageEmbed, user *twitch.User) {
	embed.Author = &discordgo.MessageEmbedAuthor{
		URL:     fmt.Sprintf("https://twitch.tv/%s/about", user.Login),
		Name:    user.DisplayName,
		IconURL: user.ProfileImageURL,
	}
	if embed.Image == nil {
		embed.Image = &discordgo.MessageEmbedImage{}
	}
	embed.Image.Width = 1920
	embed.Image.Height = 1080
}

func setOnlineEmbed(embed *discordgo.MessageEmbed, title string, user *twitch.User, stream *twitch.Stream) {
	setDefaultEmbed(embed, user)

	if title == "" {
		embed.Title = stream.Title
	} else {
		embed.Title = title
	}
	embed.URL = fmt.Sprintf("https://twitch.tv/%s", user.Login)
	embed.Color = 9520895

	if len(embed.Fields) == 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{})
	}
	embed.Fields[0].Name = lang.GetDefault("module.twitch.embed_category")
	embed.Fields[0].Value = fmt.Sprintf("[%s](https://twitch.tv/directory/category/%s)", stream.GameName, stream.GameID)
	embed.Image.URL = strings.ReplaceAll(stream.ThumbnailURL, "{width}x{height}", "1920x1080")
}

func setOfflineEmbed(embed *discordgo.MessageEmbed, user *twitch.User) {
	setDefaultEmbed(embed, user)

	embed.Title = ""
	embed.URL = ""
	embed.Color = 2829358

	embed.Fields = nil
	embed.Image.URL = user.OfflineImageURL
}
