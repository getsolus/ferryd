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

package jobs

// JobType is a numerical representation of a kind of job
type JobType int

const (
	// Invalid is the Zero for JobType
	Invalid JobType = 0
	// Check compares the contents of a repo on Disk with the DB
	Check = 1
	// CherryPick syncs a single package from one repo to another
	CherryPick = 2
	// Clone copies an existing repo into a new one
	Clone = 3
	// Compare creates a diff of the contents of two repos
	Compare = 4
	// Create adds a new empty repo
	Create = 5
	// Delta generates missing Delta Packages for an entire repo
	Delta = 6
	// DeltaPackage generates Deltas for a single package
	DeltaPackage = 7
	// Import adds a new repo to the DB from an existing filepath
	Import = 8
	// Index regenerates the index for a repo
	Index = 9
	// Remove removes a repo from the DB but not its contents on disk
	Remove = 10
	// Rescan updates the DB with the contents of a repo on disk
	Rescan = 11
	// Sync replicates the contents of one repo into another
	Sync = 12
	// TrimObsoletes removes obsoleted packages from the repo
	TrimObsoletes = 13
	// TrimPackages removes old release of packages from the repo
	TrimPackages = 14
	// TransitPackage addss a new package to the Pool and any auto-transit repos
	TransitPackage = 15
)

var typeMap = map[JobType]string{
	Invalid:        "INVALID",
	Check:          "Check",
	CherryPick:     "Cherry-Pick",
	Clone:          "Clone",
	Compare:        "Compare",
	Create:         "Create",
	Delta:          "Delta",
	DeltaPackage:   "Delta Package",
	Import:         "Import",
	Index:          "Index",
	Remove:         "Remove",
	Rescan:         "Rescan",
	Sync:           "Scan",
	TrimObsoletes:  "Trim Obsoletes",
	TrimPackages:   "Trim Packages",
	TransitPackage: "Transit Package",
}

// IsParallel tells us if a particular JobType can be done in parallel
var IsParallel = map[JobType]bool{
	Check:        true,
	Compare:      true,
	DeltaPackage: true,
}
