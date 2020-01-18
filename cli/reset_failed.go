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
	"github.com/getsolus/ferryd/api/v1"
	"os"
)

// ResetFailed fulfills the "reset-failed" sub-command
var ResetFailed = &cmd.CMD{
	Name:  "reset-failed",
	Alias: "rf",
	Short: "Reset the failed jobs log",
	Args:  &ResetFailedArgs{},
	Run:   ResetFailedRun,
}

// ResetFailedArgs are the arguments to the "reset-failed" sub-command
type ResetFailedArgs struct{}

// ResetFailedRun executes the "reset-failed" sub-command
func ResetFailedRun(r *cmd.RootCMD, c *cmd.CMD) {
	// Convert our flags
	flags := r.Flags.(*GlobalFlags)
	//args  := c.Args.(*ResetFailedArgs)
	// Create a Client
	client := v1.NewClient(flags.Socket)
	defer client.Close()
	// Send the request
	if err := client.ResetFailed(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while resetting failed jobs: %v\n", err)
		os.Exit(1)
	}
	// Report finished
	fmt.Println("Successfully reset failed jobs")
}
