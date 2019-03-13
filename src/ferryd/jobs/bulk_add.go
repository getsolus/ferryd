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
    "strings"
)

// BulkAddJobHandler is responsible for indexing repositories and should only
// ever be used in sequential queues.
type BulkAddJobHandler Job

// NewBulkAddJob will return a job suitable for adding to the job processor
func NewBulkAddJob(repo string, srcs []string) *Job {
	return &Job{
		Type:    BulkAdd,
		SrcRepo: repo,
		SourcesList: srcs,
        Sources: strings.Join(srcs, ";"),
	}
}

// NewBulkAddJobHandler will create a job handler for the input job and ensure it validates
func NewBulkAddJobHandler(j *Job) (handler *BulkAddJobHandler, err error) {
	if len(j.SrcRepo) == 0 {
		err = fmt.Errorf("job has no repo specified")
		return
	}
    j.SourcesList = strings.Split(j.Sources, ";")
	if len(j.SourcesList) == 0 {
		err = fmt.Errorf("job has no sources specified")
		return
	}
	h := BulkAddJobHandler(*j)
	handler = &h
	return
}

// Execute will attempt the mass-import of packages passed to the job
func (j *BulkAddJobHandler) Execute(_ *JobStore, manager *core.Manager) error {
	if err := manager.AddPackages(j.SrcRepo, j.SourcesList, false); err != nil {
		return err
	}
	log.Infof("Added packages '%v' to repository '%s'\n", j.SourcesList, j.SrcRepo)
	return nil
}

// Describe returns a human readable description for this job
func (j *BulkAddJobHandler) Describe() string {
	return fmt.Sprintf("Add %v packages to repository '%s'", len(j.SourcesList), j.SrcRepo)
}

// IsSerial returns true if a job should not be run alongside other jobs
func (j *BulkAddJobHandler) IsSerial() bool {
	return true
}
