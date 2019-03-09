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

const (
    CloneTip  = 0
    CloneFull = 1
)

// CloneRepoJobHandler is responsible for cloning an existing repository
type CloneRepoJobHandler Job

// NewCloneRepoJob will return a job suitable for adding to the job processor
func NewCloneRepoJob(srcRepo, newRepo string, cloneAll bool) *Job {
	mode := CloneTip
	if cloneAll {
		mode = CloneFull
	}
	return &Job{
		Type:    CloneRepo,
        SrcRepo: srcRepo,
        DstRepo: newRepo,
        Mode:    mode,
	}
}

// NewCloneRepoJobHandler will create a job handler for the input job and ensure it validates
func NewCloneRepoJobHandler(j *Job) (handler *CloneRepoJobHandler, err error) {
	if len(j.SrcRepo) == 0 {
		err = fmt.Errorf("job has no source repo")
        return
	}
	if len(j.DstRepo) == 0 {
		err = fmt.Errorf("job has no destination repo")
        return
	}
	if j.Mode < CloneTip || j.Mode > CloneFull {
		err = fmt.Errorf("job has invalid mode: %d", j.Mode)
        return
	}
	h := CloneRepoJobHandler(*j)
    handler = &h
    return
}

// Execute attempt to clone the repoID to newClone, optionally at full depth
func (j *CloneRepoJobHandler) Execute(_ *Processor, manager *core.Manager) error {
	fullClone := false
	if j.Mode == CloneFull {
		fullClone = true
	}

	if err := manager.CloneRepo(j.SrcRepo, j.DstRepo, fullClone); err != nil {
		return err
	}
	log.WithFields(log.Fields{
        "srcRepo": j.SrcRepo,
        "dstRepo": j.DstRepo,
    }).Info("Cloned repository %s into %s")
	return nil
}

// Describe returns a human readable description for this job
func (j *CloneRepoJobHandler) Describe() string {
	return fmt.Sprintf("Clone repository '%s' into '%s'", j.SrcRepo, j.DstRepo)
}

// IsSerial returns true if a job should not be run alongside other jobs
func (J *CloneRepoJobHandler) IsSerial() bool {
    return true
}
