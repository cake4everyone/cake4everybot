// Copyright 2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package birthday

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// The set subcommand. Used when executing the slash-command "/birthday set".
type subcommandSet struct {
	Chat
	*discordgo.ApplicationCommandInteractionDataOption

	day     *discordgo.ApplicationCommandInteractionDataOption // reqired
	month   *discordgo.ApplicationCommandInteractionDataOption // reqired
	year    *discordgo.ApplicationCommandInteractionDataOption // optional
	visible *discordgo.ApplicationCommandInteractionDataOption // optional
}

// Constructor for subcommandSet, the struct for the slash-command "/birthday set".
func (cmd Chat) subcommandSet() subcommandSet {
	subcommand := cmd.Interaction.ApplicationCommandData().Options[0]
	return subcommandSet{
		Chat:                                    cmd,
		ApplicationCommandInteractionDataOption: subcommand,
	}
}

func (cmd subcommandSet) handler() {
	for _, opt := range cmd.Options {
		switch opt.Name {
		case lang.GetDefault(tp + "option.set.option.day"):
			cmd.day = opt
		case lang.GetDefault(tp + "option.set.option.month"):
			cmd.month = opt
		case lang.GetDefault(tp + "option.set.option.year"):
			cmd.year = opt
		case lang.GetDefault(tp + "option.set.option.visible"):
			cmd.visible = opt
		}
	}

	switch cmd.Interaction.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		cmd.autocompleteHandler()
	case discordgo.InteractionApplicationCommand:
		cmd.interactionHandler()
	}
}

func (cmd subcommandSet) autocompleteHandler() {
	var choices []*discordgo.ApplicationCommandOptionChoice

	if cmd.day != nil && cmd.day.Focused {
		start := cmd.day.Value.(string)

		var m int
		if cmd.month != nil {
			m = int(cmd.month.IntValue())
		}
		if m == 0 {
			m = 1
		}
		leapYear := cmd.year == nil || cmd.year.IntValue()%4 == 0

		choices = dayChoices(start, m, leapYear)
	} else if cmd.month != nil && cmd.month.Focused {
		start := cmd.month.Value.(string)

		var d int
		if cmd.day != nil {
			d = int(cmd.day.IntValue())
		}
		leapYear := cmd.year == nil || cmd.year.IntValue()%4 == 0

		choices = monthChoices(start, d, leapYear)
	} else if cmd.year != nil && cmd.year.Focused {
		start := cmd.year.Value.(string)

		var d, m int
		if cmd.day != nil {
			d = int(cmd.day.IntValue())
		}
		if cmd.month != nil {
			m = int(cmd.month.IntValue())
		}
		if m == 0 {
			m = 1
		}

		choices = yearChoices(start, d, m)
	}

	cmd.ReplyAutocomplete(choices)
}

// executes when running the subcommand
func (cmd subcommandSet) interactionHandler() {
	authorID, err := strconv.ParseUint(cmd.user.ID, 10, 64)
	if err != nil {
		log.Printf("Error on parse author id of birthday command: %v\n", err)
		cmd.ReplyError()
		return
	}

	b := birthdayEntry{
		ID:      authorID,
		Day:     int(cmd.day.IntValue()),
		Month:   int(cmd.month.IntValue()),
		Visible: true,
	}
	if cmd.year != nil {
		b.Year = int(cmd.year.IntValue())
	}
	if cmd.visible != nil {
		b.Visible = cmd.visible.IntValue() == 1
	}

	embed := util.AuthoredEmbed(cmd.Session, cmd.member, tp+"display")

	b.time, err = time.Parse(time.DateOnly, fmt.Sprintf("%04d-%02d-%02d", b.Year, b.Month, b.Day))
	if err != nil {
		log.Printf("WARNING: User (%d) entered an invalid date: %v\n", authorID, err)
		embed.Description = lang.Get(tp+"msg.invalid_date", lang.FallbackLang())
		embed.Color = 0xFF0000
		cmd.ReplyHiddenEmbed(embed)
		return
	}

	hasBDay, err := cmd.hasBirthday(b.ID)
	if err != nil {
		log.Printf("Error on getting birthday data: %v\n", err)
		cmd.ReplyError()
		return
	}

	if hasBDay {
		err = cmd.handleUpdate(b, embed)
		if err != nil {
			log.Printf("Error on update birthday: %v\n", err)
			cmd.ReplyError()
			return
		}
	} else {
		err = cmd.setBirthday(b)
		if err != nil {
			log.Printf("Error on set birthday: %v\n", err)
			cmd.ReplyError()
			return
		}

		var age string
		if b.Year > 0 {
			age = fmt.Sprintf(" (%d)", b.Age()+1)
		}

		embed.Description = lang.Get(tp+"msg.set", lang.FallbackLang())
		embed.Fields = []*discordgo.MessageEmbedField{{
			Name:   lang.Get(tp+"msg.set.date", lang.FallbackLang()),
			Value:  b.String(),
			Inline: true,
		}, {
			Name:   lang.Get(tp+"msg.next", lang.FallbackLang()),
			Value:  fmt.Sprintf("<t:%d:R>%s", b.NextUnix(), age),
			Inline: true,
		}}
		embed.Color = 0x00FF00
	}

	if b.Visible {
		cmd.ReplyEmbed(embed)
	} else {
		util.AddReplyHiddenField(embed)
		cmd.ReplyHiddenEmbed(embed)
	}
}

// seperate handler for an update of the birthday
func (cmd subcommandSet) handleUpdate(b birthdayEntry, e *discordgo.MessageEmbed) error {
	before, err := cmd.updateBirthday(b)
	if err != nil {
		return err
	}

	if b == before {
		var age string
		if b.Year > 0 {
			age = fmt.Sprintf(" (%d)", b.Age()+1)
		}

		e.Description = lang.Get(tp+"msg.set.update.no_changes", lang.FallbackLang())
		e.Fields = []*discordgo.MessageEmbedField{{
			Name:   lang.Get(tp+"msg.set.date", lang.FallbackLang()),
			Value:  b.String(),
			Inline: true,
		}, {
			Name:   lang.Get(tp+"msg.next", lang.FallbackLang()),
			Value:  fmt.Sprintf("<t:%d:R>%s", b.NextUnix(), age),
			Inline: true,
		}}
		e.Color = 0x696969
		return nil
	}

	e.Description = lang.Get(tp+"msg.set.update", lang.FallbackLang())
	e.Color = 0xfcb100

	const (
		DAY   = 1 << iota // when day is changed
		MONTH             // when month is changed
		YEAR              // when year is changed

		NODAY   = MONTH | YEAR       // when month and year is changed
		NOMONTH = DAY | YEAR         // when day and year is changed
		NOYEAR  = DAY | MONTH        // when day and month is changed
		ALL     = DAY | MONTH | YEAR // when day month and year is changed
	)

	// bit field of 4 bits to determin which values have changed
	//  1st (LSB) => day changed
	//  2nd       => month changed
	//  3rd       => year changed
	//  4th (MSB) => visibility changed
	var changedBits int = util.Btoi(before.Day != b.Day) |
		util.Btoi(before.Month != b.Month)<<1 |
		util.Btoi(before.Year != b.Year)<<2

	f := &discordgo.MessageEmbedField{Inline: false}
	switch changedBits {
	// set field when only day is changed
	case DAY:
		f.Name = lang.Get(tp+"msg.set.update.day", lang.FallbackLang())
		f.Value = fmt.Sprintf("%d -> %d", before.Day, b.Day)
	// set field when only month is changed
	case MONTH:
		f.Name = lang.Get(tp+"msg.set.update.month", lang.FallbackLang())
		mNameBefore := lang.GetSliceElement(tp+"month", before.Month-1, lang.FallbackLang())
		mName := lang.GetSliceElement(tp+"month", b.Month-1, lang.FallbackLang())
		f.Value = fmt.Sprintf("%s -> %s", mNameBefore, mName)
	// set field when only year is changed
	case YEAR:
		if before.Year == 0 {
			f.Name = lang.Get(tp+"msg.set.update.year.add", lang.FallbackLang())
			f.Value = fmt.Sprintf("%d", b.Year)
		} else if b.Year == 0 {
			f.Name = lang.Get(tp+"msg.set.update.year.remove", lang.FallbackLang())
			wasYear := lang.Get(tp+"msg.set.update.year.was", lang.FallbackLang())
			f.Value = fmt.Sprintf(wasYear, before.Year)
		} else {
			f.Name = lang.Get(tp+"msg.set.update.year", lang.FallbackLang())
			f.Value = fmt.Sprintf("%d -> %d", before.Year, b.Year)
		}
	// set field when any two or all three are changed
	case NOYEAR, NOMONTH, NODAY, ALL:
		f.Name = lang.Get(tp+"msg.set.update.date", lang.FallbackLang())
		f.Value = fmt.Sprintf("%s -> %s", before, b)
		f.Inline = true
	// set field when all three are remain the same (only visibility is changed)
	default:
		f.Name = lang.Get(tp+"msg.set.update.date.unchanged", lang.FallbackLang())
		f.Value = b.String()
		f.Inline = true
	}
	e.Fields = []*discordgo.MessageEmbedField{f}

	if !f.Inline {
		util.AddEmbedField(e,
			lang.Get(tp+"msg.set.date", lang.FallbackLang()),
			b.String(),
			true,
		)
	}

	var age string
	if b.Year > 0 {
		age = fmt.Sprintf(" (%d)", b.Age()+1)
	}

	util.AddEmbedField(e,
		lang.Get(tp+"msg.next", lang.FallbackLang()),
		fmt.Sprintf("<t:%d:R>%s", b.NextUnix(), age),
		true,
	)

	if before.Visible != b.Visible {
		var visibility string
		if b.Visible {
			key := tp + "msg.set.update.visibility.true"
			visibility = lang.Get(key, lang.FallbackLang())
		} else {
			key := tp + "msg.set.update.visibility.false"
			visibility = lang.Get(key, lang.FallbackLang())

			mentionCmd := util.MentionCommand(tp+"base", tp+"option.remove")
			visibility = fmt.Sprintf(visibility, mentionCmd)
		}
		util.AddEmbedField(e,
			lang.Get(tp+"msg.set.update.visibility", lang.FallbackLang()),
			visibility,
			false,
		)
	}

	return nil
}
