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
	"github.com/getsolus/ferryd/manifest"
	"github.com/getsolus/ferryd/repo/pkgs"
	"github.com/getsolus/ferryd/repo/releases"
	"github.com/jmoiron/sqlx"
)

// Check makes sure the DB matches disk
func (r *Repo) Check(tx *sqlx.Tx) (d *Diff, err error) {
	// TODO: Implement
	return nil, errors.New("Function not implemented")
}

// Delta generates missing deltas and removes unneeded ones
func Delta(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// DeltaPackage generates missing deltas and removes unneeded ones for a single package
func DeltaPackage(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Index regenerates the index for a repo
func Index(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Import adds all of the files in a repo to the DB
func (r *Repo) Import(tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Link updates the links for a package that has already been updated in the pool and DB
func (r *Repo) Link(tx *sqlx.Tx, add, del []*releases.Release) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Remove deletes all of the DB records for this repo
func Remove(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// Remove Packages
	if _, err := tx.Exec(pkgs.RemoveByRepo, r.ID); err != nil {
		return err
	}
	// Remove Repo record
	_, err := tx.NamedExec(RemoveRepo, r)
	return err
}

// Rescan checks for differences between the DB and disk and updated the DB
func Rescan(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Transit copies new packages into a repo, creates missing deltas, removes old deltas, and add releases to the DB
func (r *Repo) Transit(tx *sqlx.Tx, m *manifest.Manifest) (add, del []*releases.Release, err error) {
	// TODO: Implement
	err = errors.New("Function not implemented")
	return
}

// TrimObsolete removes obsolete packages for a repo
func TrimObsolete(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// TrimPackages removes packages which are older than "max" releases from the latest
func TrimPackages(r *Repo, j *jobs.Job, tx *sqlx.Tx) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}
