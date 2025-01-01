package random

import (
	"math/rand/v2"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/util"
)

// The set subcommand. Used when executing the slash-command "/random coin".
type subcommandCoin struct {
	randomBase
	*Chat
	data *discordgo.ApplicationCommandInteractionDataOption
}

func (rb randomBase) subcommandCoin() subcommandCoin {
	return subcommandCoin{randomBase: rb}
}

// Constructor for subcommandCoin, the struct for the slash-command "/random coin".
func (cmd *Chat) subcommandCoin() subcommandCoin {
	var subcommand *discordgo.ApplicationCommandInteractionDataOption
	if cmd.Interaction != nil {
		subcommand = cmd.Interaction.ApplicationCommandData().Options[0]
	}
	return subcommandCoin{
		randomBase: cmd.randomBase,
		Chat:       cmd,
		data:       subcommand,
	}
}

func (cmd subcommandCoin) appCmd() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.coin"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.coin"),
		Description:              lang.GetDefault(tp + "option.coin.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.coin.description"),
	}
}

func (cmd subcommandCoin) handle() {
	cmd.ReplyComplex(cmd.flip())
}

func (cmd subcommandCoin) handleComponent(ids []string) {
	switch id := util.ShiftL(ids); id {
	case "reflip":
		if cmd.originalAuthor.ID != cmd.user.ID {
			cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.error.not_author"), cmd.originalAuthor.Mention())
			return
		}
		cmd.ReplyComplexUpdate(cmd.flip())
		return
	default:
		log.Printf("Unknown component interaction ID in subcommand coin: %s %s", id, ids)
	}
}

func (cmd subcommandCoin) flip() (data *discordgo.InteractionResponseData) {
	data = &discordgo.InteractionResponseData{}

	emoji, err := util.GetConfigEmoji(cmd.Session, "random.coin.flip")
	if err != nil {
		log.Printf("ERROR: could not get emoji: %+v", err)
		cmd.ReplyError()
		return
	}
	data.Content = emoji.MessageFormat()

	reflipButton := util.CreateButtonComponent(
		"random.coin.reflip",
		"",
		discordgo.PrimaryButton,
		util.GetConfigComponentEmoji("random.coin.reflip"))
	reflipButton.Disabled = true
	data.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{reflipButton}}}

	go func() {
		time.Sleep(2 * time.Second)
		reflipButton.Disabled = false
		defer cmd.Session.InteractionResponseEdit(cmd.Interaction.Interaction, util.MessageComplexWebhookEdit(data))

		side := "heads"
		if rand.IntN(2) == 1 {
			side = "tails"
		}
		var emoji *discordgo.Emoji
		emoji, err = util.GetConfigEmoji(cmd.Session, "random.coin."+side)
		if err != nil {
			log.Printf("Warning: could not get emoji: %+v", err)
			data.Content = side
			return
		}
		data.Content = emoji.MessageFormat()
	}()

	return data
}
