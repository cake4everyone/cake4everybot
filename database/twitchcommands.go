package database

import "database/sql"

// TwitchCommand is a representation of a twitch command response.
type TwitchCommand struct {
	// The type of the response
	ResponseType TwitchCommandResponse
	// The content of the response
	Response string
}

// TwitchCommandResponse is the type of a twitch command response.
type TwitchCommandResponse byte

// TwitchCommandResponse type constants
const (
	TwitchCommandResponseChat TwitchCommandResponse = iota
	TwitchCommandResponseMention

	TwitchCommandResponseFunc = 255
)

// GetTwitchCommand gets a twitch command from the database.
func GetTwitchCommand(channelName, command string) (cmd TwitchCommand, err error) {
	err = QueryRow("SELECT response_type,response FROM twitchcommands WHERE channel_name=? AND command=?", channelName, command).Scan(&cmd.ResponseType, &cmd.Response)
	if err == sql.ErrNoRows {
		return TwitchCommand{}, sql.ErrNoRows
	} else if err != nil {
		log.Printf("ERROR: failed to get twitch command from database: %v", err)
		return TwitchCommand{}, err
	}
	return cmd, nil
}
