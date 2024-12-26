package twitch

import (
	"cake4everybot/logger"
	"cake4everybot/tools/streamelements"
)

var (
	log *logger.Logger = logger.New("Event/Twitch")
	se  *streamelements.Streamelements
)
