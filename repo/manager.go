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

package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Manager is responsible for carrying out changes to the repositories
type Manager struct {
	db *sqlx.DB
}

/**************************/
/* MANAGER ONLY FUNCTIONS */
/**************************/

// NewManager creates the manager and opens the repo DB
func NewManager() *Manager {
	manager := &Manager{}
	// TODO: Open the repo DB
	return manager
}

// Close shuts-down the manager and closes its database
func (m *Manager) Close() error {
	// TODO: Close the repo DB
	return fmt.Errorf("Function not inplemented")
}

/*************************/
/* SINGLE REPO FUNCTIONS */
/*************************/

// CreateRepo sets up a new repo
func (m *Manager) CreateRepo(name string) error {
	return fmt.Errorf("Function not inplemented")
}

// DeltaRepo generates missing package deltas
func (m *Manager) DeltaRepo(name string) error {
	return fmt.Errorf("Function not inplemented")
}

// IndexRepo generates a new package index
func (m *Manager) IndexRepo(name string) error {
	return fmt.Errorf("Function not inplemented")
}

// RemoveRepo deletes a repo from the DB and optionally removes it from disk
func (m *Manager) RemoveRepo(name string, purge bool) error {
	return fmt.Errorf("Function not inplemented")
}

// TrimRepo removes old package releases and their deltas
func (m *Manager) TrimRepo(name string, release int) error {
	return fmt.Errorf("Function not inplemented")
}

// TrimRepoObsoletes removes obsolete packages and their deltas
func (m *Manager) TrimRepoObsoletes(name, pkg string) error {
	return fmt.Errorf("Function not inplemented")
}

/***************************/
/* MULTIPLE REPO FUNCTIONS */
/***************************/

// Clone creates a new repo as a copy of and existing repo
func (m *Manager) Clone(oldName, newName string) error {
	return fmt.Errorf("Function not inplemented")
}

// Compare reports on the differences between two repos
func (m *Manager) Compare(oldName, newName string) (same, oDiff, nDiff []Package, err error) {
	return nil, nil, nil, fmt.Errorf("Function not inplemented")
}

// Sync compares two repos and makes changes so that "new" matches "old"
func (m *Manager) Sync(oldName, newName string) error {
	return fmt.Errorf("Function not inplemented")
}

/*********************/
/* PACKAGE FUNCTIONS */
/*********************/

// TransitPackage processes an incoming package manifest, adding the package to "instant transit" repos
func (m *Manager) TransitPackage(name, pkg string) error {
	return fmt.Errorf("Function not inplemented")
}
