package twitch

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	logger "log"
	"net/http"
	"slices"
	"time"

	"github.com/spf13/viper"
)

// rawEvent represents the http body comming with a call to the /twitch_pubsub enpoints
type rawEvent struct {
	// Challenge cointains the string to return when receiving a webhook callback verification.
	// Otherwise it is an empty string
	Challenge string `json:"challenge"`

	// Subscription contains the informations this event is about.
	Subscription Subscription `json:"subscription"`

	// Event is the actual event.
	//
	// It is not set in a webhook callback verification.
	Event interface{} `json:"event"`
}

var (
	log          = logger.New(logger.Writer(), "[WebTwitch] ", logger.LstdFlags|logger.Lmsgprefix)
	lastMessages = make([]string, 10)
)

// HandlePost is the HTTP/POST handler for the Twitch PubSub endpoint.
//
// It is called to handle a webhook comming from twitch. This could be a hub challenge verification
// or a event notification.
func HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// before anything, check the hash
	if !verifyTwitchMessage(r.Header, body) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var rEvent rawEvent
	err = json.Unmarshal(body, &rEvent)
	if err != nil {
		log.Printf("Failed to unmarshal body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messageType := r.Header.Get("Twitch-Eventsub-Message-Type")
	switch messageType {
	case "webhook_callback_verification":
		handleVerification(w, r, rEvent)
		return
	case "notification":
		data, _ := json.Marshal(rEvent.Event)
		if err := handleNotification(data, rEvent.Subscription.Type); err != nil {
			log.Printf("Error on notification event: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "revocation":
		log.Println("requested to revoke")
	default:
		log.Printf("Unknown message type '%s'", messageType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func verifyTwitchMessage(header http.Header, body []byte) bool {
	// read more at https://dev.twitch.tv/docs/eventsub/handling-webhook-events/#verifying-the-event-message
	msgID := header.Get("Twitch-Eventsub-Message-Id")
	msgTime := header.Get("Twitch-Eventsub-Message-Timestamp")
	data := append([]byte(msgID+msgTime), body...)

	h := hmac.New(sha256.New, []byte(viper.GetString("twitch.webhookSecret")))
	h.Write(data)
	hmacData := h.Sum(nil)

	hmacHex := make([]byte, 128)
	n := hex.Encode(hmacHex, hmacData)
	hmacHex = append([]byte("sha256="), hmacHex[:n]...)

	signature := []byte(header.Get("Twitch-Eventsub-Message-Signature"))
	if !hmac.Equal(hmacHex, signature) {
		log.Printf("calculated sha does not match, got '%s' want '%s'", hmacHex, signature)
		return false
	}

	t, err := time.Parse(time.RFC3339, msgTime)
	if err != nil {
		log.Printf("Error parsing timestamp '%s': %v", msgTime, err)
		return false
	}

	if time.Until(t) < -10*time.Minute {
		log.Printf("message is older than 10 minutes: %v (%s)", time.Until(t), t)
		return false
	}

	if slices.Contains(lastMessages, msgID) {
		log.Printf("message id already verified: id '%s' found in last messages: %v", msgID, lastMessages)
		return false
	}
	lastMessages = append(lastMessages[1:], msgID)
	return true
}

func handleVerification(w http.ResponseWriter, _ *http.Request, rEvent rawEvent) {
	if rEvent.Challenge == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	broadcaster := rEvent.Subscription.Condition["broadcaster_user_id"]
	if subscribtions[broadcaster] {
		log.Printf("Declined verification for broadcaster '%s'!", broadcaster)
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("{\"conflict\":\"that broadcaster is not allowed\"}"))
		return
	}

	log.Printf("Accepted '%s v%s' for channel %s", rEvent.Subscription.Type, rEvent.Subscription.Version, broadcaster)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(rEvent.Challenge))
}

func handleNotification(data []byte, notificationType string) error {
	switch notificationType {
	case "channel.update":
		var e ChannelUpdateEvent
		err := json.Unmarshal(data, &e)
		if err != nil {
			return fmt.Errorf("parse channel update: %v", err)
		}
		go dcChannelUpdateHandler(dcSession, &e)
	case "stream.online":
		var e StreamOnlineEvent
		err := json.Unmarshal(data, &e)
		if err != nil {
			return fmt.Errorf("parse online: %v", err)
		}
		go dcStreamOnlineHandler(dcSession, &e)
	case "stream.offline":
		var e StreamOfflineEvent
		err := json.Unmarshal(data, &e)
		if err != nil {
			return fmt.Errorf("parse offline: %v", err)
		}
		go dcStreamOfflineHandler(dcSession, &e)
	default:
		log.Printf("Unhandled notification type '%s': %s", notificationType, data)
	}
	return nil
}
