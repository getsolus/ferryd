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
	"sort"

	"github.com/spf13/cobra"

	"github.com/getsolus/ferryd/src/libferry"
)

var listReposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List the currently known repositories",
	Long:  "List the currently known repositories",
	Run:   listRepos,
}

func init() {
	ListCmd.AddCommand(listReposCmd)
}

func listRepos(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "list-repos takes no arguments\n")
		return
	}

	client := libferry.NewClient(socketPath)
	defer client.Close()

	repos, err := client.GetRepos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting repos: %v\n", err)
		return
	}
	sort.Strings(repos)
	if len(repos) == 0 {
		fmt.Printf("No repositories have been created yet.\n\n")
		fmt.Println("Create one with 'ferryctl create-repo $name'.")
		return
	}
	fmt.Printf("Currently registered repositories: \n\n")
	for _, repo := range repos {
		fmt.Printf(" * %v\n", repo)
	}
}
