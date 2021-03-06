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

// ResetQueue fulfills the "reset-queue" sub-command
var ResetQueue = &cmd.CMD{
	Name:  "reset-queue",
	Alias: "rq",
	Short: "Cancel all pending jobs",
	Args:  &ResetQueueArgs{},
	Run:   ResetQueueRun,
}

// ResetQueueArgs are the arguments to the "reset-queue" sub-command
type ResetQueueArgs struct{}

// ResetQueueRun executes the "reset-queue" sub-command
func ResetQueueRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	//args  := c.Args.(*ResetQueueArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Send the request
	if err := client.ResetQueued(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while cancelling queued jobs: %v\n", err)
		os.Exit(1)
	}
	// Report finished
	fmt.Println("Successfully reset queued jobs")
}
