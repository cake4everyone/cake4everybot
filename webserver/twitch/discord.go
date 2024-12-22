package twitch

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kesuaheli/twitchgo"
)

var dcSession *discordgo.Session
var tSession *twitchgo.Session
var dcChannelUpdateHandler func(*discordgo.Session, *twitchgo.Session, *ChannelUpdateEvent)
var dcStreamOnlineHandler func(*discordgo.Session, *twitchgo.Session, *StreamOnlineEvent)
var dcStreamOfflineHandler func(*discordgo.Session, *twitchgo.Session, *StreamOfflineEvent)
var subscribtions = make(map[string]bool)

// SetDiscordSession sets the discordgo.Session to use for calling
// event handlers.
func SetDiscordSession(s *discordgo.Session) {
	dcSession = s
}

// SetTwitchSession sets the twitchgo.Session to use for calling
// event handlers.
func SetTwitchSession(t *twitchgo.Session) {
	tSession = t
}

// SetDiscordChannelUpdateHandler sets the function to use when calling event
// handlers.
func SetDiscordChannelUpdateHandler(f func(*discordgo.Session, *twitchgo.Session, *ChannelUpdateEvent)) {
	dcChannelUpdateHandler = f
}

// SetDiscordStreamOnlineHandler sets the function to use when calling event
// handlers.
func SetDiscordStreamOnlineHandler(f func(*discordgo.Session, *twitchgo.Session, *StreamOnlineEvent)) {
	dcStreamOnlineHandler = f
}

// SetDiscordStreamOfflineHandler sets the function to use when calling event
// handlers.
func SetDiscordStreamOfflineHandler(f func(*discordgo.Session, *twitchgo.Session, *StreamOfflineEvent)) {
	dcStreamOfflineHandler = f
}

// SubscribeChannel subscribe to the event listener for new videos of
// the given channel id.
func SubscribeChannel(channelID string) {
	if !subscribtions[channelID] {
		subscribtions[channelID] = true
		log.Printf("subscribed '%s' for announcements", channelID)
	}
}

// UnsubscribeChannel removes the given channel id from the
// subscription list and no longer sends events.
func UnsubscribeChannel(channelID string) {
	if subscribtions[channelID] {
		delete(subscribtions, channelID)
		log.Printf("unsubscribed '%s' from announcements", channelID)
	}
}
