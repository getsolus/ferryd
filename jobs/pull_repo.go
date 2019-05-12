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

// PullRepoJobHandler is responsible for cloning an existing repository
type PullRepoJobHandler Job

// NewPullRepoJob will return a job suitable for adding to the job processor
func NewPullRepoJob(sourceID, targetID string) *Job {
	return &Job{
		Type:    PullRepo,
		SrcRepo: sourceID,
		DstRepo: targetID,
	}
}

// NewPullRepoJobHandler will create a job handler for the input job and ensure it validates
func NewPullRepoJobHandler(j *Job) (handler *PullRepoJobHandler, err error) {
	if len(j.SrcRepo) == 0 {
		err = fmt.Errorf("job is missing source repo")
		return
	}
	if len(j.DstRepo) == 0 {
		err = fmt.Errorf("job is missing destination repo")
		return
	}
	h := PullRepoJobHandler(*j)
	handler = &h
	return
}

// Execute will attempt to pull the repos
func (j *PullRepoJobHandler) Execute(s *JobStore, manager *core.Manager) error {
	changedNames, err := manager.PullRepo(j.SrcRepo, j.DstRepo)
	if err != nil {
		log.Warnf("Failed to pull repo '%s' into '%s', reason: '%s'\n", j.SrcRepo, j.DstRepo, err.Error())
		return err
	}

	log.Goodf("Successfully pulled repo '%s' into '%s'\n", j.SrcRepo, j.DstRepo)

	// Create delta job in this repository on the changed names
	// Don't cause indexing because it causes noise
	for _, pkg := range changedNames {
		s.Push(NewDeltaIndexJob(j.DstRepo, pkg))
	}

	return nil
}

// Describe returns a human readable description for this job
func (j *PullRepoJobHandler) Describe() string {
	return fmt.Sprintf("Pull repository '%s' into '%s'", j.SrcRepo, j.DstRepo)
}
