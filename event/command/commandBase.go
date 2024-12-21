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

package command

import (
	"cake4everybot/modules/adventcalendar"
	"cake4everybot/modules/birthday"
	"cake4everybot/modules/info"
	"cake4everybot/modules/secretsanta"
	"cake4everybot/util"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command is an interface wrapper for all commands. Including chat-comamnds (slash-commands),
// message-commands, and user-commands.
type Command interface {
	// Definition of a command.
	// E.g., name, description, options, subcommands.
	AppCmd() *discordgo.ApplicationCommand

	// Function of a command.
	// All things that should happen at execution.
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)

	// Sets the registered command ID for internal uses after uploading to discord
	SetID(id string)
	// Gets the registered command ID
	GetID() string
}

// CommandMap holds all active commands. It maps them from a unique name identifier to the
// corresponding Command.
//
// Here the name is used, because Discord uses the name too to identify seperate commands. When a
// command is beeing registered with a name that already is beeing registerd as a command by this
// application (bot), then the new one will simply overwrite it and automatically ungerister the old
// one.
var CommandMap map[string]Command

func init() {
	CommandMap = make(map[string]Command)
}

// Register registers all application commands
func Register(s *discordgo.Session, guildID string) error {

	// This is the list of commands to use. Add a command via simply appending the struct (which
	// must implement the Command interface) to the list, i.e.:
	//
	// commandsList = append(commandsList, command.MyCommand{})
	var commandsList []Command

	// chat (slash) commands
	commandsList = append(commandsList, &birthday.Chat{})
	commandsList = append(commandsList, &info.Chat{})
	commandsList = append(commandsList, &adventcalendar.Chat{})
	commandsList = append(commandsList, &secretsanta.Chat{})
	commandsList = append(commandsList, &secretsanta.MsgCmd{})
	// messsage commands
	// user commands
	commandsList = append(commandsList, &birthday.UserShow{})

	// early return when there're no commands to add, and remove all previously registered commands
	if len(commandsList) == 0 {
		removeUnusedCommands(s, guildID, nil)
		return nil
	}

	// make an array of ApplicationCommands and perform a bulk change using it
	appCommandsList := make([]*discordgo.ApplicationCommand, 0, len(commandsList))
	for _, cmd := range commandsList {
		appCommandsList = append(appCommandsList, cmd.AppCmd())
		CommandMap[cmd.AppCmd().Name] = cmd
	}
	commandNames := make([]string, 0, len(CommandMap))
	for k := range CommandMap {
		commandNames = append(commandNames, k)
	}

	log.Printf("Adding used commands: [%s]...\n", strings.Join(commandNames, ", "))
	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, appCommandsList)
	if err != nil {
		return fmt.Errorf("failed on bulk overwrite commands: %v", err)
	}

	for _, cmd := range createdCommands {
		CommandMap[cmd.Name].SetID(cmd.ID)
	}

	removeUnusedCommands(s, guildID, createdCommands)

	// set the utility map
	cmdIDMap := make(map[string]string)
	for k, v := range CommandMap {
		cmdIDMap[k] = v.GetID()
	}
	util.SetCommandMap(cmdIDMap)

	return err
}

func removeUnusedCommands(s *discordgo.Session, guildID string, createdCommands []*discordgo.ApplicationCommand) {
	allRegisteredCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		log.Printf("Error while removing unused commands: Could not get registered commands from guild '%s'. Err: %v\n", guildID, err)
		return
	}
	newCmdIds := make(map[string]bool)
	for _, cmd := range createdCommands {
		newCmdIds[cmd.ID] = true
	}

	// Find unused commands by iterating over all commands and check if the ID is in the currently registered commands. If not, remove it.
	for _, cmd := range allRegisteredCommands {
		if !newCmdIds[cmd.ID] {
			err = s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
			if err != nil {
				log.Printf("Error while removing unused commands: Could not delete comand '%s' ('/%s'). Err: %v\n", cmd.ID, cmd.Name, err)
				continue
			}
			log.Printf("Removed unused command: '/%s'\n", cmd.Name)
		}
	}

}
