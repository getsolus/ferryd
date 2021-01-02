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
	"errors"
	"github.com/getsolus/ferryd/repo/archive"
	"sort"
)

var (
	// ErrInvalidRelease indicates that the relevant Release is malformed
	ErrInvalidRelease = errors.New("invalid release")
)

// Release represents a collection of Archives for a specific Release of a Package in the repo
type Release struct {
	Pkg    *archive.Archive
	Deltas archive.Archives
}

// Number gets the Release number of these archives
func (r Release) Number() int {
	if r.Pkg != nil {
		return r.Pkg.Release
	}
	return r.Deltas[0].Release
}

// Package gets the name of the package this Release belongs to
func (r Release) Package() string {
	if r.Pkg != nil {
		return r.Pkg.Package
	}
	return r.Deltas[0].Package
}

// Sort the internal list of Deltas
func (r *Release) Sort() {
	sort.Sort(r.Deltas)
}

// HasOrphans checks if the Deltas belonging to this release have been orphaned
func (r *Release) HasOrphans() bool {
	return r.Pkg == nil
}

// IsValid checks that a Release has been properly constructed
func (r *Release) IsValid() bool {
	if r.Pkg != nil && !r.Pkg.IsPackage() {
		return false
	}
	for _, delta := range r.Deltas {
		if !delta.IsDelta() {
			return false
		}
		if r.Pkg != nil && delta.Package != r.Pkg.Package {
			return false
		}
		if r.Pkg != nil && delta.Release != r.Pkg.Release {
			return false
		}
	}
	return true
}
