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

package manager

import (
	"errors"
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
)

/***************************/
/* MULTIPLE REPO FUNCTIONS */
/***************************/

// CherryPick syncs a single package from one repo to another
func (m *Manager) CherryPick(src, dest, pkg string) (int, error) {
	// Validate the arguments
	if len(src) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if len(dest) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	if len(pkg) == 0 {
		return -1, errors.New("job is missing a package name")
	}
	// Create a new job instance
	j := &jobs.Job{
		Type: jobs.CherryPick,
		Src:  src,
		Dst:  dest,
		Pkg:  pkg,
	}
	// Add the job to the DB
	id, err := m.store.Push(j)
	// Return the id of the new Job
	return int(id), err
}

// CherryPickExecute carries out a CherryPick Job
func (m *Manager) CherryPickExecute(j *jobs.Job) error {
	// Validate the arguments
	if len(j.Src) == 0 {
		return errors.New("job is missing a source repo")
	}
	if len(j.Dst) == 0 {
		return errors.New("job is missing a destination repo")
	}
	if len(j.Pkg) == 0 {
		return errors.New("job is missing a package name")
	}
	// Begin a DB Transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to start DB transaction, reason: '%s'", err.Error())
	}
	// Get the source Repo instance
	src, err := repo.Get(tx, j.Src)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the source Repo entry from the DB, reason: '%s'", err.Error())
	}
	// Get the source Repo instance
	dst, err := repo.Get(tx, j.Dst)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the destination Repo entry from the DB, reason: '%s'", err.Error())
	}
	// CherryPick a single package from one repo to the other
	var diff *repo.Diff
	if diff, err = src.CherryPick(tx, dst, j.Pkg); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to cherry pick '%s', reason: '%s'", j.Pkg, err.Error())
	}
	// End the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit the transaction, reason: '%s'", err.Error())
	}
	// Save the diff into the job
	j.Results, err = diff.MarshalBinary()
	if err != nil {
		return fmt.Errorf("Failed to convert Diff to binary for saving, reason: '%s'", err.Error())
	}
	return nil
}

// Clone creates a new repo as a copy of and existing repo
func (m *Manager) Clone(src, dst string) (int, error) {
	// Validate the arguments
	if len(src) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if len(dst) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	// Create a new job instance
	j := &jobs.Job{
		Type: jobs.Clone,
		Src:  src,
		Dst:  dst,
	}
	// Add the job to the DB
	id, err := m.store.Push(j)
	// Return the ID of the new job
	return int(id), err
}

// CloneExecute carries out a clone job
func (m *Manager) CloneExecute(j *jobs.Job) error {
	// Create the new repo
	if err := m.CreateExecute(j); err != nil {
		return err
	}
	// Sync from the existing repo
	return m.SyncExecute(j)
}

// Compare reports on the differences between two repos
func (m *Manager) Compare(left, right string) (int, error) {
	// Validate the arguments
	if len(left) == 0 {
		return -1, errors.New("job is missing a left repo")
	}
	if len(right) == 0 {
		return -1, errors.New("job is missing a right repo")
	}
	// Create a new job instance
	j := &jobs.Job{
		Type: jobs.Compare,
		Src:  left,
		Dst:  right,
	}
	// Add the job to the DB
	id, err := m.store.Push(j)
	// Return the id of the new Job
	return int(id), err
}

// CompareExecute carries out a Sync Job
func (m *Manager) CompareExecute(j *jobs.Job) error {
	// Validate the arguments
	if len(j.Src) == 0 {
		return errors.New("job is missing a left repo")
	}
	if len(j.Dst) == 0 {
		return errors.New("job is missing a right repo")
	}
	// Begin a DB Transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to start DB transaction, reason: '%s'", err.Error())
	}
	// Get the source Repo instance
	left, err := repo.Get(tx, j.Src)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the left Repo entry from the DB, reason: '%s'", err.Error())
	}
	// Get the source Repo instance
	right, err := repo.Get(tx, j.Dst)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the right Repo entry from the DB, reason: '%s'", err.Error())
	}
	// Compare the left repo with the right repo
	var diff *repo.Diff
	if diff, err = left.Compare(tx, right); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to compare repos, reason: '%s'", err.Error())
	}
	// End the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit the transaction, reason: '%s'", err.Error())
	}
	// Save the diff into the Job
	j.Results, err = diff.MarshalBinary()
	if err != nil {
		return fmt.Errorf("Failed to convert Diff to binary for saving, reason: '%s'", err.Error())
	}
	return nil
}

// Sync compares two repos and makes changes so that "new" matches "old"
func (m *Manager) Sync(src, dst string) (int, error) {
	// Validate the arguments
	if len(src) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if len(dst) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	// Create a new job instance
	j := &jobs.Job{
		Type: jobs.Sync,
		Src:  src,
		Dst:  dst,
	}
	// Add the job to the DB
	id, err := m.store.Push(j)
	// Return the id of the new job
	return int(id), err
}

// SyncExecute carries out a Sync job
func (m *Manager) SyncExecute(j *jobs.Job) error {
	// Validate the arguments
	if len(j.Src) == 0 {
		return errors.New("job is missing a source repo")
	}
	if len(j.Dst) == 0 {
		return errors.New("job is missing a destination repo")
	}
	// Begin a DB Transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to start DB transaction, reason: '%s'", err.Error())
	}
	// Get the source Repo instance
	src, err := repo.Get(tx, j.Src)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the source Repo entry from the DB, reason: '%s'", err.Error())
	}
	// Get the source Repo instance
	dst, err := repo.Get(tx, j.Dst)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the destination Repo entry from the DB, reason: '%s'", err.Error())
	}
	// Sync all packages from one repo to the other
	var diff *repo.Diff
	if diff, err = src.Sync(tx, dst); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to sync, reason: '%s'", err.Error())
	}
	// End the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit the transaction, reason: '%s'", err.Error())
	}
	// Save the Diff into the job results
	j.Results, err = diff.MarshalBinary()
	if err != nil {
		return fmt.Errorf("Failed to convert Diff to binary for saving, reason: '%s'", err.Error())
	}
	return nil
}

// Repos provides a summary of all available repos
func (m *Manager) Repos() (l repo.FullSummary, err error) {
	var s repo.Summary
	// start tx
	tx, err := m.db.Beginx()
	if err != nil {
		return
	}
	// get list of repos
	rs, err := repo.All(tx)
	if err != nil {
		goto CLEANUP
	}
	// get summary for each repo
	for _, r := range rs {
		s, err = r.Summarize(tx)
		if err != nil {
			goto CLEANUP
		}
		l = append(l, s)
	}

CLEANUP:
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return
}
