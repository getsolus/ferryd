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

package releases

// ReleaseSet represents a collection of releases for a package in the repos
type ReleaseSet struct {
	Pkg    Release
	Deltas []Release
}

// NewReleaseSet returns an empty ReleaseSet
func NewReleaseSet() *ReleaseSet {
	return &ReleaseSet{
		Deltas: make([]Release, 0),
	}
}

// Len returns the size of the list of Deltas for this release
func (r *ReleaseSet) Len() int {
	return len(r.Deltas)
}

// Less compares the "from_release" number for two deltas
func (r *ReleaseSet) Less(i, j int) bool {
	return r.Deltas[i].From > r.Deltas[j].From
}

// Swap exchanges deltas for sorting
func (r *ReleaseSet) Swap(i, j int) {
	r.Deltas[i], r.Deltas[j] = r.Deltas[j], r.Deltas[i]
}

// Compare finds all of the differences and similarities between two ReleaseSets
func (r *ReleaseSet) Compare(r2 *ReleaseSet) (left, right, same []Release) {
	// Everything in r2 is not in r
	if r == nil {
		right = append(right, r2.Pkg)
		right = append(right, r2.Deltas...)
		return
	}
	// Everything in r is not in r2
	if r2 == nil {
		left = append(left, r.Pkg)
		left = append(left, r.Deltas...)
		return
	}
	// Individually examine releases
	same = append(same, r.Pkg)
	for _, d1 := range r.Deltas {
		found := false
		for _, d2 := range r2.Deltas {
			if found {
				continue
			}
			// In Both R1 and R2
			if d1.From == d2.From {
				same = append(same, d1)
				found = true
			}
		}
		// Found in R1, but not R2
		if !found {
			left = append(left, d1)
		}
	}
	for _, d2 := range r2.Deltas {
		found := false
		for _, d1 := range r.Deltas {
			if found {
				continue
			}
			if d1.From == d2.From {
				found = true
			}
		}
		// Found in R2, but not R1
		if !found {
			right = append(right, d2)
		}
	}
	return
}
