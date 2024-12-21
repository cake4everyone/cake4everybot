package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"encoding/json"
	"fmt"
	logger "log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => adventcalendar
	tp = "discord.command.secretsanta."
)

var log = logger.New(logger.Writer(), "[SecretSanta] ", logger.LstdFlags|logger.Lmsgprefix)

type secretSantaBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

// getPlayers returns the list of players for the current guild. If it is the first time, it loads
// the players from the file or creates an empty file.
func (ssb secretSantaBase) getPlayers() (map[string]*player, error) {
	if allPlayers != nil {
		return allPlayers[ssb.Interaction.GuildID], nil
	}

	log.Println("First time getting players. Loading from file...")
	playersPath := viper.GetString("event.secretsanta.players")
	playersData, err := os.ReadFile(playersPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read players file: %v", err)
		}
		allPlayers = make(AllPlayers)
		playersData, err = json.Marshal(allPlayers)
		if err != nil {
			return nil, fmt.Errorf("marshal players file: %v", err)
		}
		err = os.WriteFile(playersPath, playersData, 0644)
		if err != nil {
			return nil, fmt.Errorf("write players file: %v", err)
		}
		log.Printf("Created players file: %s\n", playersPath)
		return map[string]*player{}, nil
	}
	allPlayersUnresolved := AllPlayersUnresolved{}
	err = json.Unmarshal(playersData, &allPlayersUnresolved)
	if err != nil {
		return nil, fmt.Errorf("unmarshal players file: %v", err)
	}
	err = allPlayersUnresolved.Resolve(ssb.Session)
	if err != nil {
		return nil, fmt.Errorf("resolve players file: %v", err)
	}
	log.Printf("Got %d guilds from file", len(allPlayers))

	return allPlayers[ssb.Interaction.GuildID], nil
}

// getSantaForPlayer returns the santa player of the given player for the current guild i.e. the
// player whose match is the given player.
//
// It panics if the player doesn't exist.
func (ssb secretSantaBase) getSantaForPlayer(playerID string) *player {
	for _, player := range allPlayers[ssb.Interaction.GuildID] {
		if player.Match.User.ID == playerID {
			return player
		}
	}
	panic(fmt.Sprintf("tried to get santa for player that doesn't exist: '%s' in guild '%s'", playerID, ssb.Interaction.GuildID))
}

// setPlayers sets the players for the current guild.
func (ssb secretSantaBase) setPlayers(players map[string]*player) (err error) {
	if _, err = ssb.getPlayers(); err != nil {
		return err
	}

	allPlayers[ssb.Interaction.GuildID] = players
	playersData, err := json.Marshal(allPlayers)
	if err != nil {
		return fmt.Errorf("marshal players file: %v", err)
	}
	err = os.WriteFile(viper.GetString("event.secretsanta.players"), playersData, 0644)
	if err != nil {
		return fmt.Errorf("write players file: %v", err)
	}
	return nil
}

// inviteMessage returns the message to send to the player to invite them to play.
func (ssb secretSantaBase) inviteMessage(p *player) *discordgo.MessageSend {
	var components []discordgo.MessageComponent
	components = append(components, util.CreateButtonComponent(
		fmt.Sprintf("secretsanta.invite.show_match.%s", ssb.Interaction.GuildID),
		lang.GetDefault(tp+"msg.invite.button.show_match"),
		discordgo.PrimaryButton,
		util.GetConfigComponentEmoji("secretsanta.invite.show_match"),
	))
	santaPlayer := ssb.getSantaForPlayer(p.User.ID)
	if sendPackageState := santaPlayer.SendPackage; sendPackageState == 0 {
		components = append(components, util.CreateButtonComponent(
			fmt.Sprintf("secretsanta.invite.set_address.%s", ssb.Interaction.GuildID),
			lang.GetDefault(tp+"msg.invite.button.set_address"),
			discordgo.SecondaryButton,
			util.GetConfigComponentEmoji("secretsanta.invite.set_address"),
		))
	} else if sendPackageState == 1 {
		if santaPlayer.PackageTracking != "" {
			components = append(components, util.CreateButtonComponent(
				fmt.Sprintf("secretsanta.invite.show_package_tracking.%s", ssb.Interaction.GuildID),
				lang.GetDefault(tp+"msg.invite.button.show_package_tracking"),
				discordgo.SecondaryButton,
				util.GetConfigComponentEmoji("secretsanta.invite.show_package_tracking"),
			))
		}
		components = append(components, util.CreateButtonComponent(
			fmt.Sprintf("secretsanta.invite.received_package.%s", ssb.Interaction.GuildID),
			lang.GetDefault(tp+"msg.invite.button.received_package"),
			discordgo.SuccessButton,
			util.GetConfigComponentEmoji("secretsanta.invite.received_package"),
		))
	} else if sendPackageState == 2 && santaPlayer.PackageTracking != "" {
		components = append(components, util.CreateButtonComponent(
			fmt.Sprintf("secretsanta.invite.show_package_tracking.%s", ssb.Interaction.GuildID),
			lang.GetDefault(tp+"msg.invite.button.show_package_tracking"),
			discordgo.SecondaryButton,
			util.GetConfigComponentEmoji("secretsanta.invite.show_package_tracking"),
		))
	}

	return &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{p.InviteEmbed(ssb.Session)},
		Components: []discordgo.MessageComponent{discordgo.ActionsRow{Components: components}},
	}
}

// updateInviteMessage updates the invite message for the player.
func (ssb secretSantaBase) updateInviteMessage(p *player) (DMChannel *discordgo.Channel, msg *discordgo.Message, ok bool) {
	DMChannel, err := ssb.Session.UserChannelCreate(p.User.ID)
	if err != nil {
		log.Printf("ERROR: could not create DM channel for user %s: %+v", p.User.ID, err)
		return DMChannel, msg, false
	}

	if p.MessageID == "" {
		msg, err = ssb.Session.ChannelMessageSendComplex(DMChannel.ID, ssb.inviteMessage(p))
		if err != nil {
			log.Printf("ERROR: could not send invite message for %s: %+v", p.DisplayName(), err)
			return DMChannel, msg, false
		}
		p.MessageID = msg.ID
	} else {
		msg, err = ssb.Session.ChannelMessageEditComplex(util.MessageComplexEdit(ssb.inviteMessage(p), DMChannel.ID, p.MessageID))
		if err != nil {
			log.Printf("ERROR: could not update bot message for %s '%s/%s': %+v", p.DisplayName(), ssb.Interaction.ChannelID, p.MessageID, err)
			return DMChannel, msg, false
		}
	}
	return DMChannel, msg, true
}

// player is a player in the secret santa game
type player struct {
	*discordgo.Member

	// Match is the matched player
	Match *player
	// Address is the address of the player
	Address string
	// MessageID is the message the bot sent to the player
	MessageID string
	// PendingNudge is true if the player has received a nugde from their santa and they haven't changed their
	// address yet i.e. the nudge is still pending.
	PendingNudge bool
	// SendPackage is the state of the package sent to the player.
	//
	//   0 = not sent
	//   1 = sent
	//   2 = sent and received by their partner
	SendPackage int
	// PackageTracking is a shipment tracking reference provided by the player for their partner
	PackageTracking string
}

// InviteEmbed returns an embed for the player to be sent by the bot.
func (player *player) InviteEmbed(s *discordgo.Session) (e *discordgo.MessageEmbed) {
	var matchValue, addressValue, sendPackageValue = "❌", "❌", "❓"
	if player != nil && player.Match.Address != "" {
		if player.Match.PendingNudge {
			matchValue = fmt.Sprintf("%s %s", "⌛", lang.GetDefault(tp+"msg.invite.nudge_match.pending"))
		} else {
			matchValue = "✅"
		}
	}
	if player != nil && player.Address != "" {
		if player.PendingNudge {
			addressValue = fmt.Sprintf("%s %s", "⚠️", lang.GetDefault(tp+"msg.invite.nudge_received"))
		} else {
			addressValue = "✅"
		}
	}
	if player != nil {
		switch player.SendPackage {
		case 0:
			sendPackageValue = lang.GetDefault(tp + "msg.invite.send_package.status.not_sent")
		case 1:
			sendPackageValue = lang.GetDefault(tp + "msg.invite.send_package.status.sent")
		case 2:
			sendPackageValue = lang.GetDefault(tp + "msg.invite.send_package.status.received")
		}
	}

	e = &discordgo.MessageEmbed{
		Title:       lang.GetDefault(tp + "msg.invite.title"),
		Description: lang.GetDefault(tp + "msg.invite.description"),
		Color:       0x690042,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   lang.GetDefault(tp + "msg.invite.set_address.match"),
				Value:  matchValue,
				Inline: true,
			},
			{
				Name:   lang.GetDefault(tp + "msg.invite.set_address"),
				Value:  addressValue,
				Inline: true,
			},
			{
				Name:   lang.GetDefault(tp + "msg.invite.send_package.status"),
				Value:  sendPackageValue,
				Inline: false,
			},
		},
	}
	util.SetEmbedFooter(s, tp+"display", e)
	return e
}

type playerUnresolved struct {
	MatchID         string `json:"match"`
	Address         string `json:"address"`
	MessageID       string `json:"message"`
	PendingNudge    bool   `json:"pending_nudge,omitempty"`
	SendPackage     int    `json:"send_package,omitempty"`
	PackageTracking string `json:"package_tracking,omitempty"`
}

// AllPlayers is a map from guild ID to a list of players
type AllPlayers map[string]map[string]*player

// allPlayers is the current state of all players.
// See [AllPlayers]
var allPlayers AllPlayers

// MarshalJSON implements json.Marshaler
func (allPlayers AllPlayers) MarshalJSON() ([]byte, error) {
	m := make(AllPlayersUnresolved)
	for guildID, players := range allPlayers {
		m[guildID] = make(map[string]*playerUnresolved)
		for userID, player := range players {
			var matchID string
			if player.Match != nil {
				matchID = player.Match.User.ID
			}
			m[guildID][userID] = &playerUnresolved{
				MatchID:         matchID,
				Address:         player.Address,
				MessageID:       player.MessageID,
				PendingNudge:    player.PendingNudge,
				SendPackage:     player.SendPackage,
				PackageTracking: player.PackageTracking,
			}
		}
	}
	return json.Marshal(m)
}

// AllPlayersUnresolved is a map from guild ID to a list of unresolved players.
// Unresolved players have no member but only an ID
type AllPlayersUnresolved map[string]map[string]*playerUnresolved

// Resolve resolves allPlayersUnresolved into allPlayers
func (allPlayersUnresolved AllPlayersUnresolved) Resolve(s *discordgo.Session) (err error) {
	allPlayers = make(AllPlayers)
	for guildID, unresolvedPlayers := range allPlayersUnresolved {
		allPlayers[guildID] = make(map[string]*player)
		for userID, up := range unresolvedPlayers {
			member, err := s.GuildMember(guildID, userID)
			if err != nil {
				return fmt.Errorf("failed to get guild member %s/%s: %v", guildID, userID, err)
			}
			allPlayers[guildID][userID] = &player{
				Member:          member,
				Match:           allPlayers[guildID][up.MatchID],
				Address:         up.Address,
				MessageID:       up.MessageID,
				PendingNudge:    up.PendingNudge,
				SendPackage:     up.SendPackage,
				PackageTracking: up.PackageTracking,
			}
		}
		for userID, rp := range allPlayers[guildID] {
			if rp.Match != nil {
				continue
			}
			rp.Match = allPlayers[guildID][unresolvedPlayers[userID].MatchID]
		}
	}
	return nil
}

var blacklist map[string][]string

// loadBlacklist loads the blacklist from the configured file path.
func loadBlacklist() (err error) {
	blacklistPath := viper.GetString("event.secretsanta.blacklist")
	blacklistData, err := os.ReadFile(blacklistPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("read blacklist file: %w", err)
		}
		if err = os.WriteFile(blacklistPath, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("write blacklist file: %w", err)
		}
		return nil
	}

	blacklist = make(map[string][]string)
	if err = json.Unmarshal(blacklistData, &blacklist); err != nil {
		return fmt.Errorf("unmarshal blacklist file: %w", err)
	}
	return nil
}

// derangementMatch matches the players in a way that no one gets matched to themselves.
func derangementMatch(players map[string]*player) (map[string]*player, error) {
	loadBlacklist()
	n := len(players)
	playersSlice := make([]*player, 0, n)
	for _, p := range players {
		p.Match = p
		playersSlice = append(playersSlice, p)
	}

	for i := 0; i < n-1; i++ {
		j := i + rand.Intn(n-i-1) + 1
		if !util.ContainsString(blacklist[playersSlice[i].User.ID], playersSlice[j].Match.User.ID) {
			log.Printf("%s would match blacklisted %s. Attempt to generate new match.", playersSlice[i].DisplayName(), playersSlice[j].Match.DisplayName())
			// HACK: very simple attempt to avoid blacklisted matches
			// just in case a match is blacklisted, generate a new match. If it is still
			// blacklisted, then it will be caught by the blacklist check below.
			j = i + rand.Intn(n-i-1) + 1
		}
		playersSlice[i].Match, playersSlice[j].Match = playersSlice[j].Match, playersSlice[i].Match
	}

	for id, blacklisted := range blacklist {
		if player, ok := players[id]; ok && util.ContainsString(blacklisted, player.Match.User.ID) {
			return nil, fmt.Errorf("'%s' has a blacklisted match: '%s'", player.DisplayName(), player.Match.DisplayName())
		}
	}

	return players, nil
}
