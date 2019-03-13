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
	log "github.com/DataDrake/waterlog"
)

// RemoveSourceJobHandler is responsible for removing packages by identifiers
type RemoveSourceJobHandler Job

// NewRemoveSourceJob will return a job suitable for adding to the job processor
func NewRemoveSourceJob(repoID, source string, release int) *Job {
	return &Job{
		Type:    RemoveSource,
		SrcRepo: repoID,
		Sources: []string{source},
		Release: release,
	}
}

// NewRemoveSourceJobHandler will create a job handler for the input job and ensure it validates
func NewRemoveSourceJobHandler(j *Job) (handler *RemoveSourceJobHandler, err error) {
	if len(j.SrcRepo) == 0 {
		err = fmt.Errorf("job is missing source repo")
		return
	}
	if len(j.Sources) == 0 {
		err = fmt.Errorf("job is missing source package")
		return
	}
	if j.Release == 0 {
		err = fmt.Errorf("job has invalid release number: 0")
		return
	}
	h := RemoveSourceJobHandler(*j)
	handler = &h
	return
}

// Execute will remove the source&rel match from the repo
func (j *RemoveSourceJobHandler) Execute(_ *Processor, manager *core.Manager) error {
	if err := manager.RemoveSource(j.SrcRepo, j.Sources[0], j.Release); err != nil {
		return err
	}
	log.Goodf("Successfully removed release '%d' of source '%s' from repo '%s'\n", j.Release, j.Sources[0], j.SrcRepo)
	return nil
}

// Describe returns a human readable description for this job
func (j *RemoveSourceJobHandler) Describe() string {
	return fmt.Sprintf("Remove sources by id '%s' (rel: %d) in '%s'", j.Sources[0], j.Release, j.SrcRepo)
}

// IsSerial returns true if a job should not be run alongside other jobs
func (J *RemoveSourceJobHandler) IsSerial() bool {
	return true
}
