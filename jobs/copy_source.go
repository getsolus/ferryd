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

// CopySourceJobHandler is responsible for removing packages by identifiers
type CopySourceJobHandler Job

// NewCopySourceJob will return a job suitable for adding to the job processor
func NewCopySourceJob(srcRepo, dstRepo, source string, release int) *Job {
	return &Job{
		Type:    CopySource,
		SrcRepo: srcRepo,
		DstRepo: dstRepo,
		Sources: source,
		Release: release,
	}
}

// NewCopySourceJobHandler will create a job handler for the input job and ensure it validates
func NewCopySourceJobHandler(j *Job, running bool) (handler *CopySourceJobHandler, errs []error) {
	if len(j.SrcRepo) == 0 {
		errs = append(errs, fmt.Errorf("Job '%d' is missing source repo", j.ID))
	}
	if len(j.DstRepo) == 0 {
		errs = append(errs, fmt.Errorf("Job '%d' is missing destination repo", j.ID))
	}
	if len(j.Sources) == 0 {
		errs = append(errs, fmt.Errorf("Kob '%d' is missing source name", j.ID))
	}
	if j.Release == 0 || j.Release < -1 {
		errs = append(errs, fmt.Errorf("job '%d' has invalid release number: '%d'", j.ID, j.Release))
	}
	if len(errs) == 0 && !running {
		log.Info("Copy of release '%d' of source '%s' from repo '%s' to '%s' requested\n", j.Release, j.Sources, j.SrcRepo, j.DstRepo)
	}
	h := CopySourceJobHandler(*j)
	handler = &h
	return
}

// Execute will copy the source&rel match from the repo to the target
func (j *CopySourceJobHandler) Execute(_ *JobStore, manager *core.Manager) error {
	if err := manager.CopySource(j.SrcRepo, j.DstRepo, j.Sources, j.Release); err != nil {
		return err
	}
	log.Goodf("Successfully copied release '%d' of source '%s' from repo '%s' to '%s'\n", j.Release, j.Sources, j.SrcRepo, j.DstRepo)
	return nil
}

// Describe returns a human readable description for this job
func (j *CopySourceJobHandler) Describe() string {
	return fmt.Sprintf("Copy sources by id '%s' (rel: %d) in '%s' to '%s'", j.Sources, j.Release, j.SrcRepo, j.DstRepo)
}
