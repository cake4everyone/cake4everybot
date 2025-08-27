package database

import (
	"database/sql"
	"fmt"
	"regexp"
)

// TwitchTimer is a representation of a twitch chat timer.
type TwitchTimer struct {
	ChannelName  string
	Title        *regexp.Regexp
	Minutes      byte
	ResponseType TwitchCommandResponse
	Response     string
}

// GetAllTwitchTimers gets all the timers from the database.
func GetAllTwitchTimers() (twitchTimers []TwitchTimer, err error) {
	resp, err := Query("SELECT channel_name,title,minutes,response_type,response FROM twitchtimers")
	if err == sql.ErrNoRows {
		return twitchTimers, nil
	} else if err != nil {
		return nil, fmt.Errorf("get timers from database: %w", err)
	}

	for resp.Next() {
		var tt TwitchTimer
		var title sql.NullString
		err = resp.Scan(&tt.ChannelName, &title, &tt.Minutes, &tt.ResponseType, &tt.Response)
		if err != nil {
			return nil, fmt.Errorf("scan timers row: %w", err)
		}
		if title.Valid {
			tt.Title, err = regexp.Compile(title.String)
			if err != nil {
				return nil, fmt.Errorf("compile title regexp: %w", err)
			}
		}
		twitchTimers = append(twitchTimers, tt)
	}
	return twitchTimers, nil
}
