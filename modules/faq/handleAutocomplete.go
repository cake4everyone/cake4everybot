package faq

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (cmd Chat) handleAutocomplete(input string) {
	faqs, err := cmd.getAllFAQs()
	if err != nil {
		log.Printf("ERROR: getting all FAQs: %v", err)
		cmd.ReplyError()
		return
	}

	options := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(faqs))
	for question := range faqs {
		if input == "" || strings.Contains(strings.ToLower(question), strings.ToLower(input)) {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  question,
				Value: question,
			})
		}
	}
	cmd.ReplyAutocomplete(options)
}
