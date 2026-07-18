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

package reaction

import "github.com/bwmarrin/discordgo"

type Reaction interface {

	// All things that should happen when reacting to a message.
	AddHandler() func(s *discordgo.Session, i *discordgo.MessageReactionAdd)

	// All things that should happen when removing a reaction from a message.
	RemoveHandler() func(s *discordgo.Session, i *discordgo.MessageReactionRemove)
}

var ReactionMap map[string]Reaction

func init() {
	ReactionMap = make(map[string]Reaction)
}
