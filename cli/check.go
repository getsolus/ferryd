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

// Check fulfills the "check" sub-command
var Check = &cmd.CMD{
	Name:  "check",
	Alias: "ch",
	Short: "Compare the disk contents of a repo with the DB",
	Args:  &CheckArgs{},
	Run:   CheckRun,
}

// CheckArgs are the arguments to the "check" sub-command
type CheckArgs struct {
	Repo string `desc:"Repo to compare against"`
}

// CheckRun executes the "check" sub-command
func CheckRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*CheckArgs)
	// Create a client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	d, j, err := client.Check(args.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while checking repo: %v\n", err)
		os.Exit(1)
	}
	// Print the job summary
	j.Print()
	// Print the diff
	d.Print(os.Stdout, false, !flags.NoColor)
}
