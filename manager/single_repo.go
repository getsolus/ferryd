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

package manager

import (
	"errors"
	"fmt"
	"github.com/getsolus/ferryd/config"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"github.com/getsolus/ferryd/util"
	"os"
	"path/filepath"
)

/*************************/
/* SINGLE REPO FUNCTIONS */
/*************************/

// Check compares an existing repo on Disk with its DB
func (m *Manager) Check(name string) (int, error) {
	// Validate the job arguments
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	// Create the job
	j := &jobs.Job{
		Type: jobs.Check,
		Src:  name,
	}
	// Add the job to the DB
	id, err := m.store.Push(j)
	// Return the new Job ID
	return int(id), err
}

// CheckExecute carries out a Check job
func (m *Manager) CheckExecute(j *jobs.Job) error {
	var d *repo.Diff
	// Validate arguments
	if len(j.Src) == 0 {
		return errors.New("job is missing a source repo")
	}
	// Start transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}
	// Get repo by name
	r, err := repo.Get(tx, j.Src)
	if err != nil {
		goto ROLLBACK
	}
	// Run the check
	d, err = r.Check(tx)
	if err != nil {
		goto ROLLBACK
	}
	// End transaction
	tx.Commit()
	// Save the result
	j.Results, err = d.MarshalBinary()
	return err

ROLLBACK:
	tx.Rollback()
	return err
}

// Create sets up a new repo
func (m *Manager) Create(name string, instant bool) (int, error) {
	// Validate the job arguments
	if len(name) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	// Create a new job
	max := 0
	if instant {
		max = 1
	}
	j := &jobs.Job{
		Type: jobs.Create,
		Dst:  name,
		Max:  max,
	}
	// Add it to the DB
	id, err := m.store.Push(j)
	// Return the new Job ID
	return int(id), err
}

// CreateExecute carries out a Create job
func (m *Manager) CreateExecute(j *jobs.Job) error {
	// Validate the job arguments
	if len(j.Dst) == 0 {
		return errors.New("job is missing a destination repo")
	}
	// Create the repo directory
	repoDir := append(config.Current.RepoPath(), j.Dst)
	if err := os.Mkdir(filepath.Join(repoDir...), 00755); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("repo directory for '%s' already exists", j.Dst)
		}
		return err
	}
	// Create the assets directory
	assetsDir := filepath.Join(append(config.Current.AssetPath(), j.Dst)...)
	poolAssets := filepath.Join(append(config.Current.AssetPath(), "pool")...)
	if err := util.CopyDir(poolAssets, assetsDir, false); err != nil {
		return fmt.Errorf("Failed to create assets dir, reason: '%s'", err.Error())
	}
	// Add the repo to the DB
	// Create a DB transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to create transaction, reason: '%s'", err.Error())
	}
	// Create a new repo object
	r := &repo.Repo{
		Name:           j.Dst,
		InstantTransit: j.Max == 1,
	}
	// Insert into the DB
	if err = r.Create(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create repo entry in DB, reason: '%s'", err.Error())
	}
	// End the transaction
	return tx.Commit()
}

// Delta generates missing package deltas for an entire repo
func (m *Manager) Delta(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Delta,
		Dst:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// DeltaExecute carries out a Delta job
func (m *Manager) DeltaExecute(j *jobs.Job) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// DeltaPackage generates missing package deltas for a single package
func (m *Manager) DeltaPackage(dst, pkg string) (int, error) {
	if len(dst) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	if len(pkg) == 0 {
		return -1, errors.New("job is missing a package name")
	}
	j := &jobs.Job{
		Type: jobs.DeltaPackage,
		Dst:  dst,
		Pkg:  pkg,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// DeltaPackageExecute carries out a DeltaPackage job
func (m *Manager) DeltaPackageExecute(j *jobs.Job) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Import adds an existing repo to the database
func (m *Manager) Import(name string, instant bool) (int, error) {
	// Validate the arguments
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	// Create a new job instance
	max := 0
	if instant {
		max = 1
	}
	j := &jobs.Job{
		Type: jobs.Import,
		Src:  name,
		Max:  max,
	}
	// Insert the new job into the DB
	id, err := m.store.Push(j)
	// Return the job ID
	return int(id), err
}

// ImportExecute carries out an Import job
func (m *Manager) ImportExecute(j *jobs.Job) error {
	// Validate the job arguments
	if len(j.Dst) == 0 {
		return errors.New("job is missing a destination repo")
	}
	// Create the repo directory
	repoDir := append(config.Current.RepoPath(), j.Dst)
	if _, err := os.Stat(filepath.Join(repoDir...)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("repo directory for '%s' does not exist", j.Dst)
		}
		return err
	}
	// Add the repo to the DB
	// Create a DB transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to create transaction, reason: '%s'", err.Error())
	}
	// Create a new repo object
	r := &repo.Repo{
		Name:           j.Dst,
		InstantTransit: j.Max == 1,
	}
	// Insert into the DB
	if err = r.Create(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create repo entry in DB, reason: '%s'", err.Error())
	}
	// End the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to create repo entry in DB, reason: '%s'", err.Error())
	}
	// Scan and add all of the package to the DB
	return m.RescanExecute(j)
}

// Index generates a new package index
func (m *Manager) Index(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Index,
		Dst:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// IndexExecute carries out an Index job
func (m *Manager) IndexExecute(j *jobs.Job) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// Remove deletes a repo from the DB
func (m *Manager) Remove(name string) (int, error) {
	// Validate the arguments
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	// Create a new job instance
	j := &jobs.Job{
		Type: jobs.Remove,
		Src:  name,
	}
	// Add the job to the DB
	id, err := m.store.Push(j)
	// Get the ID of the new job
	return int(id), err
}

// RemoveExecute carries out a Remove job
func (m *Manager) RemoveExecute(j *jobs.Job) error {
	// Validate the arguments
	if len(j.Src) == 0 {
		return errors.New("job is missing a source repo")
	}
	// Begin a DB Transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to start DB transaction, reason: '%s'", err.Error())
	}
	// Get the Repo instance
	r, err := repo.Get(tx, j.Src)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the Repo entry from the DB, reason: '%s'", err.Error())
	}
	// Remove the repo
	if err = r.Remove(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to remove the Repo from the DB, reason: '%s'", err.Error())
	}
	// Save the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to remove the Repo from the DB, reason: '%s'", err.Error())
	}
	return nil
}

// Rescan rebuild the database for an existing repo
func (m *Manager) Rescan(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.Rescan,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// RescanExecute carries out a Rescan job
func (m *Manager) RescanExecute(j *jobs.Job) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// TrimPackages removes old package releases and their deltas
func (m *Manager) TrimPackages(name string, max int) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if max < 1 {
		return -1, errors.New("max releases must be at least 1")
	}
	j := &jobs.Job{
		Type: jobs.TrimPackages,
		Src:  name,
		Max:  max,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// TrimPackagesExecute carries out a TrimPackages job
func (m *Manager) TrimPackagesExecute(j *jobs.Job) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}

// TrimObsoletes removes obsolete packages and their deltas
func (m *Manager) TrimObsoletes(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.TrimObsoletes,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// TrimObsoletesExecute carries out the TrimObsoletes job
func (m *Manager) TrimObsoletesExecute(j *jobs.Job) error {
	// TODO: Implement
	return errors.New("Function not implemented")
}
