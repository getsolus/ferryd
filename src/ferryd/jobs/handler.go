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
	"ferryd/core"
	"fmt"
)

// A JobHandler is created for each JobEntry, to provide specialised handling
// of the job type
type JobHandler interface {

	// Execute will attempt to execute the given job
	Execute(proc *Processor, m *core.Manager) error

	// Describe will return an appropriate description for the job
	Describe() string

    // IsSerial returns true if a job should not be run alongside other jobs
    IsSerial() bool
}

// NewJobHandler will return a handler that is loaded only during the execution
// of a previously serialised job
func NewJobHandler(j *Job) (JobHandler, error) {
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
	case Delta:
		return NewDeltaJobHandler(j, false)
	case DeltaRepo:
		return NewDeltaRepoJobHandler(j)
	case DeltaIndex:
		return NewDeltaJobHandler(j, true)
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
