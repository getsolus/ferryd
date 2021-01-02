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

package release

import (
	"github.com/getsolus/ferryd/repo/archive"
	"github.com/jmoiron/sqlx"
	"sort"
)

// Releases is a list of Release
type Releases []Release

// GetReleases retrieves all of the Releases for a package in a repo
func GetReleases(tx *sqlx.Tx, repo, pkg string) (rs Releases, err error) {
	var as archive.Archives
	if err = tx.Get(&as, GetPkgArchives, repo, pkg); err != nil {
		return
	}
	sort.Sort(as)
	// Sort Archives into Releases
	var r *Release
	for _, a := range as {
		if r == nil {
			r = &Release{}
		} else if r.Package() != a.Package {
			rs = append(rs, *r)
			r = &Release{}
		}
		if a.IsPackage() {
			d := a.Copy()
			r.Pkg = &d
		} else if a.IsDelta() {
			r.Deltas = append(r.Deltas, a.Copy())
		}
	}
	// Sort releases by Release number
	sort.Sort(rs)
	return
}

// Len returns the length of the Releases
func (rs Releases) Len() int {
	return len(rs)
}

// Less compares the "release" number for two ReleaseSets
func (rs Releases) Less(i, j int) bool {
	return rs[i].Number() < rs[j].Number()
}

// Swap exchanges releases for sorting
func (rs Releases) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}
