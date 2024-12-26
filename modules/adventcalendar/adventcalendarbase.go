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

package adventcalendar

import (
	"cake4everybot/logger"
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => adventcalendar
	tp = "discord.command.adventcalendar."
)

var log = logger.New("Advent")

type adventcalendarBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}
