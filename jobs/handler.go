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
	"fmt"
	"github.com/getsolus/ferryd/core"
)

// A JobHandler is created for each JobEntry, to provide specialised handling
// of the job type
type JobHandler interface {

	// Execute will attempt to execute the given job
	Execute(s *JobStore, m *core.Manager) error

	// Describe will return an appropriate description for the job
	Describe() string
}

// NewJobHandler will return a handler that is loaded only during the execution
// of a previously serialised job
func NewJobHandler(j *Job, running bool) (JobHandler, []error) {
	switch j.Type {
	case BulkAdd:
		return NewBulkAddJobHandler(j, running)
	case CopySource:
		return NewCopySourceJobHandler(j, running)
	case CloneRepo:
		return NewCloneRepoJobHandler(j, running)
	case CreateRepo:
		return NewCreateRepoJobHandler(j, running)
	case DeleteRepo:
		return NewDeleteRepoJobHandler(j, running)
	case Delta:
		return NewDeltaJobHandler(j)
	case DeltaRepo:
		return NewDeltaRepoJobHandler(j, running)
	case DeltaIndex:
		return NewDeltaJobHandler(j)
	case IndexRepo:
		return NewIndexRepoJobHandler(j, running)
	case RemoveSource:
		return NewRemoveSourceJobHandler(j, running)
	case PullRepo:
		return NewPullRepoJobHandler(j, running)
	case TransitProcess:
		return NewTransitJobHandler(j)
	case TrimObsolete:
		return NewTrimObsoleteJobHandler(j, running)
	case TrimPackages:
		return NewTrimPackagesJobHandler(j, running)
	default:
		return nil, []error{fmt.Errorf("unknown job type '%d'", j.Type)}
	}
}
