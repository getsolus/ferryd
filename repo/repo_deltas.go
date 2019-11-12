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
	"github.com/jmoiron/sqlx"
)

// RepoDelta is an entry in the RepoDeltas Table
type RepoDelta struct {
	RepoID  int `db:"repo_id"`
	DeltaID int `db:"delta_id"`
}

// Insert creates a new RepoDelta in the DB
func (rd *RepoDelta) Insert(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(insertRepoDelta, rd)
	return err
}

// RepoDeltasSchema is the SQLite3 schema for the RepoDeltas table
const RepoDeltasSchema = `
CREATE TABLE IF NOT EXISTS repo_deltas (
    repo_id  INTEGER,
    delta_id INTEGER,
)
`

// Queries for retrieving RepoDeltas of a particular status
const (
	repoDeltas = "SELECT * FROM repos_deltas WHERE repo_id=:repo_id"
	deltaRepos = "SELECT * FROM repos_deltas WHERE delta_id=:delta_id"
)

// Query for creating a new RepoDeltas
const insertRepoDelta = `
INSERT INTO repo_deltas (
    repo_id, delta_id
) VALUES (
    repo_id, delta_id
)
`

// Queries for removing a delta_package
const (
	deleteRepoDelta     = "DELETE * FROM repo_deltas WHERE repo_id=:repo_id AND delta_id=:delta_id"
	deleteAllDeltasRepo = "DELETE * FROM repo_deltas WHERE repo_id=:repo_id"
)
