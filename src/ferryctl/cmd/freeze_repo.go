//
// Copyright © 2025 Solus Project
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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getsolus/ferryd/src/libferry"
)

var (
	freezeRepoCmd = &cobra.Command{
		Use:   "freeze [repo]",
		Short: "protect a repository against changes",
		Long:  "Freeze a repository, disallowing changes to that repository",
		Run:   freezeRepo,
		Args:  cobra.ExactArgs(1),
	}
	unfreezeRepoCmd = &cobra.Command{
		Use:   "unfreeze [repo]",
		Short: "disable protection of a repository",
		Long:  "Disable a repository freeze, allowing changes to that repository",
		Run:   unfreezeRepo,
		Args:  cobra.ExactArgs(1),
	}
)

func init() {
	RootCmd.AddCommand(freezeRepoCmd, unfreezeRepoCmd)
}

func freezeRepo(_ *cobra.Command, args []string) {
	client := libferry.NewClient(socketPath)
	defer client.Close()

	if err := client.FreezeRepo(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error while freezing repo: %v\n", err)
		return
	}
}

func unfreezeRepo(_ *cobra.Command, args []string) {
	client := libferry.NewClient(socketPath)
	defer client.Close()

	if err := client.UnfreezeRepo(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error while unfreezing repo: %v\n", err)
		return
	}
}
