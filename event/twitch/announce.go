package twitch

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/util"
	webTwitch "github.com/cake4everyone/cake4everybot/webserver/twitch"
	"github.com/kesuaheli/twitchgo"
)

// HandleChannelUpdate is the event handler for the "channel.update" event from twitch.
func HandleChannelUpdate(s *discordgo.Session, t *twitchgo.Session, e *webTwitch.ChannelUpdateEvent) {
	HandleStreamAnnouncementChange(s, t, e.BroadcasterUserID, e.Title, false)
}

// HandleStreamOnline is the event handler for the "stream.online" event from twitch.
func HandleStreamOnline(s *discordgo.Session, t *twitchgo.Session, e *webTwitch.StreamOnlineEvent) {
	HandleStreamAnnouncementChange(s, t, e.BroadcasterUserID, "", true)
}

// HandleStreamOffline is the event handler for the "stream.offline" event from twitch.
func HandleStreamOffline(s *discordgo.Session, t *twitchgo.Session, e *webTwitch.StreamOfflineEvent) {
	HandleStreamAnnouncementChange(s, t, e.BroadcasterUserID, "", false)
}

// HandleStreamAnnouncementChange is a general event handler for twitch events, that should update
// the discord announcement embed.
func HandleStreamAnnouncementChange(s *discordgo.Session, t *twitchgo.Session, platformID, title string, sendNotification bool) {
	announcements, err := database.GetAnnouncement(database.AnnouncementPlatformTwitch, platformID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	for _, announcement := range announcements {
		err = updateAnnouncementMessage(s, t, announcement, title, sendNotification)
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

func updateAnnouncementMessage(s *discordgo.Session, t *twitchgo.Session, announcement *database.Announcement, title string, sendNotification bool) error {
	msg, err := getAnnouncementMessage(s, announcement)
	if err != nil {
		return fmt.Errorf("get announcement in channel '%s': %v", announcement, err)
	}

	var (
		embed  *discordgo.MessageEmbed
		user   *twitchgo.User
		stream *twitchgo.Stream
	)

	if msg == nil || len(msg.Embeds) == 0 {
		embed = &discordgo.MessageEmbed{}
		util.SetEmbedFooter(s, "module.twitch.embed_footer", embed)
	} else {
		embed = msg.Embeds[0]
	}
	users, err := t.GetUsersByID(announcement.PlatformID)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("get users: found no user with ID '%s'", announcement.PlatformID)
	}
	user = users[0]
	streams, err := t.GetStreamsByID(announcement.PlatformID)
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

	if sendNotification {
		notificationContent := announcement.Notification
		if announcement.Notification == "" {
			notificationContent = user.DisplayName
		} else if strings.Contains(announcement.Notification, "%s") {
			notificationContent = fmt.Sprintf(announcement.Notification, user.DisplayName)
		}
		if announcement.RoleID != "" {
			notificationContent += (&discordgo.Role{ID: announcement.RoleID}).Mention()
		}
		msgNotification, err := s.ChannelMessageSend(announcement.ChannelID, notificationContent)
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

func setDefaultEmbed(embed *discordgo.MessageEmbed, user *twitchgo.User) {
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

func setOnlineEmbed(embed *discordgo.MessageEmbed, title string, user *twitchgo.User, stream *twitchgo.Stream) {
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

func setOfflineEmbed(embed *discordgo.MessageEmbed, user *twitchgo.User) {
	setDefaultEmbed(embed, user)

	embed.Title = ""
	embed.URL = ""
	embed.Color = 2829358

	embed.Fields = nil
	embed.Image.URL = user.OfflineImageURL
}
