package twitch

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kesuaheli/twitchgo"
)

var (
	dcSession              *discordgo.Session
	tSession               *twitchgo.Session
	dcChannelUpdateHandler func(*discordgo.Session, *twitchgo.Session, *ChannelUpdateEvent)
	dcStreamOnlineHandler  func(*discordgo.Session, *twitchgo.Session, *StreamOnlineEvent)
	dcStreamOfflineHandler func(*discordgo.Session, *twitchgo.Session, *StreamOfflineEvent)
	subscriptions          = make(map[string]bool)
)

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
	if !subscriptions[channelID] {
		subscriptions[channelID] = true
		log.Printf("subscribed '%s' for announcements", channelID)
	}
}

// UnsubscribeChannel removes the given channel id from the
// subscription list and no longer sends events.
func UnsubscribeChannel(channelID string) {
	if subscriptions[channelID] {
		delete(subscriptions, channelID)
		log.Printf("unsubscribed '%s' from announcements", channelID)
	}
}

// RefreshSubscriptions sends subscription requests for all registered channels
// to subscribe stream events.
func RefreshSubscriptions() {
	var (
		subscribedChannelUpdate = make(map[string]bool)
		subscribedStreamOnline  = make(map[string]bool)
		subscribedStreamOffline = make(map[string]bool)
	)

	getSubscriptions, err := tSession.GetSubscriptions(false)
	if err != nil {
		log.Printf("Error on getting subscribed channels for %s: %v", err, twitchgo.EventChannelUpdate)
		return
	}
	for _, s := range getSubscriptions {
		if s.Status == twitchgo.SubscriptionStatusWebhookCallbackVerificationFailed {
			err = tSession.DeleteSubscription(s.ID)
			if err != nil {
				log.Printf("Error on deleting failed subscription '%s' for %s (%s): %v", s.Type, s.Condition["broadcaster_user_id"], s.ID, err)
			} else {
				log.Printf("Deleted failed subscription '%s' for %s", s.Type, s.Condition["broadcaster_user_id"])
			}
			continue
		}

		switch s.Type {
		case twitchgo.EventChannelUpdate:
			subscribedChannelUpdate[s.Condition["broadcaster_user_id"]] = true
		case twitchgo.EventStreamOnline:
			subscribedStreamOnline[s.Condition["broadcaster_user_id"]] = true
		case twitchgo.EventStreamOffline:
			subscribedStreamOffline[s.Condition["broadcaster_user_id"]] = true
		}
	}

	for broadcasterID := range subscriptions {
		if !subscribedChannelUpdate[broadcasterID] {
			subscribe(broadcasterID, twitchgo.EventChannelUpdate)
		}
		if !subscribedStreamOnline[broadcasterID] {
			subscribe(broadcasterID, twitchgo.EventStreamOnline)
		}
		if !subscribedStreamOffline[broadcasterID] {
			subscribe(broadcasterID, twitchgo.EventStreamOffline)
		}
	}
}

func subscribe(broadcasterID string, event twitchgo.SubscriptionType) {
	err := tSession.SubscribeToEvent(broadcasterID, CALLBACKURL, event)
	if err != nil {
		log.Printf("Error on subscribing '%s' for %s: %v", event, broadcasterID, err)
	} else {
		log.Printf("Requested subscription to '%s' for %s", event, broadcasterID)
	}
}
