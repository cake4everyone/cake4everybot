package secretsanta

// handleSubcommandShow handles the functionality of the show subcommand
func (cmd Chat) handleSubcommandShow() {
	players, err := cmd.getPlayers()
	if err != nil {
		cmd.ReplyError()
		return
	}

	var list string
	for _, p := range players {
		list += "- " +
			p.Mention() +
			" - (`" + p.User.Username + "`)" +
			"\n"
	}
	cmd.ReplyHiddenSimpleEmbed(0x690042, list)
}
