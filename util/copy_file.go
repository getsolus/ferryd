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
	"io"
	"os"
)

func actualCopyFile(source, dest string, st os.FileInfo) (err error) {
	var src *os.File
	var dst *os.File
	// Open the source file for reading
	if src, err = os.Open(source); err != nil {
		return
	}
	defer src.Close()
	// Open the destination file for writing
	if dst, err = os.OpenFile(dest, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, st.Mode()); err != nil {
		return
	}
	// Copy the file contents
	if _, err = io.Copy(dst, src); err != nil {
		dst.Close()
		return
	}
	// Set the user and group ownership
	dst.Chown(os.Getuid(), os.Getgid())
	// Close the file
	dst.Close()
	// Set the same modification time
	os.Chtimes(dest, st.ModTime(), st.ModTime())
	return nil
}

// CopyFile will copy the file and permissions to the new target
func CopyFile(source, dest string) error {
	var err error
	var st os.FileInfo
	// Get the details for the source file
	if st, err = os.Stat(source); err != nil {
		return err
	}
	return actualCopyFile(source, dest, st)
}
