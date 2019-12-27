//
// Copyright © 2017-2019 Solus Project
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
	"fmt"
	"github.com/getsolus/ferryd/repo/releases"
	"io"
	"sort"
)

// Diff is a list of changes made to a repo
type Diff struct {
	Package string
	keys    ReleasePairs
	Lines   map[ReleasePair]string
}

// ReleasePair is used for sorting releases
type ReleasePair struct {
	From int
	To   int
}

// ReleasePairs is a sortable list of ReleasePair
type ReleasePairs []ReleasePair

// Len returns the length of a ReleasePairs for sorting
func (rps ReleasePairs) Len() int {
	return len(rps)
}

// Less compares two ReleasePair instances for sorting
func (rps ReleasePairs) Less(i, j int) bool {
	return (rps[i].To < rps[j].To) || ((rps[i].To == rps[j].To) && (rps[i].From < rps[j].From))
}

// Swap carries out swapping for sorting
func (rps ReleasePairs) Swap(i, j int) {
	rps[i], rps[j] = rps[j], rps[i]
}

// NewDiff creates a Diff from the results of a Compare operation
func NewDiff(l, r, s []releases.Release) Diff {
	d := Diff{
		Lines: make(map[ReleasePair]string),
	}
	for _, e := range l {
		rp := ReleasePair{
			To:   e.Release,
			From: e.From,
		}
		d.keys = append(d.keys, rp)
		d.Lines[rp] = fmt.Sprintf("+ %s", e.URI)
	}
	for _, e := range r {
		rp := ReleasePair{
			To:   e.Release,
			From: e.From,
		}
		d.keys = append(d.keys, rp)
		d.Lines[rp] = fmt.Sprintf("- %s", e.URI)
	}
	for _, e := range s {
		rp := ReleasePair{
			To:   e.Release,
			From: e.From,
		}
		d.keys = append(d.keys, rp)
		d.Lines[rp] = fmt.Sprintf("  %s", e.URI)
	}
	sort.Sort(d.keys)
	return d
}

// Print writes out a Diff in a human-readable format
func (d Diff) Print(out io.Writer, full, color bool) {
	plusFmt := "%s\n"
	minusFmt := "%s\n"
	if color {
		plusFmt = "\033[49;38;5;208%s\033[0m\n"
		minusFmt = "\033[49;38;5;040%s\033[0m\n"
	}
	for _, k := range d.keys {
		line := d.Lines[k]
		rs := []rune(line)
		switch rs[0] {
		case '+':
			fmt.Fprintf(out, plusFmt, line)
		case '-':
			fmt.Fprintf(out, minusFmt, line)
		default:
			if full {
				fmt.Fprintln(out, line)
			}
		}
	}
}
