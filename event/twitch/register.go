// Copyright 2023 Kesuaheli
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

package twitch

import (
	"github.com/cake4everyone/cake4everybot/event/twitch/timer"
	"github.com/cake4everyone/cake4everybot/tools/streamelements"
	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

// Register is setting up the twitch bot. Like joining channels and other stuff that is available
// after the bot is connected
func Register(t *twitchgo.Session) (err error) {
	channels := viper.GetStringSlice("twitch.channels")
	for _, channel := range channels {
		t.JoinChannel(channel)
	}
	log.Printf("Channel list set to %v\n", channels)
	err = timer.RegisterTimer(t, channels)
	if err != nil {
		return err
	}

	se = streamelements.New(viper.GetString("streamelements.token"))
	return nil
}
