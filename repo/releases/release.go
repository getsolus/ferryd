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

package releases

import (
	"github.com/jmoiron/sqlx"
)

// Release represents a single Release of a package in the repos
type Release struct {
	ID      int    `db:"id"`
	Package string `db:"package"`
	URI     string `db:"uri"`
	Size    int    `db:"size"`
	Hash    string `db:"hash"`
	Release int    `db:"release"`
	From    int    `db:"from_release"`
	Meta    []byte `db:"meta"`
}

// IsPackage checks if this is a valid Release for a Package
func (r *Release) IsPackage() bool {
	return (r.From == 0) && (len(r.Meta) == 0)
}

// IsDelta checks if this is a valid Release for a Delta
func (r *Release) IsDelta() bool {
	return r.From != 0
}

// IsValid checks if this release is valid at all
func (r *Release) IsValid() bool {
	return r.IsDelta() || r.IsPackage()
}

// Save or create a release entry with the current values
func (r *Release) Save(tx *sqlx.Tx) error {
	if r.ID == 0 {
		//Create
		res, err := tx.NamedExec(Insert, r)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		r.ID = int(id)
	} else {
		// Update
		_, err := tx.NamedExec(Update, r)
		return err
	}
	return nil
}
