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

package repo

import (
	"bytes"
	"encoding/gob"
	"github.com/getsolus/ferryd/repo/archive"
	"io"
)

// Diff is a list of changes made to a repo
type Diff archive.Archives

// MarshalBinary converts a Diff to its Gob encoded form
func (d *Diff) MarshalBinary() (data []byte, err error) {
	buff := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buff)
	if err = enc.Encode(d); err == nil {
		data = buff.Bytes()
	}
	return
}

// UnmarshalBinary converts a Gob encoded Diff back to its useful form
func (d *Diff) UnmarshalBinary(data []byte) error {
	buff := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buff)
	return dec.Decode(d)
}

// Print writes out a Diff in a human-readable format
func (d Diff) Print(out io.Writer, full, color bool) {
	plus := "+%s\n"
	minus := "-%s\n"
	mod := "!%s\n"
	same := " %s\n"
	// Override the format strings if printing with color
	if color {
		plus = "\033[49;38;5;040m+%s\033[0m\n"
		minus = "\033[49;38;5;208m-%s\033[0m\n"
		mod = "\033[49;38;5;220m!%s\033[0m\n"
		same = "\033[49;39m %s\033[0m\n"
	}
	// Print each line
	for _, a := range d {
		a.PrintDiff(out, plus, minus, mod, same)
	}
}
