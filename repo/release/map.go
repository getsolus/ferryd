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

// Map contains a list of Releases, keyed by Package name
type Map map[string]Releases

// GetAllReleases retrieves all of the Releases for all packages in a repo
func GetAllReleases(tx *sqlx.Tx, repo string) (m Map, err error) {
	var as archive.Archives
	if err = tx.Get(&as, GetRepoArchives, repo); err != nil {
		return
	}
	sort.Sort(as)
	// Sort Archives into Releases
	var r *Release
	var rs Releases
	for _, a := range as {
		if r == nil {
			r = &Release{}
		} else if r.Package() != a.Package {
			m[r.Package()] = rs
			rs = make(Releases, 0)
			r = &Release{}
		} else if r.Number() != a.Release {
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
	return
}
