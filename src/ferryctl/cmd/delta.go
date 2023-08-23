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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/getsolus/ferryd/src/libferry"
)

var deltaCmd = &cobra.Command{
	Use:   "delta [repoName]",
	Short: "Create deltas",
	Long:  "Rebuild all deltas for a repo",
	Run:   delta,
}

func init() {
	RootCmd.AddCommand(deltaCmd)
}

func delta(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: delta [repoName]\n")
		return
	}

	client := libferry.NewClient(socketPath)
	defer client.Close()

	if err := client.DeltaRepo(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating deltas: %v\n", err)
		return
	}
}
