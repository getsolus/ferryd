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
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"github.com/jmoiron/sqlx"
)

// Manager is responsible for carrying out changes to the repositories
type Manager struct {
	db    *sqlx.DB
	store *jobs.Store
	pool  *Pool
}

/**************************/
/* MANAGER ONLY FUNCTIONS */
/**************************/

// NewManager creates the manager and opens the repo DB
func NewManager(store *jobs.Store) *Manager {
	manager := &Manager{
		store: store,
	}
	// Open the DB
	manager.db = repo.OpenDB()
	return manager
}

// Close shuts-down the manager and closes its database
func (m *Manager) Close() error {
	m.pool.Close()
	if err := m.store.Close(); err != nil {
		return err
	}
	return m.db.Close()
}
