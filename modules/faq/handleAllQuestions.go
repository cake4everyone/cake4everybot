package faq

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/util"
)

func (c Component) handleAllQuestions() {
	faqs, err := c.getAllFAQs()
	if err != nil {
		log.Printf("ERROR: getting all FAQs: %v", err)
		c.ReplyError()
		return
	}

	e := &discordgo.MessageEmbed{
		Color: 0xFAB1FD,
		Title: "FAQs",
	}
	util.SetEmbedFooter(c.Session, tp+"display", e)

	var components []discordgo.MessageComponent
	var i int
	for question := range faqs {
		i++
		util.AddEmbedField(e, fmt.Sprintf("%d", i), question, true)
		components = append(components, util.CreateButtonComponent(fmt.Sprintf("faq.show_question.%s", question), fmt.Sprint(i), discordgo.PrimaryButton, nil))
	}
	components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: components}}

	c.ReplyComponentsEmbedUpdate(components, e)
}
