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

package config

import (
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/spf13/viper"
)

var log = logger.New("Config")

// Load loads the given configuration file as the global config. It
// also loads:
//   - the languages from lang.Load() (see cake4everybot/data/lang)
func Load(config string) {
	log.Println("Loading configuration file(s)...")
	log.Printf("Loading config '%s'\n", config)

	viper.AddConfigPath(".")
	viper.SetConfigFile(config)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Could not load config file '%s': %v", config, err)
	}

	names := viper.GetStringSlice("additionalConfigs")
	for _, n := range names {
		log.Printf("Loading additional config '%s'...\n", n)
		viper.SetConfigFile(n)
		err = viper.MergeInConfig()
		if err != nil {
			log.Printf("Counld not load additional config '%s': %v\n", n, err)
		}
	}

	log.Println("Loaded configuration file(s)!")
}
