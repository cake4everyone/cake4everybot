package faq

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

func (c Component) handleShowQuestion(question string) {
	faqs, err := c.getAllFAQs()
	if err != nil {
		log.Printf("ERROR: getting all FAQs: %v", err)
		c.ReplyError()
		return
	}

	e := &discordgo.MessageEmbed{
		Color: 0xFAB1FD,
		Title: question,
	}
	util.SetEmbedFooter(c.Session, tp+"display", e)

	var ok bool
	e.Description, ok = faqs[question]
	if !ok {
		log.Printf("ERROR: could not find question: '%s'", question)
		log.Printf("Available questions: %v", faqs)
		c.ReplyError()
		return
	}

	var components []discordgo.MessageComponent
	components = append(components, util.CreateButtonComponent(
		"faq.all_questions",
		lang.GetDefault(tp+"msg.button.all_questions"),
		discordgo.SecondaryButton,
		util.GetConfigComponentEmoji("faq.all_questions"),
	))
	components = append(components, util.CloseButtonComponent())

	components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: components}}

	c.ReplyComponentsEmbedUpdate(components, e)
}
