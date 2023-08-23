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

package core

import (
	"os"
	"path/filepath"

	"github.com/getsolus/ferryd/src/libdb"
)

// A Manager is the the singleton responsible for slip management
type Manager struct {
	db   libdb.Database     // Our main database
	ctx  *Context           // Context shares all our path assignments
	pool *Pool              // Our main pool for eopkgs
	repo *RepositoryManager // Repo management

	IncomingPath string // Incoming directory
}

// NewManager will attempt to instaniate a manager for the given path,
// which will yield an error if the database cannot be opened for access.
func NewManager(path string) (*Manager, error) {
	ctx, err := NewContext(path)
	if err != nil {
		return nil, err
	}

	// Open the database if we can
	db, err := libdb.Open(ctx.DbPath)
	if err != nil {
		return nil, err
	}

	// Need incoming to monitor uploads
	incomingPath := filepath.Join(ctx.BaseDir, IncomingPathComponent)
	if err := os.MkdirAll(incomingPath, 00755); err != nil {
		return nil, err
	}

	m := &Manager{
		db:           db,
		ctx:          ctx,
		pool:         &Pool{},
		repo:         &RepositoryManager{},
		IncomingPath: incomingPath,
	}

	// Initialise the buckets in a one-time
	if err = m.initComponents(); err != nil {
		m.Close()
		return nil, err
	}

	return m, nil
}

// initComponents will ensure all initial buckets are create in the toplevel
// namespace, to require less complexity further down the line
func (m *Manager) initComponents() error {
	// Components to bring up
	components := []Component{
		m.pool,
		m.repo,
	}

	// Create all root-level buckets in a single transaction
	return m.db.Update(func(db libdb.Database) error {
		for _, component := range components {
			if err := component.Init(m.ctx, db); err != nil {
				return err
			}
		}
		return nil
	})
}

// Close will close and clean up any associated resources, such as the
// underlying database.
func (m *Manager) Close() {
	if m.db == nil {
		return
	}
	// Components to tear down
	components := []Component{
		m.pool,
		m.repo,
	}
	for _, component := range components {
		component.Close()
	}
	m.db.Close()
	m.db = nil
}
