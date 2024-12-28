package random

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
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
		cmd.ReplyComplexUpdate(cmd.flip())
		return
	default:
		log.Printf("Unknown component interaction ID in subcommand coin: %s %s", id, ids)
	}
}

func (cmd subcommandCoin) flip() (m *discordgo.InteractionResponseData) {
	m = &discordgo.InteractionResponseData{}
	side := "heads"
	if rand.IntN(2) == 1 {
		side = "tails"
	}
	m.Content = util.GetConfigEmoji("random.coin." + side).MessageFormat()

	reflipButton := util.CreateButtonComponent(
		"random.coin.reflip",
		"",
		discordgo.PrimaryButton,
		util.GetConfigComponentEmoji("random.coin.reflip"))
	m.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{reflipButton}}}

	return m
}
