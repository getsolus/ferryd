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
	"fmt"
	"os"
)

// CreateDir makes a directory if it doesn't exist
func CreateDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not stat directory '%s', reason: %s", path, err.Error())
		}
		if err = os.Mkdir(path, 0755); err != nil {
			return fmt.Errorf("could not create directory '%s', reason: %s", path, err.Error())
		}
	}
	if err := os.Chown(path, os.Getuid(), os.Getgid()); err != nil {
		return fmt.Errorf("failed to set directory ownership, reason: %s", err.Error())
	}
	return nil
}
