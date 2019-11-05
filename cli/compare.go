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
}

// CompareRun executes the "compare" sub-command
func CompareRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*CompareArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	if err := client.Compare(args.Left, args.Right); err != nil {
		fmt.Fprintf(os.Stderr, "Error while comparing repos: %v\n", err)
		os.Exit(1)
	}
}
