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

// ResetCompleted fulfills the "reset-completed" sub-command
var ResetCompleted = &cmd.CMD{
	Name:  "reset-completed",
	Alias: "rc",
	Short: "Reset the list of completed jobs",
	Args:  &ResetCompletedArgs{},
	Run:   ResetCompletedRun,
}

// ResetCompletedArgs are the arguments to the "reset-completed" sub-command
type ResetCompletedArgs struct{}

// ResetCompletedRun executes the "reset-completed" sub-command
func ResetCompletedRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	// args  := c.Args.(*ResetCompletedArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	if err := client.ResetCompleted(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while resetting complete jobs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully reset completed jobs")
}
