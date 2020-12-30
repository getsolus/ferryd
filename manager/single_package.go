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
	"github.com/getsolus/ferryd/manifest"
	"github.com/getsolus/ferryd/repo"
)

/*********************/
/* PACKAGE FUNCTIONS */
/*********************/

// TransitPackage processes an incoming package manifest, adding the package to "instant transit" repos
func (m *Manager) TransitPackage(pkg string) (int, error) {
	// Check arguments
	if len(pkg) == 0 {
		return -1, errors.New("job is missing a package")
	}
	// Create new job
	j := &jobs.Job{
		Type: jobs.TransitPackage,
		Pkg:  pkg,
	}
	// Add to the DB
	return m.store.Push(j)
}

// TransitPackageExecute carries out a TransitPackage job
func (m *Manager) TransitPackageExecute(j *jobs.Job) error {
	// Check arguments
	if len(j.Pkg) == 0 {
		return errors.New("job is missing a package")
	}
	// Read the manifest
	manifest, err := manifest.NewManifest(j.Pkg)
	if err != nil {
		return fmt.Errorf("Failed to read in the manifest, reason: '%s'", err.Error())
	}
	// Verify the manifest
	if err = manifest.Verify(); err != nil {
		return fmt.Errorf("Failed to verify the manifest, reason: '%s'", err.Error())
	}
	// Create a DB transaction
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("Failed to create transaction, reason: '%s'", err.Error())
	}
	// Get the list of repos
	rs, err := repo.All(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to get the list of repos, reason: '%s'", err.Error())
	}
	// Find pool
	var pool *repo.Repo
	for _, r := range rs {
		if r.Name == "pool" {
			pool = r
			break
		}
	}
	if pool == nil {
		tx.Rollback()
		return errors.New("Could not find a DB entry for the pool")
	}
	// Copy the package files into the pool, create deltas, and add releases to the DB
	add, del, err := pool.Transit(tx, manifest)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to transit into the pool, reason: '%s'", err.Error())
	}
	// End the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Failed to end the transaction, reason: '%s'", err.Error())
	}
	// For each repo with instant_transit=true
	for _, r := range rs {
		// Skip pool
		if r.Name == "pool" {
			continue
		}
		// Create a DB transaction
		tx, err := m.db.Beginx()
		if err != nil {
			return fmt.Errorf("Failed to create transaction, reason: '%s'", err.Error())
		}
		// Copy in the new packages
		if err = r.Link(tx, add, del); err != nil {
			return fmt.Errorf("Failed to link new packages, reason: '%s'", err.Error())
		}
		// Re-Index
		if err = repo.Index(r, j, tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to reindex the repo '%s', reason: '%s'", r.Name, err.Error())
		}
		// End the transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("Failed to end the transaction, reason: '%s'", err.Error())
		}
	}
	return nil
}
