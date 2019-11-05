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

// Sync fulfills the "sync" sub-command
var Sync = &cmd.CMD{
	Name:  "sync",
	Alias: "sr",
	Short: "Sync an existing repository into another repository",
	Args:  &SyncArgs{},
	Run:   SyncRun,
}

// SyncArgs are the arguments to the "sync" sub-command
type SyncArgs struct {
	Source string `desc:"Repo to sync from"`
	Dest   string `desc:"Repo to sync into"`
}

// SyncRun executes the "sync" sub-command
func SyncRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	args := c.Args.(*SyncArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	if err := client.Sync(args.Source, args.Dest); err != nil {
		fmt.Fprintf(os.Stderr, "Error while syncing: %v\n", err)
		os.Exit(1)
	}
}
