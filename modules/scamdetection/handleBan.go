package scamdetection

import "github.com/cake4everyone/cake4everybot/util"

func (c Component) handleBan(ids []string) {
	authorID := util.ShiftL(ids)

	log.Printf("ban user %s", authorID)
}
