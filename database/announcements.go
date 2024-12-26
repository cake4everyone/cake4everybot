package database

// Announcement is a representation of a Discord announcement channel.
//
// It can be obtained by GetAnnouncement for a given channel on a platform.
type Announcement struct {
	GuildID      string
	ChannelID    string
	MessageID    string
	RoleID       string
	Platform     Platform
	PlatformID   string
	Notification string
}

// Platform is the type of platform a Announcement can be made
type Platform uint16

// Platform types for announcements
const (
	AnnouncementPlatformDiscord Platform = iota
	AnnouncementPlatformTwitch
	AnnouncementPlatformYoutube
)

// String implements the fmt.Stringer interface
func (p Platform) String() string {
	switch p {
	case AnnouncementPlatformDiscord:
		return "Discord"
	case AnnouncementPlatformTwitch:
		return "Twitch"
	case AnnouncementPlatformYoutube:
		return "YouTube"
	default:
		return ""
	}
}

// GoString implements the fmt.GoStringer interface
func (p Platform) GoString() string {
	switch p {
	case AnnouncementPlatformDiscord:
		return "discord"
	case AnnouncementPlatformTwitch:
		return "twitch"
	case AnnouncementPlatformYoutube:
		return "youtube"
	default:
		return ""
	}
}

// GetAllAnnouncementIDs returns all platform IDs for a given platform.
//
// If no result matches the given platform the returned error will be sql.ErrNoRows.
// Other errors may exist.
func GetAllAnnouncementIDs(platform Platform) (platformIDs []string, err error) {
	rows, err := Query("SELECT DISTINCT platform_id FROM announcements WHERE platform=?", platform)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	platformIDs = make([]string, 0)
	for rows.Next() {
		var platformID string
		if err := rows.Scan(&platformID); err != nil {
			return nil, err
		}
		platformIDs = append(platformIDs, platformID)
	}
	return platformIDs, nil
}

// GetAnnouncement reads all Discord announcement channels from the database for a given channel ID
// on a platform.
// A platform could be "twitch" or "youtube".
//
// If no result matches the given platform and channel ID the returned error will be sql.ErrNoRows.
// Other errors may exist.
func GetAnnouncement(platform Platform, platformID string) ([]*Announcement, error) {
	rows, err := Query("SELECT guild_id,channel_id,message_id,role_id,notification FROM announcements WHERE platform=? AND platform_id=?", platform, platformID)
	if err != nil {
		return []*Announcement{}, err
	}
	defer rows.Close()
	announcements := make([]*Announcement, 0)
	for rows.Next() {
		var guildID, channelID, messageID, roleID, notification string
		if err := rows.Scan(&guildID, &channelID, &messageID, &roleID); err != nil {
			return []*Announcement{}, err
		}
		announcements = append(announcements, &Announcement{guildID, channelID, messageID, roleID, platform, platformID, notification})
	}
	return announcements, err
}

// UpdateAnnouncementMessage updates the message id of a with newID.
func (a *Announcement) UpdateAnnouncementMessage(newID string) error {
	_, err := Exec("UPDATE announcements SET message_id=? WHERE guild_id=? AND channel_id=? AND message_id=? AND role_id=? AND platform=? AND platform_id=? AND notification=?",
		newID,
		a.GuildID, a.ChannelID, a.MessageID, a.RoleID, a.Platform, a.PlatformID, a.Notification,
	)
	a.MessageID = newID
	return err
}
