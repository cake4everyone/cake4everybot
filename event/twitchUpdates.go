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
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/event/twitch"
	"github.com/cake4everyone/cake4everybot/util"
	webTwitch "github.com/cake4everyone/cake4everybot/webserver/twitch"
	"github.com/kesuaheli/twitchgo"
)

func addTwitchListeners(s *discordgo.Session, t *twitchgo.Session, webChan chan struct{}) {
	webTwitch.SetDiscordSession(s)
	webTwitch.SetTwitchSession(t)
	webTwitch.SetDiscordChannelUpdateHandler(twitch.HandleChannelUpdate)
	webTwitch.SetDiscordStreamOnlineHandler(twitch.HandleStreamOnline)
	webTwitch.SetDiscordStreamOfflineHandler(twitch.HandleStreamOffline)

	err := util.ForAllPlatformIDs(database.AnnouncementPlatformTwitch, func(channelID string) {
		webTwitch.SubscribeChannel(channelID)
	})
	if err != nil {
		log.Printf("Error on subscribing to Twitch channels: %v", err)
	}

	go func() {
		<-webChan
		webTwitch.RefreshSubscriptions()
	}()
}
