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

var pullRepoCmd = &cobra.Command{
	Use:   "pull [sourceRepo] [destRepo]",
	Short: "pull an existing repository",
	Long:  "Clone an existing repository into a new repository",
	Run:   pullRepo,
}

func init() {
	RootCmd.AddCommand(pullRepoCmd)
}

func pullRepo(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: pull [sourceRepo] [destRepo]\n")
		return
	}

	client := libferry.NewClient(socketPath)
	defer client.Close()

	if err := client.PullRepo(args[0], args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error while pulling repo: %v\n", err)
		return
	}
}
