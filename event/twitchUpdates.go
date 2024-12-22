// Copyright 2024 Kesuaheli
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

package event

import (
	"cake4everybot/event/twitch"
	webTwitch "cake4everybot/webserver/twitch"

	"github.com/bwmarrin/discordgo"
	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

func addTwitchListeners(s *discordgo.Session, t *twitchgo.Session) {
	webTwitch.SetDiscordSession(s)
	webTwitch.SetTwitchSession(t)
	webTwitch.SetDiscordChannelUpdateHandler(twitch.HandleChannelUpdate)
	webTwitch.SetDiscordStreamOnlineHandler(twitch.HandleStreamOnline)
	webTwitch.SetDiscordStreamOfflineHandler(twitch.HandleStreamOffline)

	channels := viper.GetStringSlice("announce.twitch")
	for _, channelID := range channels {
		webTwitch.SubscribeChannel(channelID)
	}
}
