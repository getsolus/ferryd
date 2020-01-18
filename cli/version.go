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
)

// Version fulfills the "version" sub-command
var Version = &cmd.CMD{
	Name:  "version",
	Alias: "v",
	Short: "Get the version of ferryd",
	Args:  &VersionArgs{},
	Run:   VersionRun,
}

// VersionArgs are the arguments to the "version" sub-command
type VersionArgs struct{}

// VersionRun executes the "version" sub-command
func VersionRun(r *cmd.RootCMD, c *cmd.CMD) {
	//flags := r.Flags.(*GlobalFlags)
	//args  := c.Args.(*VersionArgs)

	fmt.Printf("ferryd version: %s\n", v1.Version)
}
