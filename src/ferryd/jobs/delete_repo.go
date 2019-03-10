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
	log "github.com/sirupsen/logrus"
)

// DeleteRepoJobHandler is responsible for creating new repositories and should only
// ever be used in sequential queues.
type DeleteRepoJobHandler Job

// NewDeleteRepoJob will return a job suitable for adding to the job processor
func NewDeleteRepoJob(id string) *Job {
	return &Job{
		Type:    DeleteRepo,
        SrcRepo: id,
	}
}

// NewDeleteRepoJobHandler will create a job handler for the input job and ensure it validates
func NewDeleteRepoJobHandler(j *Job) (handler *DeleteRepoJobHandler, err error) {
	if len(j.SrcRepo) == 0 {
		err = fmt.Errorf("job is missing a source repo")
        return
	}
	h := DeleteRepoJobHandler(*j)
    handler = &h
    return
}

// Execute will delete an existing repository
func (j *DeleteRepoJobHandler) Execute(_ *Processor, manager *core.Manager) error {
	if err := manager.DeleteRepo(j.SrcRepo); err != nil {
		return err
	}
	log.WithFields(log.Fields{"repo": j.SrcRepo}).Info("Deleted repository")
	return nil
}

// Describe returns a human readable description for this job
func (j *DeleteRepoJobHandler) Describe() string {
	return fmt.Sprintf("Delete repository '%s'", j.SrcRepo)
}

// IsSerial returns true if a job should not be run alongside other jobs
func (J *DeleteRepoJobHandler) IsSerial() bool {
    return true
}
