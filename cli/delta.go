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

package cli

import (
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
	"github.com/getsolus/ferryd/api/v1"
	"github.com/getsolus/ferryd/repo"
	"os"
)

// Delta fulfills the "delta" sub-command
var Delta = &cmd.CMD{
	Name:  "delta",
	Alias: "dr",
	Short: "Generate missing deltas and cleanup outdated ones",
	Args:  &DeltaArgs{},
	Run:   DeltaRun,
}

// DeltaArgs are the arguments to the "delta" sub-command
type DeltaArgs struct {
	Repo string `desc:"Repo for updating deltas"`
}

// DeltaRun executes the "delta" sub-command
func DeltaRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*DeltaArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	j, err := client.Delta(args.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while updating repo deltas: %v\n", err)
		os.Exit(1)
	}
	// Print the job Summary
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
