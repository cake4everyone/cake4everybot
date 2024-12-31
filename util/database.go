package util

import (
	"database/sql"
	"fmt"

	"github.com/cake4everyone/cake4everybot/database"
)

// ForAllPlatformIDs calls the given function for all platform IDs of the given platform.
// If no platform IDs exist for the given platform, the function will not be called.
func ForAllPlatformIDs(platform database.Platform, f func(platformID string)) error {
	channels, err := database.GetAllAnnouncementIDs(platform)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("get all %s announcement ids: %w", platform, err)
	}

	for _, channelID := range channels {
		f(channelID)
	}
	return nil
}
