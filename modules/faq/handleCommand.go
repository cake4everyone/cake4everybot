package faq

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

func (cmd Chat) handleCommand(question string) {
	faqs, err := cmd.getAllFAQs()
	if err != nil {
		log.Printf("ERROR: getting all FAQs: %v", err)
		cmd.ReplyError()
		return
	}

	if len(faqs) == 0 {
		cmd.ReplyHiddenSimpleEmbed(0xFAB1FD, lang.GetDefault(tp+"msg.no_questions"))
		return
	}

	e := &discordgo.MessageEmbed{
		Color: 0xFAB1FD,
		Title: question,
	}
	util.SetEmbedFooter(cmd.Session, tp+"display", e)

	if question != "" {
		var ok bool
		e.Description, ok = faqs[question]
		if !ok {
			cmd.ReplyHiddenSimpleEmbedf(0xFAB1FD, lang.GetDefault(tp+"msg.question_not_found"), question)
			return
		}
		cmd.ReplyEmbed(e)
		return
	}

	e.Title = "FAQs"
	var components []discordgo.MessageComponent
	var i int
	for question := range faqs {
		i++
		util.AddEmbedField(e, fmt.Sprintf("%d", i), question, true)
		components = append(components, util.CreateButtonComponent(fmt.Sprintf("faq.show_question.%s", question), fmt.Sprint(i), discordgo.PrimaryButton, nil))
	}
	components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: components}}

	cmd.ReplyComponentsEmbed(components, e)
}
