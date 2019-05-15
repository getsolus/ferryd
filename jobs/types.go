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
	// BulkAdd is a sequential job which will attempt to add all of the packages
	BulkAdd = 1

	// CopySource is a sequential job to copy from one repo to another
	CopySource = 2

	// CloneRepo is a sequential job which will attempt to clone a repo
	CloneRepo = 3

	// CreateRepo is a sequential job which will attempt to create a new repo
	CreateRepo = 4

	// DeleteRepo is a sequential job which will attempt to delete a repository
	DeleteRepo = 5

	// Delta is a parallel job which will attempt the construction of deltas for
	// a given package name + repo
	Delta = 6

	// DeltaIndex is created in response to transit manifest events, and will
	// cause the repository to be reindexed after each delta job continues
	DeltaIndex = 7

	// DeltaRepo is a sequential job which creates Delta jobs for every package in
	// a repo
	DeltaRepo = 8

	// IndexRepo is a sequential job that requests the repository be re-indexed
	IndexRepo = 9

	// PullRepo is a sequential job that will attempt to pull a repo
	PullRepo = 10

	// RemoveSource is a sequential job that will attempt removal of packages
	RemoveSource = 11

	// TransitProcess is a sequential job that will process the incoming uploads
	// directory, dealing with each .tram upload
	TransitProcess = 12

	// TrimObsolete is a sequential job to permanently remove obsolete packages
	// from a repo
	TrimObsolete = 13

	// TrimPackages is a sequential job to trim fat from a repository
	TrimPackages = 14
)

var Mapping = map[string]JobType{
	"BulkAdd":        BulkAdd,
	"CopySource":     CopySource,
	"CloneRepo":      CloneRepo,
	"CreateRepo":     CreateRepo,
	"DeleteRepo":     DeleteRepo,
	"Delta":          Delta,
	"DeltaRepo":      DeltaRepo,
	"IndexRepo":      IndexRepo,
	"PullRepo":       PullRepo,
	"RemoveSource":   RemoveSource,
	"TransitProcess": TransitProcess,
	"TrimObsolete":   TrimObsolete,
	"TrimPackages":   TrimPackages,
}

// IsParallel tells us if a particular JobType can be done in parallel
var IsParallel = map[JobType]bool{
	CopySource: true,
	Delta:      true,
	DeltaIndex: true,
}
