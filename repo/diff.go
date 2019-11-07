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

// Diff is a list of changes made to a repo
type Diff []string

// Print writes out a Diff in a human-readable format
func (d *Diff) Print(out io.Writer, full, color bool) {
    plusFmt := "%s\n"
    minusFmt := "%s\n"
    if color {
        plusFmt  = "\033[49;38;5;208%s\033[0m\n"
        minusFmt = "\033[49;38;5;040%s\033[0m\n"
    }
    for _, change := range d {
        switch r := []rune(change)[0]
        case '+':
            fmt.Fprintf(out, plusFmt, change)
        case '-':
            fmt.Fprintf(out, minusFmt, change)
        default:
            if full {
                fmt.Fprintln(out, change)
            }
        }
    }
}
