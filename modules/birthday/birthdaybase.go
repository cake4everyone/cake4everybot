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

package birthday

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cake4everyone/cake4everybot/data/lang"
	"github.com/cake4everyone/cake4everybot/database"
	"github.com/cake4everyone/cake4everybot/logger"
	"github.com/cake4everyone/cake4everybot/util"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => birthday
	tp = "discord.command.birthday."
)

var log = logger.New("Birthday")

type birthdayBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

type birthdayEntry struct {
	ID          uint64 `database:"id"`
	Day         int    `database:"day"`
	Month       int    `database:"month"`
	Year        int    `database:"year"`
	Visible     bool   `database:"visible"`
	time        time.Time
	GuildIDsRaw string `database:"guilds"`
	GuildIDs    []string
}

// Returns a readable Form of the date
func (b birthdayEntry) String() string {
	if b.Year == 0 {
		month := lang.GetSliceElement(tp+"month", b.Month-1, lang.FallbackLang())
		return fmt.Sprintf("%d. %s", b.Day, month)
	}
	return fmt.Sprintf("<t:%d:D>", b.time.Unix())
}

// DOW returns the day of the week with
//
//	MON = 0
//	TUE = 1
//	SUN = 6
//
// C code by sakamoto@sm.sony.co.jp (Tomohiko Sakamoto) posted on
// comp.lang.c on March 10th, 1993:
//
//	dow(m,d,y){y-=m<3;return(y+y/4-y/100+y/400+"-bed=pen+mad."[m]+d)%7;}
func (b birthdayEntry) DOW() int {
	b.Year = b.Year - util.Btoi(b.Month < 3)
	monthKey := "WLAN or Fence"
	return int(b.Year+b.Year/4-b.Year/100+b.Year/400+int(monthKey[b.Month])+b.Day) % 7
}

// NextUnix returns the unix timestamp in seconds of the next birthday. This is just a shorthand for
// Next().Unix().
//
// See Next() for more.
func (b birthdayEntry) NextUnix() int64 {
	return b.Next().Unix()
}

// Next returns the time.Time object of the next birthday.
func (b birthdayEntry) Next() time.Time {
	years := time.Now().Year() - b.Year
	nextTime := b.time.AddDate(years, 0, 0)
	if time.Until(nextTime) <= 0 {
		nextTime = b.time.AddDate(years+1, 0, 0)
	}
	return nextTime
}

// ParseTime tries to parse the date (b.Day, b.Month, b.Year) to a time.Time object.
func (b *birthdayEntry) ParseTime() (err error) {
	b.time, err = time.Parse(time.DateOnly, fmt.Sprintf("%04d-%02d-%02d", b.Year, b.Month, b.Day))
	return err
}

// Age returns the current age of the user. If no year is set, it returns 0.
func (b birthdayEntry) Age() int {
	if b.Year == 0 {
		return 0
	}
	return b.Next().Year() - b.Year - 1
}

// ParseGuildIDs splits the guild IDs into a slice and stores them in b.GuildIDs.
func (b *birthdayEntry) ParseGuildIDs() {
	b.GuildIDs = strings.Split(b.GuildIDsRaw, ",")
}

// IsInGuild returns true if the guildID is in b.GuildIDs.
// If guildID is empty, IsInGuild returns true.
func (b birthdayEntry) IsInGuild(guildID string) bool {
	if guildID == "" {
		return true
	}
	return util.ContainsString(b.GuildIDs, guildID)
}

// SetGuild sets the guildID in the birthday entry.
func (b *birthdayEntry) SetGuild(guildID string) {
	b.GuildIDsRaw += guildID
	b.ParseGuildIDs()
}

// AddGuild adds the guildID to the birthday entry.
func (b *birthdayEntry) AddGuild(guildID string) error {
	if util.ContainsString(b.GuildIDs, guildID) {
		return nil
	} else if len(b.GuildIDs) >= 3 {
		return fmt.Errorf("this entry already has %d guilds", len(b.GuildIDs))
	}
	b.GuildIDsRaw += "," + guildID
	b.GuildIDsRaw = strings.Trim(b.GuildIDsRaw, ", ")
	b.ParseGuildIDs()
	return nil
}

// IsEqual returns true if b and b2 are equal.
//
// That is, if all of the following are true
//  1. They have the same user ID.
//  2. They are on the same date.
//  3. They have the same visibility.
//  4. They have the same guilds in (any order).
func (b birthdayEntry) IsEqual(b2 birthdayEntry) bool {
	if b.ID != b2.ID || b.Day != b2.Day || b.Month != b2.Month || b.Year != b2.Year || b.Visible != b2.Visible {
		return false
	}

	// check for same guilds in any order
	for _, guildID := range b.GuildIDs {
		if !util.ContainsString(b2.GuildIDs, guildID) {
			return false
		}
	}
	for _, guildID := range b2.GuildIDs {
		if !util.ContainsString(b.GuildIDs, guildID) {
			return false
		}
	}
	return true
}

// getBirthday copies all birthday fields into the struct pointed at by b.
//
// If the user from b.ID is not found it returns sql.ErrNoRows.
func (cmd birthdayBase) getBirthday(b *birthdayEntry) (err error) {
	row := database.QueryRow("SELECT day,month,year,visible,guilds FROM birthdays WHERE id=?", b.ID)
	err = row.Scan(&b.Day, &b.Month, &b.Year, &b.Visible, &b.GuildIDsRaw)
	if err != nil {
		return err
	}
	b.ParseGuildIDs()
	return b.ParseTime()
}

// hasBirthday returns true whether the given user id has entered their birthday.
func (cmd birthdayBase) hasBirthday(id uint64) (hasBirthday bool, err error) {
	err = database.QueryRow("SELECT EXISTS(SELECT id FROM birthdays WHERE id=?)", id).Scan(&hasBirthday)
	return hasBirthday, err
}

// setBirthday inserts a new database entry with the values from b.
func (cmd birthdayBase) setBirthday(b *birthdayEntry) (err error) {
	b.SetGuild(cmd.Interaction.GuildID)
	_, err = database.Exec("INSERT INTO birthdays(id,day,month,year,visible,guilds) VALUES(?,?,?,?,?);", b.ID, b.Day, b.Month, b.Year, b.Visible, b.GuildIDsRaw)
	return err
}

// updateBirthday updates an existing database entry with the values from b.
func (cmd birthdayBase) updateBirthday(b *birthdayEntry) (before birthdayEntry, err error) {
	before.ID = b.ID
	if err = cmd.getBirthday(&before); err != nil {
		return birthdayEntry{}, fmt.Errorf("trying to get old birthday: %v", err)
	}
	b.GuildIDsRaw = before.GuildIDsRaw
	b.ParseGuildIDs()

	err = b.AddGuild(cmd.Interaction.GuildID)
	if err != nil {
		return birthdayEntry{}, fmt.Errorf("adding guild '%s' to birthday entry: %v", cmd.Interaction.GuildID, err)
	}

	// early return if nothing changed
	if b.IsEqual(before) {
		return before, nil
	}

	var (
		updateNames []string
		updateVars  []any
		oldV        reflect.Value = reflect.ValueOf(before)
		v           reflect.Value = reflect.ValueOf(*b)
	)
	for i := 0; i < v.NumField(); i++ {
		var (
			oldF = oldV.Field(i)
			f    = v.Field(i)
			p    = f.Type().PkgPath()
		)
		//skip unexported fields
		if len(p) != 0 {
			continue
		}

		tag := v.Type().Field(i).Tag.Get("database")
		if tag == "" {
			continue
		}
		if f.Interface() != oldF.Interface() {
			updateNames = append(updateNames, tag)
			updateVars = append(updateVars, f.Interface())
		}

	}

	if len(updateNames) == 0 {
		return before, nil
	}

	updateString := strings.Join(updateNames, "=?,") + "=?"
	_, err = database.Exec("UPDATE birthdays SET "+updateString+" WHERE id="+fmt.Sprint(b.ID)+";", updateVars...)
	return before, err
}

// removeBirthday deletes the existing birthday entry for the given id and returns the previously
// entered birthday.
func (cmd birthdayBase) removeBirthday(id uint64) (birthdayEntry, error) {
	b := birthdayEntry{ID: id}
	err := cmd.getBirthday(&b)
	if err != nil {
		return b, err
	}

	_, err = database.Exec("DELETE FROM birthdays WHERE id=?;", b.ID)
	return b, err
}

// getBirthdaysMonth return a sorted slice of birthday entries that matches the given guildID and month.
func (cmd birthdayBase) getBirthdaysMonth(guildID string, month int) (birthdays []birthdayEntry, err error) {
	var numOfEntries int64
	err = database.QueryRow("SELECT COUNT(*) FROM birthdays WHERE month=?", month).Scan(&numOfEntries)
	if err != nil {
		return nil, err
	}

	birthdays = make([]birthdayEntry, 0, numOfEntries)
	if numOfEntries == 0 {
		return birthdays, nil
	}

	rows, err := database.Query("SELECT id,day,year,visible,guilds FROM birthdays WHERE month=?", month)
	if err != nil {
		return birthdays, err
	}
	defer rows.Close()

	for rows.Next() {
		b := birthdayEntry{Month: month}
		err = rows.Scan(&b.ID, &b.Day, &b.Year, &b.Visible, &b.GuildIDsRaw)
		if err != nil {
			return birthdays, err
		}
		b.ParseGuildIDs()

		if !b.Visible || !b.IsInGuild(guildID) {
			continue
		}

		err = b.ParseTime()
		if err != nil {
			return birthdays, err
		}

		birthdays = append(birthdays, b)
	}

	sort.Slice(birthdays, func(i, j int) bool {
		return birthdays[i].Day < birthdays[j].Day
	})

	return birthdays, nil
}

// getBirthdaysDate return a slice of birthday entries that matches the given guildID and date.
func getBirthdaysDate(guildID string, day int, month int) (birthdays []birthdayEntry, err error) {
	var numOfEntries int64
	err = database.QueryRow("SELECT COUNT(*) FROM birthdays WHERE day=? AND month=?", day, month).Scan(&numOfEntries)
	if err != nil {
		return nil, err
	}

	birthdays = make([]birthdayEntry, 0, numOfEntries)
	if numOfEntries == 0 {
		return birthdays, nil
	}

	rows, err := database.Query("SELECT id,year,visible,guilds FROM birthdays WHERE day=? AND month=?", day, month)
	if err != nil {
		return birthdays, err
	}
	defer rows.Close()

	for rows.Next() {
		b := birthdayEntry{Day: day, Month: month}
		err = rows.Scan(&b.ID, &b.Year, &b.Visible, &b.GuildIDsRaw)
		if err != nil {
			return birthdays, err
		}
		b.ParseGuildIDs()

		if !b.Visible || !b.IsInGuild(guildID) {
			continue
		}

		err = b.ParseTime()
		if err != nil {
			return birthdays, err
		}

		birthdays = append(birthdays, b)
	}

	return birthdays, nil
}
