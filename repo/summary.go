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
	"database/sql"
	"fmt"
	"io"
)

// Summary is a brief description of a single Repo
type Summary struct {
	Name        string        `json:"name"`
	Packages    sql.NullInt64 `json:"packages"`
	Deltas      sql.NullInt64 `json:"deltas"`
	ArchiveSize sql.NullInt64 `json:"size"`
	Used        uint64        `json:"used"`
	Free        uint64        `json:"free"`
}

// Print writes out a Summary in a human-readable format
func (s *Summary) Print(out io.Writer, single bool) {
	// Don't try to print a null summary
	if s == nil {
		fmt.Fprintln(out, "No summary found.")
		return
	}
	if single {
		// No indent
		fmt.Fprintf(out, "Name: %s\n", s.Name)
		if s.Packages.Valid {
			fmt.Fprintf(out, "\tPackages: %d\n", s.Packages.Int64)
		}
		if s.Deltas.Valid {
			fmt.Fprintf(out, "\t  Deltas: %d\n", s.Deltas.Int64)
		}
		if s.ArchiveSize.Valid {
			fmt.Fprintf(out, "\t    Size: %d\n", s.ArchiveSize.Int64)
		}
		fmt.Fprintln(out, "\tFilesystem:")
		fmt.Fprintf(out, "\t      Used: %d\n", s.Used)
		fmt.Fprintf(out, "\t      Free: %d\n", s.Free)
		fmt.Fprintln(out)
	} else {
		// One Indent
		fmt.Fprintf(out, "\tName: %s\n", s.Name)
		if s.Packages.Valid {
			fmt.Fprintf(out, "\t\tPackages: %d\n", s.Packages.Int64)
		}
		if s.Deltas.Valid {
			fmt.Fprintf(out, "\t\t  Deltas: %d\n", s.Deltas.Int64)
		}
		if s.ArchiveSize.Valid {
			fmt.Fprintf(out, "\t\t    Size: %d\n", s.ArchiveSize.Int64)
		}
		fmt.Fprintln(out, "\t\tFilesystem:")
		fmt.Fprintf(out, "\t\t      Used: %d\n", s.Used)
		fmt.Fprintf(out, "\t\t      Free: %d\n", s.Free)
		fmt.Fprintln(out)
	}
}

// FullSummary is a brief description of all Repos
type FullSummary []Summary

// Print writes out a FullSummary in a human-readable format
func (f FullSummary) Print(out io.Writer) {
	fmt.Println("Repositories:")
	for _, s := range f {
		s.Print(out, false)
	}
}
