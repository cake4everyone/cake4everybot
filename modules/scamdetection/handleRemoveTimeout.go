package scamdetection

import "github.com/cake4everyone/cake4everybot/util"

func (c Component) handleRemoveTimeout(ids []string) {
	authorID := util.ShiftL(ids)

	log.Printf("remove timeout from user %s", authorID)
}
