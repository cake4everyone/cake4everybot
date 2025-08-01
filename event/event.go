// Copyright 2022-2023 Kesuaheli
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
	"github.com/cake4everyone/cake4everybot/event/command"
	"github.com/cake4everyone/cake4everybot/event/component"
	"github.com/cake4everyone/cake4everybot/event/modal"
	"github.com/cake4everyone/cake4everybot/event/twitch"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/kesuaheli/twitchgo"
)

var log = logger.New("Event")

// PostRegister registers all events, like commands after the bots are started.
func PostRegister(dc *discordgo.Session, t *twitchgo.Session, guildID string) error {
	err := command.Register(dc, guildID)
	if err != nil {
		return err
	}
	component.Register()
	modal.Register()

	twitch.Register(t)

	return nil
}

// AddListeners adds all event handlers to the given bots.
func AddListeners(dc *discordgo.Session, t *twitchgo.Session, webChan chan struct{}) {
	dc.AddHandler(handleInteractionCreate)
	addVoiceStateListeners(dc)

	t.OnChannelCommandMessage("", false, twitch.HandleGeneralCommand)
	t.OnChannelMessage(twitch.MessageHandler)

	addYouTubeListeners(dc)
	addTwitchListeners(dc, t, webChan)
	addScheduledTriggers(dc, t, webChan)
}
