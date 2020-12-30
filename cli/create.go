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

// Create fulfills the "create-repo" sub-command
var Create = &cmd.CMD{
	Name:  "create-repo",
	Alias: "cr",
	Short: "Create a new empty repository",
	Args:  &CreateArgs{},
	Run:   CreateRun,
}

// CreateArgs are the arguments to the "create-repo" sub-command
type CreateArgs struct {
	Repo    string `desc:"Name of the repo to create"`
	Instant bool   `desc:"Decide whether or not a repo should have instant transit"`
}

// CreateRun executes the "create-repo" sub-command
func CreateRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*CreateArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	j, err := client.Create(args.Repo, args.Instant)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating repo: %v\n", err)
		os.Exit(1)
	}
	// Print a summary of the job
	j.Print()
}
