// Package status keeps track of the bots internal status.
//
// # Copyright 2022-2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package status

import "time"

// The point in time where the bot starts
var startTime time.Time

func init() {
	startTime = time.Now()
}

// GetStartTime returns the time were the bot stared
func GetStartTime() time.Time {
	return startTime
}

// GetUptime returns the duration since the bot stared
func GetUptime() time.Duration {
	return time.Since(startTime)
}
