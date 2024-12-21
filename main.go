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

package main

import (
	"context"
	logger "log"
	"os/signal"
	"syscall"

	"cake4everybot/config"
	"cake4everybot/database"
	"cake4everybot/event"
	"cake4everybot/twitch"
	"cake4everybot/webserver"

	"github.com/bwmarrin/discordgo"
	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

var log = logger.New(logger.Writer(), "[MAIN] ", logger.LstdFlags|logger.Lmsgprefix)

const banner string = "\n" +
	"   ______      __        __ __  ______                                     \n" +
	"  / ____/___ _/ /_____  / // / / ____/   _____  _______  ______  ____  ___ \n" +
	" / /   / __ `/ //_/ _ \\/ // /_/ __/ | | / / _ \\/ ___/ / / / __ \\/ __ \\/ _ \\\n" +
	"/ /___/ /_/ / ,< /  __/__  __/ /___ | |/ /  __/ /  / /_/ / /_/ / / / /  __/\n" +
	"\\____/\\__,_/_/|_|\\___/  /_/ /_____/ |___/\\___/_/   \\__, /\\____/_/ /_/\\___/ \n" +
	"                                                  /____/                   \n" +
	"      ____  _                          __            ____        __        \n" +
	"     / __ \\(_)_____________  _________/ /           / __ )____  / /_       \n" +
	"    / / / / / ___/ ___/ __ \\/ ___/ __  /  ______   / __  / __ \\/ __/       \n" +
	"   / /_/ / (__  ) /__/ /_/ / /  / /_/ /  /_____/  / /_/ / /_/ / /_         \n" +
	"  /_____/_/____/\\___/\\____/_/   \\__,_/           /_____/\\____/\\__/         \n" +
	"\n" +
	"Version: v%s\n" +
	"%s\n" +
	"Copyright 2022-2023 Kesuaheli\n\n"

func init() {
	config.Load("config.yaml")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()
	webChan := make(chan struct{})

	log.Printf(banner, viper.GetString("version"), viper.GetString("discord.credits"))

	database.Connect()
	defer database.Close()

	// initializing Discord and Twitch bots
	discordBot, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.Fatalf("invalid discord bot parameters: %v", err)
	}

	discordBot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in to Discord as %s#%s\n", s.State.User.Username, s.State.User.Discriminator)
	})

	twitch.Connect()
	twitchBot := twitchgo.New(viper.GetString("twitch.name"), viper.GetString("twitch.token"))

	// adding listeners for events
	event.AddListeners(discordBot, twitchBot, webChan)

	// open connection and login to Discord and Twitch
	log.Println("Logging in to Discord")
	err = discordBot.Open()
	if err != nil {
		log.Fatalf("could not open the discord session: %v", err)
	}
	defer discordBot.Close()

	log.Println("Logging in to Twitch")
	err = twitchBot.Connect()
	if err != nil {
		log.Fatalf("could not open the twitch connection: %v", err)
	}
	defer twitchBot.Close()

	// register all events.
	err = event.PostRegister(discordBot, twitchBot, viper.GetString("discord.guildID"))
	if err != nil {
		log.Printf("Error registering events: %v\n", err)
	}

	log.Println("Starting webserver...")
	addr := ":8080"
	webserver.Run(addr, webChan)

	// Wait to end the bot
	log.Println("Press Ctrl+C to exit")
	<-ctx.Done()

	log.Println("\nGracefully shutting down. Byee")
}
