//
// Copyright © 2017-2020 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package jobs

import (
	"database/sql"
	"time"
)

// NullString extends sql.NullString with JSON-compatible Marshall/Unmarshall
type NullString sql.NullString

// MarshalText converts a NullString to "null" or its valid value
func (ns NullString) MarshalText() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return []byte(ns.String), nil
}

// UnmarshalText parses a NullString from text and determines if it is a valid value
func (ns *NullString) UnmarshalText(text []byte) error {
	ns.Valid = false
	ns.String = ""
	v := string(text)
	if v == "null" {
		return nil
	}
	ns.Valid = true
	ns.String = v
	return nil
}

// NullTime extends sql.NullTime with JSON-compatible Marshall/Unmarshall
type NullTime sql.NullTime

// MarshalText converts a NullTime to "null" or its valid value
func (nt NullTime) MarshalText() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return []byte(nt.Time.UTC().Format("2006-01-02T15:04:05Z")), nil
}

// UnmarshalText parses a NullTime from text and determines if it is a valid value
func (nt *NullTime) UnmarshalText(text []byte) error {
	nt.Valid = false
	nt.Time = time.Time{}
	v := string(text)
	if v == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", v)
	if err != nil {
		return err
	}
	nt.Valid = true
	nt.Time = t
	return nil
}
