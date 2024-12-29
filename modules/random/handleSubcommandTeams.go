package random

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
)

// The set subcommand. Used when executing the slash-command "/random teams".
type subcommandTeams struct {
	randomBase
	*Chat
	data *discordgo.ApplicationCommandInteractionDataOption

	members    *discordgo.ApplicationCommandInteractionDataOption // required
	teamSize   *discordgo.ApplicationCommandInteractionDataOption // optional
	teamAmount *discordgo.ApplicationCommandInteractionDataOption // optional
}

func (rb randomBase) subcommandTeams() subcommandTeams {
	return subcommandTeams{randomBase: rb}
}

// Constructor for subcommandTeams, the struct for the slash-command "/random teams".
func (cmd *Chat) subcommandTeams() subcommandTeams {
	var subcommand *discordgo.ApplicationCommandInteractionDataOption
	if cmd.Interaction != nil {
		subcommand = cmd.Interaction.ApplicationCommandData().Options[0]
	}
	return subcommandTeams{
		randomBase: cmd.randomBase,
		Chat:       cmd,
		data:       subcommand,
	}
}

func (cmd subcommandTeams) appCmd() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		cmd.optionMembers(),
		cmd.optionTeamSize(),
		cmd.optionTeamAmount(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.teams"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.teams"),
		Description:              lang.GetDefault(tp + "option.teams.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.teams.description"),
		Options:                  options,
	}
}

func (cmd subcommandTeams) optionMembers() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionRole,
		Name:                     lang.GetDefault(tp + "option.teams.option.members"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.teams.option.members"),
		Description:              lang.GetDefault(tp + "option.teams.option.members.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.teams.option.members.description"),
		Required:                 true,
	}
}

func (cmd subcommandTeams) optionTeamSize() *discordgo.ApplicationCommandOption {
	minValueTwo := float64(2)
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.teams.option.team_size"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.teams.option.team_size"),
		Description:              lang.GetDefault(tp + "option.teams.option.team_size.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.teams.option.team_size.description"),
		Required:                 false,
		MinValue:                 &minValueTwo,
	}
}

func (cmd subcommandTeams) optionTeamAmount() *discordgo.ApplicationCommandOption {
	minValueOne := float64(1)
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.teams.option.team_amount"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.teams.option.team_amount"),
		Description:              lang.GetDefault(tp + "option.teams.option.team_amount.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.teams.option.team_amount.description"),
		Required:                 false,
		MinValue:                 &minValueOne,
	}
}

func (cmd subcommandTeams) handle() {
	for _, opt := range cmd.data.Options {
		switch opt.Name {
		case lang.GetDefault(tp + "option.teams.option.members"):
			cmd.members = opt
		case lang.GetDefault(tp + "option.teams.option.team_size"):
			cmd.teamSize = opt
		case lang.GetDefault(tp + "option.teams.option.team_amount"):
			cmd.teamAmount = opt
		}
	}

	if cmd.teamSize == nil && cmd.teamAmount == nil {
		cmd.ReplyHidden(lang.GetDefault(tp + "msg.teams.missing_option"))
		return
	} else if cmd.teamSize != nil && cmd.teamAmount != nil {
		cmd.ReplyHidden(lang.GetDefault(tp + "msg.teams.multiple_options"))
		return
	}

	var (
		memberRole = cmd.members.RoleValue(cmd.Session, cmd.Interaction.GuildID)
		teamSize   int
		teamAmount int
	)
	if cmd.teamSize != nil {
		teamSize = int(cmd.teamSize.IntValue())
	} else {
		teamAmount = int(cmd.teamAmount.IntValue())
	}

	members, err := cmd.getMembersWithRole(memberRole.ID)
	if err != nil {
		log.Printf("ERROR: could not get members with role '%s/%s' (%s): %+v", cmd.Interaction.GuildID, memberRole.ID, memberRole.Name, err)
		cmd.ReplyError()
	}

	data := &discordgo.InteractionResponseData{}
	if cmd.teamSize != nil {
		data = cmd.splitTeamsSize(members, teamSize)
	} else {
		data = cmd.splitTeamsN(members, teamAmount)
	}

	cmd.ReplyComplex(data)
}

func (cmd subcommandTeams) handleComponent(ids []string) {
}

// splitTeamsSize splits the members into teams of a maximum size teamSize.
//
// The last team might be smaller.
func (cmd subcommandTeams) splitTeamsSize(members []*discordgo.Member, teamSize int) (data *discordgo.InteractionResponseData) {
	data = &discordgo.InteractionResponseData{}

	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	var teams [][]*discordgo.Member
	for i := 0; i < len(members); i += teamSize {
		end := i + teamSize
		if end > len(members) {
			end = len(members)
		}
		teams = append(teams, members[i:end])
	}
	data.Embeds = teamsEmbed(cmd.Session, teams)

	resplitButton := util.CreateButtonComponent(
		fmt.Sprintf("random.teams.resplit_size.%d", teamSize),
		"",
		discordgo.PrimaryButton,
		util.GetConfigComponentEmoji("random.teams.resplit_size"))
	data.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{resplitButton}}}

	return data
}

// splitTeamsN splits the members into n teams.
func (cmd subcommandTeams) splitTeamsN(members []*discordgo.Member, n int) (data *discordgo.InteractionResponseData) {
	data = &discordgo.InteractionResponseData{}

	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	if n > len(members) {
		n = len(members)
	}
	var teams [][]*discordgo.Member = make([][]*discordgo.Member, n)
	for i, member := range members {
		teams[i%n] = append(teams[i%n], member)
	}
	data.Embeds = teamsEmbed(cmd.Session, teams)

	resplitButton := util.CreateButtonComponent(
		"random.teams.resplit_amount",
		"",
		discordgo.PrimaryButton,
		util.GetConfigComponentEmoji("random.teams.resplit_amount"))
	data.Components = []discordgo.MessageComponent{discordgo.ActionsRow{Components: []discordgo.MessageComponent{resplitButton}}}

	return data
}

func (cmd subcommandTeams) getMembersWithRole(roleID string) ([]*discordgo.Member, error) {
	var membersWithRole []*discordgo.Member
	var after string

	for {
		members, err := cmd.Session.GuildMembers(cmd.Interaction.GuildID, after, 1000)
		if err != nil {
			return nil, err
		}
		if len(members) == 0 {
			break
		}

		for _, member := range members {
			if util.ContainsString(member.Roles, roleID) {
				membersWithRole = append(membersWithRole, member)
			}
		}

		after = members[len(members)-1].User.ID
	}

	return membersWithRole, nil
}

// teamsEmbed returns one or more embeds listing the given teams.
func teamsEmbed(s *discordgo.Session, teams [][]*discordgo.Member) (embeds []*discordgo.MessageEmbed) {
	embeds = util.SplitToEmbedFields(s, teams, 0xFFD700, tp+"display", teamEmbed)
	embeds[0].Title = lang.GetDefault(tp + "msg.teams.title")

	if len(embeds[0].Fields) == 1 {
		embeds[0].Description = embeds[0].Fields[0].Value
		embeds[0].Fields = nil
	}

	return embeds
}

// teamEmbed returns the given team as an embed field.
//
// i is the team number (0-indexed) used for the field name.
func teamEmbed(team []*discordgo.Member, i int) *discordgo.MessageEmbedField {
	var value string
	for i, member := range team {
		value += fmt.Sprintf("%d. %s\n", i, member.Mention())
	}

	return &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf(lang.GetDefault(tp+"msg.teams.team"), i+1),
		Value:  value,
		Inline: true,
	}
}
