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

// Restart fulfills the "restart" sub-command
var Restart = &cmd.CMD{
	Name:  "restart",
	Alias: "r",
	Short: "Restart a running ferryd daemon",
	Args:  &RestartArgs{},
	Run:   RestartRun,
}

// RestartArgs are the arguments to the "restart" sub-command
type RestartArgs struct{}

// RestartRun executes the "restart" sub-command
func RestartRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	//args  := c.Args.(*RestartArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	if err := client.Restart(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while restarting daemon: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Daemon has been restarted successfully.")
}
