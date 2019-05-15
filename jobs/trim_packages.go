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

// TrimPackagesJobHandler is responsible for removing packages by identifiers
type TrimPackagesJobHandler Job

// NewTrimPackagesJob will return a job suitable for adding to the job processor
func NewTrimPackagesJob(repoID string, maxKeep int) *Job {
	return &Job{
		Type:    TrimPackages,
		SrcRepo: repoID,
		MaxKeep: maxKeep,
	}
}

// NewTrimPackagesJobHandler will create a job handler for the input job and ensure it validates
func NewTrimPackagesJobHandler(j *Job, running bool) (handler *TrimPackagesJobHandler, errs []error) {
	if len(j.SrcRepo) == 0 {
		errs = append(errs, fmt.Errorf("job is missing a source repository"))
	}
	if j.MaxKeep < 1 {
		errs = append(errs, fmt.Errorf("must keep at least one release of a package"))
	}
	if len(errs) == 0 && !running {
		log.Infof("Trim of packages with more than '%d' releases in repo '%s' requested\n", j.MaxKeep, j.SrcRepo)
	}
	h := TrimPackagesJobHandler(*j)
	handler = &h
	return
}

// Execute will attempt removal of excessive packages in the index
func (j *TrimPackagesJobHandler) Execute(_ *JobStore, manager *core.Manager) error {
	if err := manager.TrimPackages(j.SrcRepo, j.MaxKeep); err != nil {
		return err
	}
	log.Goodf("Successfully trimmed packages in repo '%s' with max keep '%d'\n", j.SrcRepo, j.MaxKeep)
	return nil
}

// Describe returns a human readable description for this job
func (j *TrimPackagesJobHandler) Describe() string {
	return fmt.Sprintf("Trim packages to maximum of %d in '%s'", j.MaxKeep, j.SrcRepo)
}
