package twitch

import (
	"bytes"
	"cake4everybot/webserver/oauth"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Subscription represents a single subscription to an event.
type Subscription struct {
	// ID is the unique identifier for this subscription.
	ID string `json:"id"`

	// Status is the status of this subscription e.g. it is set to "enabled" when successfully
	// verified and active.
	Status string `json:"status"`

	// Type is the acual event that triggers.
	Type string `json:"type"`

	// Version is the version number of the event defined in the Type field.
	Version string `json:"version"`

	// Condition contains a list of key-value-pairs of conditions. The 'key' gives a variable to check
	// and 'value' the value to match. For example "broadcaster_user_id":"12345" requires the event
	// to be triggered at the channel of user "12345".
	Condition map[string]string `json:"condition"`

	// Transport gives information about how this subscription is (or will be) delivered.
	Transport SubscriptionTransport `json:"transport"`

	// CreatedAt is the timestamp of creation of this subscription.
	CreatedAt time.Time `json:"created_at"`

	// The amount points this subscription costs. The cost is added to a global count. Each
	// application has a fixed amount of available points to use.
	Cost int `json:"cost"`
}

// SubscriptionTransport gives information about how a subscription is (or will be) delivered.
type SubscriptionTransport struct {
	// Method is either set to "webhook" or "websocket".
	Method string `json:"method"`

	// WebhookCallbackURI gives the complete URI of the webhook.
	//
	// Only when Method == "webhook"
	WebhookCallbackURI string `json:"callback,omitempty"`

	// WebhookSecret is the secret given with the creation of the subscription to veryfiy its
	// correctness.
	//
	// Only when Method == "webhook"
	WebhookSecret string `json:"secret,omitempty"`

	// WebSocketSessionID is the ID the welcome message returns, when connecting to the twitch
	// websocket. More information needed.
	//
	// Only when Method == "websocket"
	WebSocketSessionID string `json:"session_id,omitempty"`
}

// ChannelUpdateEvent represents a "channel.update" event from twitch.
type ChannelUpdateEvent struct {
	BroadcasterUserID           string   `json:"broadcaster_user_id"`           // The broadcaster’s user ID.
	BroadcasterUserLogin        string   `json:"broadcaster_user_login"`        // The broadcaster’s user login.
	BroadcasterUserName         string   `json:"broadcaster_user_name"`         // The broadcaster’s user display name.
	Title                       string   `json:"title"`                         // The channel’s stream title.
	Language                    string   `json:"language"`                      // The channel’s broadcast language.
	CategoryID                  string   `json:"category_id"`                   // The channel’s category ID.
	CategoryName                string   `json:"category_name"`                 // The category name.
	ContentClassificationLabels []string `json:"content_classification_labels"` // Array of content classification label IDs currently applied on the Channel. To retrieve a list of all possible IDs, use the Get Content Classification Labels API endpoint.
}

// StreamOnlineEvent represents a "stream.online" event from twitch.
type StreamOnlineEvent struct {
	ID                   string    `json:"id"`                     // The id of the stream.
	BroadcasterUserID    string    `json:"broadcaster_user_id"`    // The broadcaster’s user id.
	BroadcasterUserLogin string    `json:"broadcaster_user_login"` // The broadcaster’s user login.
	BroadcasterUserName  string    `json:"broadcaster_user_name"`  // The broadcaster’s user display name.
	Type                 string    `json:"type"`                   // The stream type. Valid values are: live, playlist, watch_party, premiere, rerun.
	StartedAt            time.Time `json:"started_at"`             // The timestamp at which the stream went online at.
}

// StreamOfflineEvent represents a "stream.offline" event from twitch.
type StreamOfflineEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`    // The broadcaster’s user id.
	BroadcasterUserLogin string `json:"broadcaster_user_login"` // The broadcaster’s user login.
	BroadcasterUserName  string `json:"broadcaster_user_name"`  // The broadcaster’s user display name.
}

// RegisterEventSubscription registers the callback URL for a twitch API event at twitch.
func RegisterEventSubscription() {
	reqURL := "https://api.twitch.tv/helix/eventsub/subscriptions"
	clientID := viper.GetString("twitch.clientID")
	clientSecret := viper.GetString("twitch.clientSecret")
	webhookSecret := viper.GetString("twitch.webhookSecret")

	for id := range subscribtions {
		log.Printf("Requesting subscription refresh for id '%s'...", id)

		var (
			body []byte
			err  error
		)
		if body, err = json.Marshal(Subscription{
			Type:    "channel.update",
			Version: "1",
			Condition: map[string]string{
				"broadcaster_user_id": id,
			},
			Transport: SubscriptionTransport{
				Method:             "webhook",
				WebhookCallbackURI: "https://webhook.cake4everyone.de/api/twitch_pubsub",
				WebhookSecret:      webhookSecret,
			},
		}); err != nil {
			log.Printf("Failed to marshal request body")
			return
		}

		req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
		if err != nil {
			log.Printf("Failed to create refresh subscription: %v", err)
			return
		}

		// App Token
		appToken := oauth.New(
			"https://id.twitch.tv/oauth2/token",
			clientID,
			clientSecret,
			"",
		)
		t, err := appToken.GenerateToken()
		if err != nil {
			log.Printf("Failed to generate token: %v", err)
			return
		}
		log.Printf("Token: %s", t)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+t)
		req.Header.Set("Client-Id", clientID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to do refresh subscription request: %v", err)
			return
		}

		resp.Request = nil
		buf, err := json.MarshalIndent(resp, "", "	")
		if err != nil {
			log.Printf("Failed to read refresh subscription response: %v", err)
			return
		}
		var name string = "response.json"
		err = os.WriteFile(name, buf, 0644)
		if err != nil {
			log.Printf("Failed to save response to file: %v", err)
			return
		}
		log.Printf("Saved reponse to file: ./%s", name)

		buf, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return
		}
		log.Printf("Body:\n%s", string(buf))

		log.Printf("Successfully refreshed subscription for channel '%s'", id)
	}
}
