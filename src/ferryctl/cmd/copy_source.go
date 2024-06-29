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
	"strconv"

	"github.com/spf13/cobra"

	"github.com/getsolus/ferryd/src/libferry"
)

var (
	skipIndex bool
)

var copySourceCmd = &cobra.Command{
	Use:   "source [fromRepo] [targetRepo] [sourceName] [releaseNumber]",
	Short: "copy packages by source name",
	Long:  "Remove an existing package set in the ferryd instance",
	Run:   copySource,
}

func init() {
	CopyCmd.AddCommand(copySourceCmd)
	CopyCmd.PersistentFlags().BoolVarP(&skipIndex, "skip-index", "si", false, "Skip updating the index of the target")
}

func copySource(cmd *cobra.Command, args []string) {
	var (
		repoID        string
		targetID      string
		sourceID      string
		sourceRelease int
	)

	switch len(args) {
	case 3:
		repoID = args[0]
		targetID = args[1]
		sourceID = args[2]
		sourceRelease = -1
	case 4:
		repoID = args[0]
		targetID = args[1]
		sourceID = args[2]

		release, err := strconv.ParseInt(args[3], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid integer: %v\n", err)
			return
		}
		if release < 1 {
			fmt.Fprintf(os.Stderr, "Release should be higher than 1\n")
			return
		}
		sourceRelease = int(release)
	default:
		fmt.Fprintf(os.Stderr, "usage: [fromRepo] [targetRepo] [sourceName] [releaseNumber]\n")
		return
	}

	client := libferry.NewClient(socketPath)
	defer client.Close()

	if err := client.CopySource(repoID, targetID, sourceID, sourceRelease, skipIndex); err != nil {
		fmt.Fprintf(os.Stderr, "Error while copying source: %v\n", err)
		return
	}
}
