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

// Import fulfills the "import" sub-command
var Import = &cmd.CMD{
	Name:  "import",
	Alias: "imp",
	Short: "Import a repo on disk into the DB",
	Args:  &ImportArgs{},
	Run:   ImportRun,
}

// ImportArgs are the arguments to the "import" sub-command
type ImportArgs struct {
	Repo    string `desc:"Repo to import"`
	Instant bool   `desc:"Indicate whether this repo should be instant transit or not"`
}

// ImportRun executes the "import" sub-command
func ImportRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*ImportArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	s, j, err := client.Import(args.Repo, args.Instant)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while importing repo: %v\n", err)
		os.Exit(1)
	}
	// Print a summary of the job
	j.Print()
	// Print the repo summary
	s.Print(os.Stdout, true)
}
