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

// TrimPackages fulfills the "trim-packages" sub-command
var TrimPackages = &cmd.CMD{
	Name:  "trim-packages",
	Alias: "tp",
	Short: "Remove up all, but the last N releases of all packages",
	Args:  &TrimPackagesArgs{},
	Run:   TrimPackagesRun,
}

// TrimPackagesArgs are the arguments to the "trim-packages" sub-command
type TrimPackagesArgs struct {
	Repo     string `desc:"Repo to trim"`
	Releases int    `desc:"Number of releases to keep"`
}

// TrimPackagesRun executes the "trim-packages" sub-command
func TrimPackagesRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*TrimPackagesArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	j, err := client.TrimPackages(args.Repo, args.Releases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while trimming packages in repo: %v\n", err)
		os.Exit(1)
	}
	j.Print()
}
