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

package cli

import (
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/getsolus/ferryd/api"
	"os"
)

// TrimObsoletes fulfills the "trim-obsoletes" sub-command
var TrimObsoletes = &cmd.CMD{
	Name:  "trim-obsoletes",
	Alias: "to",
	Short: "Remove all obsolete packages from a repo",
	Args:  &TrimObsoletesArgs{},
	Run:   TrimObsoletesRun,
}

// TrimObsoletesArgs are the arguments to the "trim-obsoletes" sub-command
type TrimObsoletesArgs struct {
	Repo string `desc:"Repo to trim"`
}

// TrimObsoletesRun executes the "trim-obsoletes" sub-command
func TrimObsoletesRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*TrimObsoletesArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	d, j, err := client.TrimObsoletes(args.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while trimming obsolete packages in repo: %v\n", err)
		os.Exit(1)
	}
	// Print the job summary
	j.Print()
	// Print the diff
	d.Print(os.Stdout, false, !flags.NoColor)
}
