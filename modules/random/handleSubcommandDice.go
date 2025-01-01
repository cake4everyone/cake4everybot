package random

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// The set subcommand. Used when executing the slash-command "/random dice".
type subcommandDice struct {
	randomBase
	*Chat
	data *discordgo.ApplicationCommandInteractionDataOption

	diceRange *discordgo.ApplicationCommandInteractionDataOption // optional
}

func (rb randomBase) subcommandDice() subcommandDice {
	return subcommandDice{randomBase: rb}
}

// Constructor for subcommandDice, the struct for the slash-command "/random dice".
func (cmd *Chat) subcommandDice() subcommandDice {
	var subcommand *discordgo.ApplicationCommandInteractionDataOption
	if cmd.Interaction != nil {
		subcommand = cmd.Interaction.ApplicationCommandData().Options[0]
	}
	return subcommandDice{
		randomBase: cmd.randomBase,
		Chat:       cmd,
		data:       subcommand,
	}
}

func (cmd subcommandDice) appCmd() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		cmd.optionRange(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.dice"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.dice"),
		Description:              lang.GetDefault(tp + "option.dice.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.dice.description"),
		Options:                  options,
	}
}

func (cmd subcommandDice) optionRange() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.dice.option.range"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.dice.option.range"),
		Description:              lang.GetDefault(tp + "option.dice.option.range.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.dice.option.range.description"),
		Required:                 false,
		MinValue:                 util.FloatTwo(),
	}
}

func (cmd subcommandDice) handle() {
	for _, opt := range cmd.data.Options {
		switch opt.Name {
		case lang.GetDefault(tp + "option.dice.option.range"):
			cmd.diceRange = opt
		}
	}
	diceRange := 6
	if cmd.diceRange != nil {
		diceRange = int(cmd.diceRange.IntValue())
	}
	cmd.ReplyComplex(cmd.roll(diceRange))
}

func (cmd subcommandDice) handleComponent(ids []string) {
	switch id := util.ShiftL(ids); id {
	case "reroll":
		if cmd.originalAuthor.ID != cmd.user.ID {
			cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.error.not_author"), cmd.originalAuthor.Mention())
			return
		}
		diceRange, _ := strconv.Atoi(util.ShiftL(ids))
		cmd.ReplyComplexUpdate(cmd.roll(diceRange))
		return
	default:
		log.Printf("Unknown component interaction ID in subcommand dice: %s %s", id, ids)
	}
}

func (cmd subcommandDice) roll(diceRange int) (data *discordgo.InteractionResponseData) {
	data = &discordgo.InteractionResponseData{}
	var err error

	if diceRange <= 6 {
		var emoji *discordgo.Emoji
		emoji, err = util.GetConfigEmoji(cmd.Session, "random.dice.rolling")
		if err != nil {
			log.Printf("ERROR: could not get emoji: %+v", err)
			cmd.ReplyError()
			return nil
		}
		data.Content = emoji.MessageFormat()
	} else {
		data.Embeds = util.SimpleEmbed(0xFF7D00, "...")
	}

	rerollButton := util.CreateButtonComponent(
		fmt.Sprintf("random.dice.reroll.%d", diceRange),
		"",
		discordgo.PrimaryButton,
		util.GetConfigComponentEmoji("random.dice.reroll"),
	)
	rerollButton.Disabled = true
	data.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{rerollButton}}}

	go func() {
		time.Sleep(2 * time.Second)
		rerollButton.Disabled = false
		defer cmd.Session.InteractionResponseEdit(cmd.Interaction.Interaction, util.MessageComplexWebhookEdit(data))

		diceResult := rand.IntN(diceRange) + 1
		if diceRange > 6 {
			data.Embeds = util.SimpleEmbedf(0xFF7D00, lang.GetDefault(tp+"msg.dice.roll"), diceResult)
			return
		}

		var emoji *discordgo.Emoji
		emoji, err = util.GetConfigEmoji(cmd.Session, fmt.Sprintf("random.dice.%d", diceResult))
		if err != nil {
			log.Printf("Warning: could not get emoji: %+v", err)
			data.Content = fmt.Sprintf(lang.GetDefault(tp+"msg.dice.roll"), diceResult)
			return
		}
		data.Content = emoji.MessageFormat()
	}()

	return data
}
