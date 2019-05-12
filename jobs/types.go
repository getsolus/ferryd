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

	// BulkAdd is a sequential job which will attempt to add all of the packages
	BulkAdd JobType = 0

	// CopySource is a sequential job to copy from one repo to another
	CopySource = 1

	// CloneRepo is a sequential job which will attempt to clone a repo
	CloneRepo = 2

	// CreateRepo is a sequential job which will attempt to create a new repo
	CreateRepo = 3

	// DeleteRepo is a sequential job which will attempt to delete a repository
	DeleteRepo = 4

	// Delta is a parallel job which will attempt the construction of deltas for
	// a given package name + repo
	Delta = 5

	// DeltaIndex is created in response to transit manifest events, and will
	// cause the repository to be reindexed after each delta job continues
	DeltaIndex = 6

	// DeltaRepo is a sequential job which creates Delta jobs for every package in
	// a repo
	DeltaRepo = 7

	// IndexRepo is a sequential job that requests the repository be re-indexed
	IndexRepo = 8

	// PullRepo is a sequential job that will attempt to pull a repo
	PullRepo = 9

	// RemoveSource is a sequential job that will attempt removal of packages
	RemoveSource = 10

	// TransitProcess is a sequential job that will process the incoming uploads
	// directory, dealing with each .tram upload
	TransitProcess = 11

	// TrimObsolete is a sequential job to permanently remove obsolete packages
	// from a repo
	TrimObsolete = 12

	// TrimPackages is a sequential job to trim fat from a repository
	TrimPackages = 13
)

// IsParallel tells us if a particular JobType can be done in parallel
var IsParallel = map[JobType]bool{
	CopySource: true,
	Delta:      true,
	DeltaIndex: true,
}
