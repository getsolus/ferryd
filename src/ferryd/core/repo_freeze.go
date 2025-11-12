//
// Copyright © 2025 Solus Project
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

package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrFrozen is returned by write commands when the repository is frozen.
var ErrFrozen = errors.New("this repository is frozen, no changes are allowed")

// Freeze configures the repository to prevent any changes.
// This can be undone using Unfreeze.
func (r *Repository) Freeze() error {
	if err := os.WriteFile(r.freezeFile(), nil, 0o0644); err != nil {
		return fmt.Errorf("freeze repository %q: %w", r.ID, err)
	}

	return nil
}

// Unfreeze undoes Freeze, allowing any changes to the repository.
func (r *Repository) Unfreeze() error {
	if err := os.Remove(r.freezeFile()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("unfreeze repository %q: %w", r.ID, err)
	}

	return nil
}

// IsFrozen returns true if the repository has been frozen using Freeze.
func (r *Repository) IsFrozen() bool {
	_, err := os.Stat(r.freezeFile())

	return err == nil
}

func (r *Repository) freezeFile() string {
	return filepath.Join(r.path, ".frozen")
}

func (r *Repository) checkWrite() error {
	if r.IsFrozen() {
		return ErrFrozen
	}

	return nil
}
