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
	"github.com/spf13/cobra"
	"libferry"
	"os"
)

var resetCompletedCmd = &cobra.Command{
	Use:   "completed",
	Short: "reset completed logs",
	Long:  "Purge the logs for completed jobs",
	Run:   resetCompleted,
}

func init() {
	ResetCmd.AddCommand(resetCompletedCmd)
}

func resetCompleted(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "reset completed takes no arguments\n")
		return
	}

	client := libferry.NewClient(socketPath)
	defer client.Close()

	if err := client.ResetCompleted(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while resetting completed log: %v\n", err)
		return
	}
}