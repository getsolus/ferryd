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
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"os"
)

// Compare fulfills the "compare" sub-command
var Compare = &cmd.CMD{
	Name:  "compare",
	Alias: "diff",
	Short: "Calculate the differences between two repos",
	Args:  &CompareArgs{},
	Run:   CompareRun,
}

// CompareArgs are the arguments to the "compare" sub-command
type CompareArgs struct {
	Left  string `desc:"First repo to compare"`
	Right string `desc:"Second repo to compare"`
	Full  bool   `desc:"Print the full diff and not just the changes"`
}

// CompareRun executes the "compare" sub-command
func CompareRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*CompareArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Run the job
	var j *jobs.Job
	var err error
	if j, err = client.Compare(args.Left, args.Right); err != nil {
		fmt.Fprintf(os.Stderr, "Error while comparing repos: %v\n", err)
		os.Exit(1)
	}
	// Print the job summary
	j.Print()
	// Decode the Diff
	var d *repo.Diff
	if err = d.UnmarshalBinary(j.Results); err != nil {
		fmt.Fprintf(os.Stderr, "Error while decoding diff: %v\n", err)
		os.Exit(1)
	}
	// Print the diff
	d.Print(os.Stdout, args.Full, !flags.NoColor)
}
