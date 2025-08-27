package timer

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
	"github.com/kesuaheli/twitchgo"
)

var (
	log *logger.Logger = logger.New("Event/Twitch/Timer")
)

// RegisterTimer registers the twitch chat timers.
// They are loaded and parsed from the database and checked if the corresponding
// channel is in the channels list.
func RegisterTimer(t *twitchgo.Session, channels []string) (err error) {
	timers, err := database.GetAllTwitchTimers()
	if err != nil {
		return fmt.Errorf("get timers from database: %w", err)
	}
	log.Printf("Got %d timers from database\n", len(timers))
	for _, timer := range timers {
		if !util.ContainsString(channels, timer.ChannelName) {
			continue
		}
		go runTimer(t, timer)
	}

	return nil
}

func runTimer(t *twitchgo.Session, timer database.TwitchTimer) {
	const randomMultiplier = 0.1
	for {
		randomRange := int(float32(timer.Minutes*60) * randomMultiplier)
		randomOffset := time.Duration(rand.IntN(randomRange*2)-randomRange) * time.Second
		time.Sleep(time.Duration(timer.Minutes)*time.Minute + randomOffset)

		streams, err := t.GetStreamsByName(timer.ChannelName)
		if err != nil {
			log.Printf("ERROR: Could not get stream for channel '%s': %+v\n", timer.ChannelName, err)
			continue
		}
		if len(streams) == 0 || (timer.Title != nil && !timer.Title.MatchString(streams[0].Title)) {
			continue
		}

		switch timer.ResponseType {
		case database.TwitchCommandResponseChat:
			t.SendMessage(timer.ChannelName, timer.Response)
		case database.TwitchCommandResponseFunc:
			log.Printf("WARN: TwitchCommandResponseFunc (%d) is not implemented yet in twitch chat timers. Skipping timer for channel '%s'\n", timer.ResponseType, timer.ChannelName)
			return
		default:
			log.Printf("ERROR: Unknown response type '%d' in twitch timer for channel '%s'\n", timer.ResponseType, timer.ChannelName)
			return
		}
	}
}
