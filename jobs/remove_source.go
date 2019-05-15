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
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/ferryd/core"
)

// RemoveSourceJobHandler is responsible for removing packages by identifiers
type RemoveSourceJobHandler Job

// NewRemoveSourceJob will return a job suitable for adding to the job processor
func NewRemoveSourceJob(repoID, source string, release int) *Job {
	return &Job{
		Type:    RemoveSource,
		SrcRepo: repoID,
		Sources: source,
		Release: release,
	}
}

// NewRemoveSourceJobHandler will create a job handler for the input job and ensure it validates
func NewRemoveSourceJobHandler(j *Job, running bool) (handler *RemoveSourceJobHandler, errs []error) {
	if len(j.SrcRepo) == 0 {
		errs = append(errs, fmt.Errorf("job is missing source repo"))
	}
	if len(j.Sources) == 0 {
		errs = append(errs, fmt.Errorf("job is missing source package"))
	}
	if j.Release == 0 {
		errs = append(errs, fmt.Errorf("job has invalid release number: 0"))
	}
	if len(errs) == 0 && !running {
		log.Infof("Removal of release '%d' of source '%s' in repo '%s' requested", j.Release, j.SrcRepo, j.DstRepo)
	}
	h := RemoveSourceJobHandler(*j)
	handler = &h
	return
}

// Execute will remove the source&rel match from the repo
func (j *RemoveSourceJobHandler) Execute(_ *JobStore, manager *core.Manager) error {
	if err := manager.RemoveSource(j.SrcRepo, j.Sources, j.Release); err != nil {
		return err
	}
	log.Goodf("Successfully removed release '%d' of source '%s' from repo '%s'\n", j.Release, j.Sources, j.SrcRepo)
	return nil
}

// Describe returns a human readable description for this job
func (j *RemoveSourceJobHandler) Describe() string {
	return fmt.Sprintf("Remove sources by id '%s' (rel: %d) in '%s'", j.Sources, j.Release, j.SrcRepo)
}
