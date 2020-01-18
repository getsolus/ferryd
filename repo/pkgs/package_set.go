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
	"github.com/getsolus/ferryd/repo/releases"
	"github.com/jmoiron/sqlx"
	"sort"
)

// PackageSet contains all of the ReleaseSets for a package in a Repo
type PackageSet struct {
	Name string
	Repo string
	keys []int
	Sets map[int]*releases.ReleaseSet
}

// GetSet retrieves all of the release for a package in a repo and arranges them into a PackageSet
func GetSet(tx *sqlx.Tx, pkg, repo string) (ps *PackageSet, err error) {
	// Create empty PackageSet
	ps = &PackageSet{
		Name: pkg,
		Repo: repo,
		keys: make([]int, 0),
		Sets: make(map[int]*releases.ReleaseSet),
	}
	// Get All Releases
	var rs []releases.Release
	err = tx.Get(&rs, PackageReleases, repo, pkg)
	if err != nil {
		return
	}
	// Sort Releases into Release Sets
	for _, r := range rs {
		if ps.Sets[r.Release] == nil {
			ps.keys = append(ps.keys, r.Release)
			ps.Sets[r.Release] = releases.NewReleaseSet()
		}
		if !r.IsValid() {
			continue
		}
		if r.IsPackage() {
			ps.Sets[r.Release].Pkg = r
		} else if r.IsDelta() {
			ps.Sets[r.Release].Deltas = append(ps.Sets[r.Release].Deltas, r)
		}
	}
	// Sort releases by Release number
	sort.Sort(ps)
	for _, k := range ps.keys {
		sort.Sort(ps.Sets[k])
	}
	return
}

// Len returns the size of the list of releases for this package
func (p *PackageSet) Len() int {
	return len(p.keys)
}

// Less compares the "release" number for two ReleaseSets
func (p *PackageSet) Less(i, j int) bool {
	return p.keys[i] > p.keys[j]
}

// Swap exchanges releases for sorting
func (p *PackageSet) Swap(i, j int) {
	p.keys[i], p.keys[j] = p.keys[j], p.keys[i]
}

// Compare finds all of the differences and similarities between two PackageSets
func (p *PackageSet) Compare(p2 *PackageSet) (left, right, same []releases.Release) {
	// Build a list of keys
	keys := p.keys
	for _, k2 := range p2.keys {
		found := false
		for _, k := range p.keys {
			if k == k2 {
				found = true
			}
		}
		if !found {
			keys = append(keys, k2)
		}
	}
	// Sort highest to lowest
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	// Compare all ReleaseSets
	for _, k := range keys {
		r1 := p.Sets[k]
		r2 := p2.Sets[k]
		l, r, s := r1.Compare(r2)
		left = append(left, l...)
		right = append(right, r...)
		same = append(same, s...)
	}
	return
}
