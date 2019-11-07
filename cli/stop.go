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

// Stop fulfills the "stop" sub-command
var Stop = &cmd.CMD{
	Name:  "stop",
	Alias: "k",
	Short: "Stop a running ferryd daemon",
	Args:  &StopArgs{},
	Run:   StopRun,
}

// StopArgs are the arguments to the "stop" sub-command
type StopArgs struct{}

// StopRun executes the "stop" sub-command
func StopRun(r *cmd.RootCMD, c *cmd.CMD) {
	flags := r.Flags.(*GlobalFlags)
	//args  := c.Args.(*StopArgs)

	client := v1.NewClient(flags.Socket)
	defer client.Close()

	resp, err := client.Stop()
    if err != nil {
		fmt.Fprintf(os.Stderr, "Error while stopping daemon: %v\n", err)
		os.Exit(1)
	}
    if len(resp.Errors) > 0 {
		fmt.Fprintln(os.Stderr, "Error while stopping daemon:")
        resp.Print(os.Stderr)
        os.Exit(1)
    }
    fmt.Fprintln("Daemon has been stopped successfully.")
}