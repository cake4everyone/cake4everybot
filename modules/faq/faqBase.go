package faq

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

var log = logger.New("FAQ")

var (
	// lastFAQs is a cache for all the FAQs
	lastFAQs = make(map[string]map[string]string)
	// lastFAQTime is the time when the lastFAQs map was updated
	lastFAQTime time.Time
)

type faqBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

func (faq faqBase) getAllFAQs() (map[string]string, error) {
	if time.Since(lastFAQTime) < 2*time.Minute {
		return lastFAQs[faq.Interaction.GuildID], nil
	}
	delete(lastFAQs, faq.Interaction.GuildID)
	lastFAQs[faq.Interaction.GuildID] = make(map[string]string)

	row, err := database.Query("SELECT question, answer FROM faq WHERE guild_id=?", faq.Interaction.GuildID)
	if errors.Is(err, sql.ErrNoRows) {
		return lastFAQs[faq.Interaction.GuildID], nil
	} else if err != nil {
		return lastFAQs[faq.Interaction.GuildID], fmt.Errorf("getting all FAQs from guild %s: %w", faq.Interaction.GuildID, err)
	}
	defer row.Close()

	for row.Next() {
		var question, answer string
		if err := row.Scan(&question, &answer); err != nil {
			return lastFAQs[faq.Interaction.GuildID], fmt.Errorf("scanning row: %w", err)
		}
		lastFAQs[faq.Interaction.GuildID][question] = answer
	}

	lastFAQTime = time.Now()
	return lastFAQs[faq.Interaction.GuildID], nil
}
