package twitch

import (
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/tools/streamelements"
)

var (
	log *logger.Logger = logger.New("Event/Twitch")
	se  *streamelements.Streamelements
)
