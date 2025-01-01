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

package util

import (
	"github.com/cake4everyone/cake4everybot/logger"
)

var log = logger.New("Util")

var (
	constIntZero   int64   = 0
	constIntOne            = constIntZero + 1
	constIntTwo            = constIntOne + 1
	constFloatZero float64 = 0
	constFloatOne          = constFloatZero + 1
	constFloatTwo          = constFloatOne + 1
)

// ContainsInt reports whether at least one of num is at least once anywhere in i.
func ContainsInt(i []int, num ...int) bool {
	for _, x := range i {
		for _, y := range num {
			if x == y {
				return true
			}
		}
	}
	return false
}

// ContainsString reports whether at least one of str is at least once anywhere in s.
func ContainsString(s []string, str ...string) bool {
	for _, x := range s {
		for _, y := range str {
			if x == y {
				return true
			}
		}
	}
	return false
}

// Btoi returns the integer for the given boolean b.
//
//	Btoi(false) => 0
//	Btoi(true) => 1
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ShiftL takes a slice and shifts all elements to the left. The first element pops out and is
// returned. If s is an empty slice the zero value of the given type is returned. If t is given it
// will be inserted at the last position instead of an element with its zero value.
func ShiftL[T any](s []T, t ...T) (first T) {
	for i, v := range s {
		if i == 0 {
			first = v
			continue
		}
		s[i-1] = s[i]
		if i == len(s)-1 {
			var last T
			if len(t) > 0 {
				last = t[0]
			}
			s[i] = last
		}
	}
	return first
}

// IntZero returns a pointer to an [int64] with the value 0.
func IntZero() *int64 {
	return &constIntZero
}

// IntOne returns a pointer to an [int64] with the value 1.
func IntOne() *int64 {
	return &constIntOne
}

// IntTwo returns a pointer to an [int64] with the value 2.
func IntTwo() *int64 {
	return &constIntTwo
}

// FloatZero returns a pointer to a [float64] with the value 0.
func FloatZero() *float64 {
	return &constFloatZero
}

// FloatOne returns a pointer to a [float64] with the value 1.
func FloatOne() *float64 {
	return &constFloatOne
}

// FloatTwo returns a pointer to a [float64] with the value 2.
func FloatTwo() *float64 {
	return &constFloatTwo
}
