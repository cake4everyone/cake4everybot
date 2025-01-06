package faq

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
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
		answer = strings.ReplaceAll(answer, "\\n", "\n")
		lastFAQs[faq.Interaction.GuildID][question] = answer
	}

	lastFAQTime = time.Now()
	return lastFAQs[faq.Interaction.GuildID], nil
}

func (faq faqBase) getFAQMessage(question string) (msg *discordgo.InteractionResponseData) {
	msg, faqs, final := faq.getBaseMessage()
	if final {
		return
	}
	msg.Embeds[0].Title = question

	answer, ok := faqs[question]
	if !ok {
		e := util.SimpleEmbedf(0, lang.GetDefault(tp+"msg.question_not_found"), question)[0]
		msg.Embeds[0].Title = e.Title
		msg.Embeds[0].Description = e.Description
		msg.Flags = discordgo.MessageFlagsEphemeral
		return
	}
	msg.Embeds[0].Description = answer

	var components []discordgo.MessageComponent
	components = append(components, util.CreateButtonComponent(
		"faq.all_questions",
		lang.GetDefault(tp+"msg.button.all_questions"),
		discordgo.SecondaryButton,
		util.GetConfigComponentEmoji("faq.all_questions"),
	))
	components = append(components, util.CloseButtonComponent())
	msg.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: components}}

	return
}

func (faq faqBase) getAllFAQsMessage() (msg *discordgo.InteractionResponseData) {
	msg, faqs, final := faq.getBaseMessage()
	if final {
		return
	}

	var components []discordgo.MessageComponent
	var i int
	for question := range faqs {
		i++
		util.AddEmbedField(
			msg.Embeds[0],
			fmt.Sprintf("%d", i),
			question,
			true,
		)
		components = append(components, util.CreateButtonComponent(
			fmt.Sprintf("faq.show_question.%s", question),
			fmt.Sprint(i),
			discordgo.PrimaryButton,
			nil,
		))
	}
	components = append(components, util.CloseButtonComponent())
	msg.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: components}}

	return
}

// getBaseMessage is a helper function to get the base message for the FAQs. It
// returns the an prefilled message with one embed.
//
// It also returns and checks the FAQs. If final is true, then the message is
// ready to send. This may be caused by an error or similar.
func (faq faqBase) getBaseMessage() (msg *discordgo.InteractionResponseData, faqs map[string]string, final bool) {
	msg = &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{{
			Color:  0xFAB1FD,
			Title:  "FAQs",
			Footer: util.EmbedFooter(faq.Session, tp+"display"),
		}},
	}

	faqs, err := faq.getAllFAQs()
	if err != nil {
		log.Printf("ERROR: getting all FAQs: %v", err)
		msg.Embeds[0].Description = lang.GetDefault(tp + "msg.error")
		msg.Flags = discordgo.MessageFlagsEphemeral
		final = true
		return
	}

	if len(faqs) == 0 {
		e := util.SimpleEmbed(0, lang.GetDefault(tp+"msg.no_questions"))[0]
		msg.Embeds[0].Title = e.Title
		msg.Embeds[0].Description = e.Description
		msg.Flags = discordgo.MessageFlagsEphemeral
		final = true
		return
	}

	return
}
