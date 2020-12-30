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
	"github.com/getsolus/ferryd/config"
	"github.com/getsolus/ferryd/util"
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

// Repo is an entry in the Repo Table
type Repo struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	InstantTransit bool   `db:"instant_transit"`
}

// Get retrieves a single repo by name
func Get(tx *sqlx.Tx, name string) (r *Repo, err error) {
	r = &Repo{}
	// Get the DB entry
	if err = tx.Get(r, GetSingle, name); err != nil {
		return
	}
	// Create the repo directory if missing
	rp := filepath.Join(append(config.Current.RepoPath(), name)...)
	if err = util.CreateDir(rp); err != nil {
		return
	}
	// Create the repo directory if missing
	ap := filepath.Join(append(config.Current.AssetPath(), name)...)
	if err = util.CreateDir(ap); err != nil {
		return
	}
	return
}

// All retrieves a list of all the repos in the DB
func All(tx *sqlx.Tx) (rs []*Repo, err error) {
	rs = make([]*Repo, 0)
	err = tx.Select(&rs, GetAll)
	return
}

// Create inserts a new repo into the database
func (r *Repo) Create(tx *sqlx.Tx) error {
	res, err := tx.NamedExec(Insert, r)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = int(id)
	return nil
}

// Summarize gets a summary for this repo
func (r *Repo) Summarize(tx *sqlx.Tx) (s Summary, err error) {
	s.Name = r.Name
	if err = tx.Get(&s.Packages, PackageCount, r.ID); err != nil {
		return
	}
	if err = tx.Get(&s.Deltas, DeltaCount, r.ID); err != nil {
		return
	}
	err = tx.Get(&s.Size, GetSize, r.ID)
	return
}
