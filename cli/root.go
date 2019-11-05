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
	"github.com/DataDrake/cli-ng/cmd"
)

// Root is the main entry point into ferryd
var Root *cmd.RootCMD

// GlobalFlags contains the flags for all commands
type GlobalFlags struct {
	Socket string `short:"s" arg:"true" long:"socket" desc:"Set the socket path to talk to ferryd"`
}

func init() {
	Root = &cmd.RootCMD{
		Name:  "ferryd",
		Short: "ferryd is the Solus package repository tool",
		Flags: &GlobalFlags{"/run/ferryd.sock"},
	}

	Root.RegisterCMD(&cmd.Help)
	Root.RegisterCMD(Version)
	// Daemon
	Root.RegisterCMD(Daemon)
	Root.RegisterCMD(Stop)
	Root.RegisterCMD(Restart)
	Root.RegisterCMD(Status)
	// Job Management
	Root.RegisterCMD(ResetCompleted)
	Root.RegisterCMD(ResetFailed)
	Root.RegisterCMD(ResetQueued)
	// Single-Repo
	Root.RegisterCMD(Check)
	Root.RegisterCMD(Create)
	Root.RegisterCMD(Delta)
	Root.RegisterCMD(Import)
	Root.RegisterCMD(Index)
	Root.RegisterCMD(Rescan)
	Root.RegisterCMD(Remove)
	Root.RegisterCMD(TrimPackages)
	Root.RegisterCMD(TrimObsoletes)
	// Multiple-Repo
	Root.RegisterCMD(CherryPick)
	Root.RegisterCMD(Clone)
	Root.RegisterCMD(Compare)
	Root.RegisterCMD(List)
	Root.RegisterCMD(Sync)
}
