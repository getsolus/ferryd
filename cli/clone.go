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
	"os"
)

// Clone fulfills the "clone" sub-command
var Clone = &cmd.CMD{
	Name:  "clone",
	Alias: "cl",
	Short: "Clone an existing repository into a new repository",
	Args:  &CloneArgs{},
	Run:   CloneRun,
}

// CloneArgs are the arguments to the "clone" sub-command
type CloneArgs struct {
	Source string `desc:"Repo to clone from"`
	Dest   string `desc:"New Repo to create and clone into"`
}

// CloneRun executes the "clone" sub-command
func CloneRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*CloneArgs)
	// Create a client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	s, j, err := client.Clone(args.Source, args.Dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while cloning repo: %v\n", err)
		os.Exit(1)
	}
	// Print the job summary
	j.Print()
	// Print repo summary
	s.Print(os.Stdout, true)
}
