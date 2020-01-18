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
	"github.com/getsolus/ferryd/jobs"
	"os"
)

// Remove fulfills the "remove-repo" sub-command
var Remove = &cmd.CMD{
	Name:  "remove-repo",
	Alias: "rr",
	Short: "Remove a repo from the DB, but leave on disk",
	Args:  &RemoveArgs{},
	Run:   RemoveRun,
}

// RemoveArgs are the arguments to the "remove-repo" sub-command
type RemoveArgs struct {
	Repo string `desc:"Repo to remove"`
}

// RemoveRun executes the "remove-repo" sub-command
func RemoveRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*RemoveArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	var j *jobs.Job
	var err error
	if j, err = client.Remove(args.Repo); err != nil {
		fmt.Fprintf(os.Stderr, "Error while removing repo: %v\n", err)
		os.Exit(1)
	}
	// Print a summary of the job
	j.Print()
}
