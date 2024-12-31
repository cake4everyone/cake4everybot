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

package lang

import (
	"strings"

	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/spf13/viper"
)

var log = logger.New("Data/Lang")
var langsMap = map[string]*viper.Viper{}

// Unify takes and returns a string wich defines a language, i.e.
// 'en_us', and changes it to a uniform format.
//
// It splits the given language at any '-' (dash), '_' (underscore),
// or ' ' (space) character, if any, and returns the first set of
// characters as lowercase.
//
//	Unify("en_us") = "en"
//	Unify("EN gb") = "en"
//	Unify("de")    = "de"
//
// This function is called on every lang input internally, so calling
// it on a lang name before passing it to a function is pointless.
func Unify(lang string) string {
	lang = strings.ToLower(lang)
	split := strings.IndexAny(lang, "-_ ")
	if split == -1 {
		return lang
	}
	return lang[:split]
}

// Load loads (and reloads) the language files defined in the global
// config.
func Load() {
	log.Println("Loading languages...")

	// clear all existing languages
	for k := range langsMap {
		delete(langsMap, k)
	}

	// read language files defined in global config
	for _, langName := range viper.GetStringSlice("languages") {
		if langName == "" {
			continue
		}

		langName = Unify(langName)

		// skip duplicates
		if _, ok := langsMap[langName]; ok {
			continue
		}

		log.Printf("Loading language '%s'...\n", langName)

		lang := viper.New()
		lang.SetConfigType("yaml")
		lang.AddConfigPath("data/lang")
		lang.SetConfigName(langName)

		err := lang.ReadInConfig()
		if err != nil {
			log.Printf("WARNING: Could not load language '%s': %v\n", langName, err)
			continue
		}

		langsMap[langName] = lang
	}

	if len(langsMap) == 0 {
		log.Fatalln("Could not load languages: needs at least one loaded language!")
	}

	log.Printf("Loaded %d language(s), with %s beeing the fallback language!\n", len(langsMap), FallbackLang())
}

// FallbackLang returns the bots fallback language, which is the
// first loaded language string in the 'languages' list from the
// global config.
//
// When 'languages' is empty or only has empty string entries, it
// fails and calls os.Exit(1).
func FallbackLang() string {
	for _, langName := range viper.GetStringSlice("languages") {
		if langName == "" {
			continue
		}

		langName = Unify(langName)
		if _, ok := langsMap[langName]; !ok {
			continue
		}
		return langName
	}

	log.Fatalln("ERROR: No languages in config")
	return ""
}

// Get returns the configured translation for key in the given
// language lang.
//   - If lang is not a loaded language, Get translates key with the
//     fallback language.
//   - If key either does not exits or is an empty string, Get
//     translates key with the fallback language.
//   - However, if lang already is the fallback language in one of the
//     cases above, Get returns key instead.
//
// In all three of these 'fail cases', Get will print a warning
// message in the log
func Get(key, lang string) string {
	if len(langsMap) == 0 {
		log.Println()
		log.Printf("ERROR: Tried to get translation, but no language loaded\n")
		log.Println()
		return key
	}

	lang = Unify(lang)

	v, ok := langsMap[lang]
	fLang := FallbackLang()
	if !ok {
		if lang == fLang {
			log.Println()
			log.Printf("ERROR: Tried to get key from fallback language ('%s'), but its not load\n", fLang)
			log.Println()
			return key
		}
		log.Printf("WARNING: language '%s' is not loaded, using '%s' as fallback instead\n", lang, fLang)
		return Get(key, fLang)
	}

	val := v.GetString(key)
	if val != "" {
		return val
	}

	if lang == fLang {
		log.Printf("WARNING: key '%s' is not defined in fallback language '%s'\n", key, lang)
		return key
	}
	log.Printf("WARNING: key '%s' is not defined in language '%s', using '%s' as fallback instead\n", key, lang, fLang)
	return Get(key, fLang)
}

// GetDefault is like Get, but with FallbackLang as language
func GetDefault(key string) string {
	return Get(key, FallbackLang())
}

// GetSliceElement returns the configured translation for index i in the
// list at key in the given language lang.
//   - If lang is not a loaded language, Get translates key with the
//     fallback language.
//   - If key either does not exits or is an empty string, Get
//     translates key with the fallback language.
//   - However, if lang already is the fallback language in one of the
//     cases above, Get returns key instead.
//   - If the list at key contains fewer items than needed, Get
//     returns key instead
//
// In all four of these 'fail cases', Get will print a warning
// message in the log
func GetSliceElement(key string, i int, lang string) string {
	if len(langsMap) == 0 {
		log.Println()
		log.Printf("ERROR: Tried to get translation, but no language loaded\n")
		log.Println()
		return key
	}

	lang = Unify(lang)

	v, ok := langsMap[lang]
	fLang := FallbackLang()
	if !ok {
		if lang == fLang {
			log.Println()
			log.Printf("ERROR: Tried to get key from fallback language ('%s'), but its not load\n", fLang)
			log.Println()
			return key
		}
		log.Printf("WARNING: language '%s' is not loaded, using '%s' as fallback instead\n", lang, fLang)
		return GetSliceElement(key, i, fLang)
	}

	s := v.GetStringSlice(key)
	if len(s) <= i {
		log.Printf("WARNING: tried to get index %d from key '%s' in lang '%s', but it has only %d items\n", i, key, lang, len(s))
		return key
	}
	val := s[i]
	if val != "" {
		return val
	}

	if lang == fLang {
		log.Printf("WARNING: key '%s' is not defined in fallback language '%s'\n", key, lang)
		return key
	}
	log.Printf("WARNING: key '%s' is not defined in language '%s', using '%s' as fallback instead\n", key, lang, fLang)
	return GetSliceElement(key, i, fLang)
}

// GetSlice is similar sto GetSliceElement, but instead returns the hole sting slice without any
// checks.
func GetSlice(key string, lang string) []string {
	if len(langsMap) == 0 {
		log.Println()
		log.Printf("ERROR: Tried to get translation, but no language loaded\n")
		log.Println()
		return []string{key}
	}

	lang = Unify(lang)

	v, ok := langsMap[lang]
	fLang := FallbackLang()
	if !ok {
		if lang == fLang {
			log.Println()
			log.Printf("ERROR: Tried to get key from fallback language ('%s'), but its not load\n", fLang)
			log.Println()
			return []string{key}
		}
		log.Printf("WARNING: language '%s' is not loaded, using '%s' as fallback instead\n", lang, fLang)
		return GetSlice(key, fLang)
	}

	return v.GetStringSlice(key)
}

// GetLangs returns all loaded languages
func GetLangs() []string {
	langs := make([]string, 0, len(langsMap))
	for lang := range langsMap {
		langs = append(langs, lang)
	}
	return langs
}

// IsLoaded returns true when the given
func IsLoaded(lang string) bool {
	lang = Unify(lang)
	_, ok := langsMap[lang]
	return ok
}
