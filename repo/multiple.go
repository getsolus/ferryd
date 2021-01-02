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

package repo

import (
	"errors"
	"github.com/getsolus/ferryd/jobs"
	"github.com/jmoiron/sqlx"
)

// CherryPick syncs a single package from this repo to another
func CherryPick(left, right *Repo, j *jobs.Job, tx *sqlx.Tx) (d *Diff, err error) {
	if d, err = Compare(left, right, j, tx); err != nil {
		return
	}
	// TODO: Implement
	return nil, errors.New("Function not implemented")
}

// Compare the contents of this repo to another
func Compare(left, right *Repo, j *jobs.Job, tx *sqlx.Tx) (d *Diff, err error) {
	// TODO: Implement, checking for j.Pkg to be set to limit the size of the comparison
	return nil, errors.New("Function not implemented")
}

// Sync all packages from this repo to another
func Sync(left, right *Repo, j *jobs.Job, tx *sqlx.Tx) (d *Diff, err error) {
	// TODO: Implement
	return nil, errors.New("Function not implemented")
}
