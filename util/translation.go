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

package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
)

// TranslateLocalization returns a pointer to a map of all translations for the given key from
// discord languages that are loaded in the lang package.
func TranslateLocalization(key string) *map[discordgo.Locale]string {
	translateMap := map[discordgo.Locale]string{}
	for locale := range discordgo.Locales {
		if !lang.IsLoaded(string(locale)) {
			continue
		}
		translateMap[locale] = lang.Get(key, string(locale))
	}
	return &translateMap
}
