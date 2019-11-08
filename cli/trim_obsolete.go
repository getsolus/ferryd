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

// TrimObsolete fulfills the "trim-obsolete" sub-command
var TrimObsolete = &cmd.CMD{
	Name:  "trim-obsolete",
	Alias: "to",
	Short: "Remove all obsolete packages from a repo",
	Args:  &TrimObsoleteArgs{},
	Run:   TrimObsoleteRun,
}

// TrimObsoleteArgs are the arguments to the "trim-obsolete" sub-command
type TrimObsoleteArgs struct {
	Repo string `desc:"Repo to trim"`
}

// TrimObsoleteRun executes the "trim-obsolete" sub-command
func TrimObsoleteRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*TrimObsoleteArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	j, err := client.TrimObsolete(args.Repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while trimming obsolete packages in repo: %v\n", err)
		os.Exit(1)
	}
	j.Print()
}
