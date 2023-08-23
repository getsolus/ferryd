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

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"

	"github.com/getsolus/ferryd/src/ferryd/core"
	"github.com/getsolus/ferryd/src/libferry"
)

// JobType is a numerical representation of a kind of job
type JobType string

const (

	// BulkAdd is a sequential job which will attempt to add all of the packages
	BulkAdd JobType = "BulkAdd"

	// CopySource is a sequential job to copy from one repo to another
	CopySource = "CopySource"

	// CloneRepo is a sequential job which will attempt to clone a repo
	CloneRepo = "CloneRepo"

	// CreateRepo is a sequential job which will attempt to create a new repo
	CreateRepo = "CreateRepo"

	// DeleteRepo is a sequential job which will attempt to delete a repository
	DeleteRepo = "DeleteRepo"

	// Delta is a parallel job which will attempt the construction of deltas for
	// a given package name + repo
	Delta = "Delta"

	// DeltaIndex is created in response to transit manifest events, and will
	// cause the repository to be reindexed after each delta job continues
	DeltaIndex = "Delta+Index"

	// DeltaRepo is a sequential job which creates Delta jobs for every package in
	// a repo
	DeltaRepo = "DeltaRepo"

	// IndexRepo is a sequential job that requests the repository be re-indexed
	IndexRepo = "IndexRepo"

	// PullRepo is a sequential job that will attempt to pull a repo
	PullRepo = "PullRepo"

	// RemoveSource is a sequential job that will attempt removal of packages
	RemoveSource = "RemoveSource"

	// TransitProcess is a sequential job that will process the incoming uploads
	// directory, dealing with each .tram upload
	TransitProcess = "TransitProcess"

	// TrimObsolete is a sequential job to permanently remove obsolete packages
	// from a repo
	TrimObsolete = "TrimObsolete"

	// TrimPackages is a sequential job to trim fat from a repository
	TrimPackages = "TrimPackages"
)

// A JobHandler is created for each JobEntry, to provide specialised handling
// of the job type
type JobHandler interface {

	// Execute will attempt to execute the given job
	Execute(proc *Processor, m *core.Manager) error

	// Describe will return an appropriate description for the job
	Describe() string
}

// JobEntry is an entry in the JobQueue
type JobEntry struct {
	id         []byte // Unique ID for this job
	sequential bool   // Private to the job implementation
	Type       JobType
	Claimed    bool
	Params     []string
	Timing     libferry.TimingInformation // Store all timing information

	// Not serialised, set by the worker on claim
	description string

	// Not serialised, stored by the worker if the job fails
	failure error
}

// Serialize uses Gob encoding to convert a JobEntry to a byte slice
func (j *JobEntry) Serialize() (result []byte, err error) {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	err = enc.Encode(j)
	if err != nil {
		return
	}
	result = buff.Bytes()
	return
}

// Deserialize use Gob decoding to convert a byte slice to a JobEntry
func Deserialize(serial []byte) (*JobEntry, error) {
	ret := &JobEntry{}
	buff := bytes.NewBuffer(serial)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(ret)
	return ret, err
}

// GetID gets the true numerical ID for this job entry
func (j *JobEntry) GetID() string {
	return fmt.Sprintf("%v", binary.BigEndian.Uint64(j.id))
}

// NewJobHandler will return a handler that is loaded only during the execution
// of a previously serialised job
func NewJobHandler(j *JobEntry) (JobHandler, error) {
	switch j.Type {
	case BulkAdd:
		return NewBulkAddJobHandler(j)
	case CopySource:
		return NewCopySourceJobHandler(j)
	case CloneRepo:
		return NewCloneRepoJobHandler(j)
	case CreateRepo:
		return NewCreateRepoJobHandler(j)
	case DeleteRepo:
		return NewDeleteRepoJobHandler(j)
	case DeltaRepo:
		return NewDeltaRepoJobHandler(j)
	case IndexRepo:
		return NewIndexRepoJobHandler(j)
	case RemoveSource:
		return NewRemoveSourceJobHandler(j)
	case PullRepo:
		return NewPullRepoJobHandler(j)
	case TransitProcess:
		return NewTransitJobHandler(j)
	case TrimObsolete:
		return NewTrimObsoleteJobHandler(j)
	case TrimPackages:
		return NewTrimPackagesJobHandler(j)
	default:
		return nil, fmt.Errorf("unknown job type '%s'", j.Type)
	}
}
