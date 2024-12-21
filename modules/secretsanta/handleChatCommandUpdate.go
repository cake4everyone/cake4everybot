package secretsanta

import (
	"cake4everybot/data/lang"
)

// handleSubcommandUpdate handles the functionality of the update subcommand
func (cmd Chat) handleSubcommandUpdate() {
	cmd.ReplyDeferedHidden()
	players, err := cmd.getPlayers()
	if err != nil {
		cmd.ReplyError()
		return
	}

	var failedToSend string
	for _, p := range players {
		if _, _, ok := cmd.updateInviteMessage(p); !ok {
			failedToSend += "\n- " + p.Mention()
		}
	}

	if failedToSend != "" {
		cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.cmd.update.error"), failedToSend)
		return
	}
	cmd.ReplyHiddenf(lang.GetDefault(tp+"msg.cmd.update.success"), len(players))
}
