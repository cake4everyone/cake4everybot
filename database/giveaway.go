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

package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GiveawayEntry represents a giveaway entry from the database
type GiveawayEntry struct {
	// Identification of the entry
	UserID string
	// The current weight or number of tickets in this entry
	Weight int
	// The day of last entry. Useful to check when only one ticket per day is allowed.
	LastEntry time.Time
	// The platform the giveaway is for
	Platform Platform
	// The platform identifier e.g. the channel id for the platform
	PlatformID string
}

// ToEmbedField formats the giveaway entry to an discord message embed field.
func (e GiveawayEntry) ToEmbedField(s *discordgo.Session, totalTickets int) (f *discordgo.MessageEmbedField) {
	var name string
	if u, err := s.User(e.UserID); err != nil {
		log.Printf("Error on getting user '%s': %v", e.UserID, err)
		name = "???"
	} else {
		name = u.Username
	}

	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  fmt.Sprintf("<@%s>\n%d tickets\nChance: %.2f%%\nlast entry: <t:%d:R>", e.UserID, e.Weight, float64(e.Weight*100)/float64(totalTickets), e.LastEntry.Unix()),
		Inline: true,
	}
}

// GetGiveawayEntry gets the giveaway entry for the given user identifier, if their last entry was
// prefixed with prefix.
//
// If an error occours or it doesn't match prefix, an emtpy GiveawayEntry is returned instead.
func GetGiveawayEntry(prefix, userID string, platform Platform, platformID string) GiveawayEntry {
	var (
		weight      int
		lastEntryID string
	)
	err := QueryRow("SELECT weight,last_entry_id FROM giveaway WHERE id=? AND platform=? AND platform_id=?", userID, platform, platformID).Scan(&weight, &lastEntryID)
	if err == sql.ErrNoRows {
		return GiveawayEntry{UserID: userID, Weight: 0, Platform: platform, PlatformID: platformID}
	}
	if err != nil {
		log.Printf("Database failed to get giveaway entries for '%s': %v", userID, err)
		return GiveawayEntry{}
	}

	if lastEntryID == "" {
		return GiveawayEntry{UserID: userID, Weight: weight, Platform: platform, PlatformID: platformID}
	}

	dateValue, ok := strings.CutPrefix(lastEntryID, prefix+"-")
	if !ok {
		return GiveawayEntry{}
	}

	lastEntry, err := time.Parse(time.DateOnly, dateValue)
	if err != nil {
		log.Printf("could not convert last_entry_id '%s' to time: %v", lastEntryID, err)
		return GiveawayEntry{}
	}
	return GiveawayEntry{userID, weight, lastEntry, platform, platformID}
}

// DeleteGiveawayEntry deletes the giveaway entry for the given user identifier from the database.
//
// If an error occours it will be returned. However if no datbase entry matched it returns err ==
// nil, not err == sql.ErrNoRows. Because sql.ErrNoRows also results in the non-existence of the
// requested row and therefore is treated as a successful call.
func DeleteGiveawayEntry(userID string, platform Platform, platformID string) error {
	_, err := Exec("DELETE FROM giveaway WHERE id=? AND platform=? AND platform_id=?", userID, platform, platformID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// AddGiveawayWeight adds amount to the given user identifier.
//
// However if their last entry wasn't prefixed with prefix, their weight will be resetted and starts
// at amount. If you dont want it to be resetted check with GetGiveawayEntry first.
//
// If there was no error the modified entry is returned. If there was an error, an emtpy
// GiveawayEntry is returned instead.
func AddGiveawayWeight(prefix, userID string, amount int, platform Platform, platformID string) GiveawayEntry {
	var (
		weight      int
		lastEntryID string
		new         bool
	)
	err := QueryRow("SELECT weight,last_entry_id FROM giveaway WHERE id=? AND platform=? AND platform_id=?", userID, platform, platformID).Scan(&weight, &lastEntryID)
	if err == sql.ErrNoRows {
		new = true
	} else if err != nil {
		log.Printf("Database failed to get giveaway weight for '%s': %v", userID, err)
		return GiveawayEntry{}
	}

	// validate prefix
	if _, ok := strings.CutPrefix(lastEntryID, prefix+"-"); !ok {
		weight = 0
	}

	weight += amount
	dateValue := time.Now().Format(time.DateOnly)
	lastEntryID = fmt.Sprintf("%s-%s", prefix, dateValue)
	lastEntry, _ := time.Parse(time.DateOnly, dateValue)

	if new {
		_, err = Exec("INSERT INTO giveaway (id,weight,last_entry_id) VALUES (?,?,?)", userID, weight, lastEntryID)
		if err != nil {
			log.Printf("Database failed to insert giveaway for '%s': %v", userID, err)
			return GiveawayEntry{}
		}
		return GiveawayEntry{userID, weight, lastEntry, platform, platformID}
	}
	_, err = Exec("UPDATE giveaway SET weight=?,last_entry_id=? WHERE id=?", weight, lastEntryID, userID)
	if err != nil {
		log.Printf("Database failed to update weight (new: %d) for '%s': %v", weight, userID, err)
		return GiveawayEntry{}
	}
	return GiveawayEntry{userID, weight, lastEntry, platform, platformID}
}

// GetAllGiveawayEntries gets all giveaway entries that matches prefix.
func GetAllGiveawayEntries(prefix string, platform Platform, platformID string) []GiveawayEntry {
	rows, err := Query("SELECT id,weight,last_entry_id FROM giveaway WHERE platform=? AND platform_id=?", platform, platformID)
	if err == sql.ErrNoRows {
		return []GiveawayEntry{}
	} else if err != nil {
		log.Printf("ERROR: could not get entries from database: %v", err)
		return []GiveawayEntry{}
	}
	defer rows.Close()

	var entries []GiveawayEntry
	for rows.Next() {
		var (
			userID      string
			weight      int
			lastEntryID string
		)
		err = rows.Scan(&userID, &weight, &lastEntryID)
		if err != nil {
			log.Printf("Warning: could not scan variables from row")
			continue
		}

		if lastEntryID == "" {
			entries = append(entries, GiveawayEntry{UserID: userID, Weight: weight, Platform: platform, PlatformID: platformID})
			continue
		}

		dateValue, ok := strings.CutPrefix(lastEntryID, prefix+"-")
		if !ok {
			continue
		}

		lastEntry, err := time.Parse(time.DateOnly, dateValue)
		if err != nil {
			log.Printf("ERROR: could not convert last_entry_id '%s' to time: %v", lastEntryID, err)
			continue
		}
		entries = append(entries, GiveawayEntry{userID, weight, lastEntry, platform, platformID})
	}
	return entries
}

// DrawGiveawayWinner takes one of the given entries and draw one winner of them. The probability
// is based on their Weight value. A higher Weight means a higher probability.
func DrawGiveawayWinner(e []GiveawayEntry) (winner GiveawayEntry, totalTickets int) {
	// skip randomizing when there is only one entry
	if len(e) == 1 {
		return e[0], e[0].Weight
	}

	var entries []GiveawayEntry
	for _, e := range e {
		for i := 0; i < e.Weight; i++ {
			entries = append(entries, e)
		}
	}
	totalTickets = len(entries)
	if totalTickets == 0 {
		return GiveawayEntry{}, 0
	}

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})
	return entries[rand.Intn(totalTickets-1)], totalTickets
}

// GiveawayPrizeType represents the type of a giveaway prize
type GiveawayPrizeType string

// GiveawayPrizeGroupSort represents the sorting type of a group of prizes
type GiveawayPrizeGroupSort string

const (
	// GiveawayPrizeTypeSingle describes a single giveaway prize
	GiveawayPrizeTypeSingle GiveawayPrizeType = "single"
	// GiveawayPrizeTypeGroup describes a group containing a pool of giveaway prizes
	GiveawayPrizeTypeGroup GiveawayPrizeType = "group"

	// GiveawayPrizeGroupOrdered defines that the pool in this group is ordered and prizes should be
	// drawn in ascending order
	GiveawayPrizeGroupOrdered GiveawayPrizeGroupSort = "ordered"
	// GiveawayPrizeGroupRandom defines that the pool in this group contains a set of prizes that
	// should be drawn in a random order
	GiveawayPrizeGroupRandom GiveawayPrizeGroupSort = "random"
)

// giveawayPrizeInterface is a helper type to un-/marshal a giveaway prize json file
type giveawayPrizeInterface interface {
	prizeType() GiveawayPrizeType
}

// GiveawayPrize represents a general giveaway prize. You can unmarshal to it as well as marshal it
// again. Various functions provide access to modify the pool of prizes.
type GiveawayPrize struct {
	giveawayPrizeInterface
	filename string
}

// NewGiveawayPrize reads the file and stores it in a new GiveawayPrize struct.
//
// When modifying something, make sure to call p.SaveFile() to save the changes back to the file.
func NewGiveawayPrize(filename string) (p GiveawayPrize, err error) {
	if filename == "" {
		return p, fmt.Errorf("argument filename cannot be empty")
	}
	p.filename = filename

	err = p.ReadFile()
	return p, err

}

// ReadFile reads the giveaway file from the configured filename and stores it in p
func (p *GiveawayPrize) ReadFile() error {
	if p == nil || p.filename == "" {
		return fmt.Errorf("cannot read to invalid GiveawayPrize! Make sure to use NewGiveawayPrize()")
	}
	data, err := os.ReadFile(p.filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &p)
}

// SaveFile saves p in the configured json file
func (p GiveawayPrize) SaveFile() error {
	if p.filename == "" {
		return fmt.Errorf("cannot save invalid GiveawayPrize! Make sure to use NewGiveawayPrize()")
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p.filename, data, 0644)
}

// HasPrizeAvailable returns whether p has at least one prize without a winner
func (p GiveawayPrize) HasPrizeAvailable() bool {
	return p.HasPrizeWon("")
}

// HasPrizeWon returns whether p has at least one prize where the winner is the same as userID
func (p GiveawayPrize) HasPrizeWon(userID string) bool {
	if p.giveawayPrizeInterface == nil {
		return false
	}

	switch t := p.giveawayPrizeInterface.(type) {
	case *GiveawayPrizeSingle:
		return t.Winner == userID
	case *GiveawayPrizeGroup:
		return t.HasPrizeWon(userID)
	}
	return false
}

// GetNextPrize returns a pointer to the next giveaway prize. You can make any changes to it but in
// in order to save it back to the file, use p.SaveFile().
//
// The next prize is determined by the next one which has an empty Winner field. If the prize
// whould be a group which sort is set to random, then one of the prizes in its pool, which has an
// empty Winner filed (or if it is also a group, then when it contains at least one prize without
// a winner), is selected.
//
// When there is no prize available ok will be false.
func (p *GiveawayPrize) GetNextPrize() (*GiveawayPrizeSingle, bool) {
	if p == nil || p.giveawayPrizeInterface == nil {
		return nil, false
	}

	switch t := p.giveawayPrizeInterface.(type) {
	case *GiveawayPrizeSingle:
		if t.Winner == "" {
			return t, true
		}
	case *GiveawayPrizeGroup:
		return t.GetNextPrize()
	}
	return nil, false
}

// UnmarshalJSON implements json.Unmarshaler
func (p *GiveawayPrize) UnmarshalJSON(data []byte) error {
	var h struct {
		Type GiveawayPrizeType `json:"type"`
	}
	err := json.Unmarshal(data, &h)
	if err != nil {
		return err
	}

	switch h.Type {
	case GiveawayPrizeTypeSingle:
		var t GiveawayPrizeSingle
		err = json.Unmarshal(data, &t)
		p.giveawayPrizeInterface = &t
		return err
	case GiveawayPrizeTypeGroup:
		var t GiveawayPrizeGroup
		err = json.Unmarshal(data, &t)
		p.giveawayPrizeInterface = &t
		return err
	default:
		return &json.UnmarshalTypeError{
			Value: string(h.Type),
			Type:  reflect.TypeOf(h.Type),
		}
	}
}

// MarshalJSON implements json.Marshaler
func (p GiveawayPrize) MarshalJSON() ([]byte, error) {
	if p.giveawayPrizeInterface == nil {
		return []byte{}, &json.MarshalerError{
			Type: reflect.TypeOf(p),
			Err:  fmt.Errorf("underlying prize is nil"),
		}
	}
	b := bytes.NewBuffer([]byte{})
	b.WriteByte('{')
	const format string = "\"%s\":\"%s\""
	b.WriteString(fmt.Sprintf(format, "type", p.prizeType()))
	buf, err := json.Marshal(p.giveawayPrizeInterface)
	if err != nil {
		return []byte{}, err
	}
	b.WriteByte(',')
	b.Write(buf[1:])
	return b.Bytes(), nil
}

// GiveawayPrizeSingle represents a single giveaway prize. Its the lowest struct from all giveaway
// structures.
type GiveawayPrizeSingle struct {
	// The name of prize
	Name string `json:"name"`

	// The identifier of the winner. An empty string means this prize has no winner yet and is
	// available.
	Winner string `json:"winner,omitempty"`
}

func (p GiveawayPrizeSingle) prizeType() GiveawayPrizeType {
	return GiveawayPrizeTypeSingle
}

// GiveawayPrizeGroup represents a pool of prizes. The behavior or the order is defined by the Sort
// field.
type GiveawayPrizeGroup struct {
	// The order of the pool. Defines in which order to read the prizes in the pool
	Sort GiveawayPrizeGroupSort `json:"sort"`
	// All prizes that belong to this group
	Pool []GiveawayPrize `json:"pool"`
}

func (pg GiveawayPrizeGroup) prizeType() GiveawayPrizeType {
	return GiveawayPrizeTypeGroup
}

// HasPrizeAvailable returns whether pg contains at least one prize without a winner
func (pg GiveawayPrizeGroup) HasPrizeAvailable() bool {
	return pg.HasPrizeWon("")
}

// HasPrizeWon returns whether pg contains at least one prize where the winner is the same as userID
func (pg GiveawayPrizeGroup) HasPrizeWon(userID string) bool {
	for _, p := range pg.Pool {
		if p.HasPrizeWon(userID) {
			return true
		}
	}
	return false
}

// GetNextPrize returns a pointer to the next giveaway prize. You can make any changes to it but in
// in order to save it back to the file, use p.SaveFile().
//
// The next prize is determined by the next one which has an empty Winner field. If the prize
// whould be a group which sort is set to random, then one of the prizes in its pool, which has an
// empty Winner filed (or if it is also a group, then when it contains at least one prize without
// a winner), is selected.
//
// When there is no prize available ok will be false.
func (pg *GiveawayPrizeGroup) GetNextPrize() (*GiveawayPrizeSingle, bool) {
	if pg == nil || len(pg.Pool) == 0 {
		return nil, false
	}

	switch pg.Sort {
	case GiveawayPrizeGroupOrdered:
		for i, p := range pg.Pool {
			if s, ok := p.GetNextPrize(); ok {
				pg.Pool[i] = p
				return s, true
			}
		}
	case GiveawayPrizeGroupRandom:
		var available []int
		for i, p := range pg.Pool {
			if p.HasPrizeAvailable() {
				available = append(available, i)
			}
		}

		switch len(available) {
		case 0:
		case 1:
			return pg.Pool[available[0]].GetNextPrize()
		default:
			rand.Shuffle(len(available), func(i, j int) {
				available[i], available[j] = available[j], available[i]
			})
			i := available[rand.Intn(len(available)-1)]
			return pg.Pool[i].GetNextPrize()
		}
	}
	return nil, false
}
