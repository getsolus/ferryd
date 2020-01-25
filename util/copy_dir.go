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

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CopyDir will copy the contents of the files from one directory to another
func CopyDir(source, dest string, recursive bool) error {
	var err error
	var fi os.FileInfo
	// Get details about the source directory
	if fi, err = os.Stat(source); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// Create destination directory
		if err = os.Mkdir(dest, fi.Mode()); err != nil {
			return err
		}
	}
	// Set ownership
	if err = os.Chown(dest, os.Getuid(), os.Getgid()); err != nil {
		return err
	}
	// Get a list of files in the source directory
	files, err := ioutil.ReadDir(source)
	for _, file := range files {
		// Generate the filepaths
		srcFile := filepath.Join(source, file.Name())
		dstFile := filepath.Join(source, file.Name())
		// Check for a directory
		if file.IsDir() {
			// handle recursion
			if recursive {
				if err = CopyDir(srcFile, dstFile, recursive); err != nil {
					return err
				}
			}
			continue
		}
		// Copy the file to destination
		if err = actualCopyFile(srcFile, dstFile, file); err != nil {
			return err
		}
	}
	return nil
}
