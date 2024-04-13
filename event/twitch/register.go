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
	"cake4everybot/tools/streamelements"

	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

// Register is setting up the twitch bot. Like joining channels and other stuff that is available
// after the bot is connected
func Register(t *twitchgo.Session) {
	channels := viper.GetStringSlice("twitch.channels")
	for _, channel := range channels {
		t.SendCommandf("JOIN #%s", channel)
	}
	log.Printf("Channel list set to %v\n", channels)

	se = streamelements.New(viper.GetString("streamelements.token"))
}
