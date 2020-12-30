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

// Index fulfills the "index" sub-command
var Index = &cmd.CMD{
	Name:  "index",
	Alias: "idx",
	Short: "Update the Index for a repo",
	Args:  &IndexArgs{},
	Run:   IndexRun,
}

// IndexArgs are the arguments to the "index" sub-command
type IndexArgs struct {
	Repo string `desc:"Repo to generate a new index for"`
}

// IndexRun executes the "index" sub-command
func IndexRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*IndexArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	j, err := client.Index(args.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while generating repo Index: %v\n", err)
		os.Exit(1)
	}
	// Print a summary of the job
	j.Print()
}
