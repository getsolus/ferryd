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

package archive

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"path/filepath"
	"strings"
)

// Status indicates whether or a release should be added, removed, or remain unchanged
type Status int

const (
	// StatusUnchanged indicates that this Archive remains unchanged in the target Repo
	StatusUnchanged Status = iota
	// StatusAdded indicates that this Archive should be added to the target Repo
	StatusAdded
	// StatusModified indicates that this Archive shoud be updated in the target Repo
	StatusModified
	// StatusRemoved indicates that this Archive should be removed from the target Repo
	StatusRemoved
)

var (
	// ErrInvalidArchive indicates that the relevant Archive DB entry is malformed
	ErrInvalidArchive = errors.New("invalid archive")
	// ErrArchiveTypeMismatch indicates that two compared Archives do not have the same type
	ErrArchiveTypeMismatch = errors.New("archive type mismatch")
)

// Archive represents a single Archive of a package in the repos
type Archive struct {
	ID      int    `db:"id"`
	Package string `db:"package"`
	URI     string `db:"uri"`
	Size    int    `db:"size"`
	Hash    string `db:"hash"`
	Release int    `db:"release"`
	To      int    `db:"to_release"`
	Meta    []byte `db:"meta"`
	Status  Status `db:"-"`
}

// Copy creates a duplicate of an existing Archive
func (a Archive) Copy() (ret Archive) {
	ret = a
	meta := make([]byte, len(a.Meta))
	copy(meta, a.Meta)
	ret.Meta = meta
	ret.Status = StatusUnchanged
	return
}

// Name returns the filename of this Archive
func (a *Archive) Name() (name string, err error) {
	if !a.IsValid() {
		err = ErrInvalidArchive
		return
	}
	name = filepath.Base(a.URI)
	return
}

// IsPackage checks if this is a valid Package Archive
func (a *Archive) IsPackage() bool {
	return a.Release > 0 && a.To == 0
}

// IsDelta checks if this is a valid Delta Archive
func (a *Archive) IsDelta() bool {
	return a.Release > 0 && a.To > a.Release
}

// IsValid checks if this release is valid at all
func (a *Archive) IsValid() bool {
	return a.IsDelta() || a.IsPackage()
}

// Compare this Archive with another based on sorting order
func (a Archive) Compare(a2 Archive) (res int) {
	if a.ID == a2.ID {
		res = 0
		return
	}
	if res = strings.Compare(a.Package, a2.Package); res != 0 {
		return
	}
	switch {
	case a.Release < a2.Release:
		res = -1
	case a.Release > a2.Release:
		res = 1
	case a.To < a2.To:
		res = -1
	case a.To > a2.To:
		res = 1
	default:
		res = 0
	}
	return
}

// Save or create an Archive with the current values
func (a *Archive) Save(tx *sqlx.Tx) (err error) {
	if a.ID == 0 {
		//Create
		res, err := tx.NamedExec(Insert, a)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		a.ID = int(id)
	} else {
		// Update
		_, err = tx.NamedExec(Update, a)
	}
	return
}

// PrintDiff prints an Archive according to its Status
func (a *Archive) PrintDiff(out io.Writer, plus, minus, mod, same string) error {
	name, err := a.Name()
	if err != nil {
		return err
	}
	switch a.Status {
	case StatusAdded:
		fmt.Fprintf(out, plus, name)
	case StatusRemoved:
		fmt.Fprintf(out, minus, name)
	case StatusModified:
		fmt.Fprintf(out, mod, name)
	default:
		fmt.Fprintf(out, same, name)
	}
	return nil
}
