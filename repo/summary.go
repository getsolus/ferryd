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
	"io"
)

// Summary is a brief description of a single Repo
type Summary struct {
	Name     string
	Packages uint64
	Deltas   uint64
	Size     uint64
}

// Print writes out a Summary in a human-readable format
func (s *Summary) Print(out io.Writer, single bool) {
	if s == nil {
		fmt.Fprintln(out, "No summary found.")
		return
	}
	if single {
		fmt.Fprintf(out, "Name: %s\n", s.Name)
		fmt.Fprintf(out, "\tPackages: %d\n", s.Packages)
		fmt.Fprintf(out, "\t  Deltas: %d\n", s.Deltas)
		fmt.Fprintf(out, "\t    Size: %d\n", s.Size)
		fmt.Fprintln(out)
	} else {
		fmt.Fprintf(out, "\tName: %s\n", s.Name)
		fmt.Fprintf(out, "\t\tPackages: %d\n", s.Packages)
		fmt.Fprintf(out, "\t\t  Deltas: %d\n", s.Deltas)
		fmt.Fprintf(out, "\t\t    Size: %d\n", s.Size)
		fmt.Fprintln(out)
	}
}

// FullSummary is a brief description of all Repos
type FullSummary []Summary

// Print writes out a FullSummary in a human-readable format
func (f FullSummary) Print(out io.Writer) {
	for _, s := range f {
		s.Print(out, false)
	}
}
