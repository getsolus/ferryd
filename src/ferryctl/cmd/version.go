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

	"github.com/getsolus/ferryd/src/libferry"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Long:  "Print the ferryctl version and exit",
	Run:   printVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion(cmd *cobra.Command, args []string) {
	// Print local version
	fmt.Printf("ferry %v\n\nCopyright © 2017-2019 Solus Project\n", libferry.Version)
	fmt.Printf("Licensed under the Apache License, Version 2.0\n")
}
