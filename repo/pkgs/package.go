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

package pkgs

import (
	"github.com/jmoiron/sqlx"
)

// Package is an entry in the Package Table
type Package struct {
	RepoID    int
	ReleaseID int
}

// Save adds a new entry to the package table
func (p *Package) Save(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(Insert, p)
	return err
}

// Remove deletes an exact match for this package in the DB
func (p *Package) Remove(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(Remove, p)
	return err
}
