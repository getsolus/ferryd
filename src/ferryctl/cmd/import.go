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
	"path/filepath"
	"strings"
)

var importCmd = &cobra.Command{
	Use:   "import [repo] [packages]",
	Short: "Bulk import packages into repository",
	Long:  "Add packages in bulk to the named repository",
	Run:   importEx,
}

func init() {
	RootCmd.AddCommand(importCmd)
}

// GetEopkgs will utilize the provided path to get any eopkgs.
// If the provided path is a file, we will validate it is an eopkg,. If the provided path is a directory, we'll recursively fetch eopkgs inside its contents.
func GetEopkgs(providedPath string) (files []string, getErr error) {
	var pathFile *os.File

	if pathFile, getErr = os.Open(providedPath); getErr == nil { // If we successfully opened the "file"
		var fileInfo os.FileInfo

		if fileInfo, getErr = pathFile.Stat(); getErr == nil { // If we successfully stat'ed the file
			if fileInfo.IsDir() { // If this is a directory
				var nestedNames []string

				if nestedNames, getErr = pathFile.Readdirnames(-1); getErr == nil { // Get all the file contents
					for _, nestedFileName := range nestedNames { // For each nestedName
						nestedFiles, nestedErr := GetEopkgs(filepath.Join(providedPath, nestedFileName)) // Get our nested contents

						if nestedErr == nil { // No error
							files = append(files, nestedFiles...) // Append our nested files
						} else { // Error happened in nested GetEopkgs
							getErr = nestedErr
						}
					}
				}
			} else { // If this is a file
				if strings.HasSuffix(fileInfo.Name(), "x86_64.eopkg") { // Not a delta
					files = append(files, providedPath)
				}
			}
		}

		pathFile.Close()
	}

	return
}

func importEx(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: import [repo] [packages]\n")
		return
	}

	client := libferry.NewClient(socketPath)
	defer client.Close()

	repoID := args[0]
	var packages []string
	for i := 1; i < len(args); i++ {
		f, err := filepath.Abs(args[i])

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to probe: %s: %v\n", f, err)
			return
		}

		if eopkgs, getErr := GetEopkgs(f); getErr == nil { // If we successfully got eopkgs
			packages = append(packages, eopkgs...)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to get eopkgs: %v", getErr)
			return
		}
	}

	if err := client.ImportPackages(repoID, packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error while importing packages: %v\n", err)
		return
	}
}
