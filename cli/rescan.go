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
	"github.com/getsolus/ferryd/api/v1"
	"github.com/getsolus/ferryd/repo"
	"os"
)

// Rescan fulfills the "rescan" sub-command
var Rescan = &cmd.CMD{
	Name:  "rescan",
	Alias: "rs",
	Short: "Rescan a repo on disk and make the DB match",
	Args:  &RescanArgs{},
	Run:   RescanRun,
}

// RescanArgs are the arguments to the "rescan" sub-command
type RescanArgs struct {
	Repo string `desc:"Repo to rescan"`
}

// RescanRun executes the "rescan" sub-command
func RescanRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*RescanArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	j, err := client.Rescan(args.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while rescanning repo: %v\n", err)
		os.Exit(1)
	}
	// print the job summary
	j.Print()
	// Decode the Diff
	var d *repo.Diff
	if err = d.UnmarshalBinary(j.Results); err != nil {
		fmt.Fprintf(os.Stderr, "Error while decoding diff: %v\n", err)
		os.Exit(1)
	}
	// Print the diff
	d.Print(os.Stdout, false, !flags.NoColor)
}
