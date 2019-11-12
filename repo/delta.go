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
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// Delta is an entry in the Delta Table
type Delta struct {
	ID          int            `db:"id"`
	PackageName sql.NullString `db:"package_name"`
	URI         sql.NullString `db:"uri"`
	Size        int            `db:"size"`
	Hash        sql.NullString `db:"hash"`
	FromRelease int            `db:"from_release"`
	ToRelease   int            `db:"to_release"`
}

// Insert creates a new delta in the DB
func (d *Delta) Insert(tx *sqlx.Tx, repoID int) error {
	// Insert New Delta
	resp, err := tx.NamedExec(insertDelta, d)
	if err != nil {
		return err
	}
	// Get ID of new Delta record
	id, err := resp.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = int(id)
	// Insert New RepoDelta to pair with repo
	rd := &RepoDelta{
		RepoID:  repoID,
		DeltaID: int(id),
	}
	return rd.Insert(tx)
}

// Equal checks if this delta is equal to another
func (d *Delta) Equal(d2 *Delta) bool {
	same := NullStringEqual(d.PackageName, d2.PackageName)
	same = same && NullStringEqual(d.URI, d2.URI)
	same = same && d.Size == d2.Size
	same = same && NullStringEqual(d.Hash, d2.Hash)
	same = same && d.FromRelease == d2.FromRelease
	same = same && d.ToRelease == d2.ToRelease
	return same
}
