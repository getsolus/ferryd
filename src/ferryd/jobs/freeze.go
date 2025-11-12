//
// Copyright © 2025 Solus Project
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

	log "github.com/sirupsen/logrus"

	"github.com/getsolus/ferryd/src/ferryd/core"
)

// FreezeRepoJobHandler is responsible for freezing an existing repository
type FreezeRepoJobHandler struct {
	repoID string
}

func NewFreezeRepoJob(repoID string) *JobEntry {
	return &JobEntry{
		sequential: true,
		Type:       FreezeRepo,
		Params:     []string{repoID},
	}
}

// NewFreezeRepoJobHandler will create a job handler for the input job and ensure it validates
func NewFreezeRepoJobHandler(j *JobEntry) (*FreezeRepoJobHandler, error) {
	if len(j.Params) != 1 {
		return nil, fmt.Errorf("job has invalid parameters")
	}
	return &FreezeRepoJobHandler{repoID: j.Params[0]}, nil
}

// Execute will attempt to Freeze the repos
func (j *FreezeRepoJobHandler) Execute(jproc *Processor, manager *core.Manager) error {
	if err := manager.FreezeRepo(j.repoID); err != nil {
		log.WithFields(log.Fields{"repo": j.repoID, "error": err}).
			Warning("Failed to freeze repository")
		return err
	}

	return nil
}

// Describe returns a human readable description for this job
func (j *FreezeRepoJobHandler) Describe() string {
	return fmt.Sprintf("Freeze repository '%s'", j.repoID)
}

// UnfreezeRepoJobHandler is responsible for unfreezing an existing repository
type UnfreezeRepoJobHandler struct {
	repoID string
}

func NewUnfreezeRepoJob(repoID string) *JobEntry {
	return &JobEntry{
		sequential: true,
		Type:       UnfreezeRepo,
		Params:     []string{repoID},
	}
}

// NewUnfreezeRepoJobHandler will create a job handler for the input job and ensure it validates
func NewUnfreezeRepoJobHandler(j *JobEntry) (*UnfreezeRepoJobHandler, error) {
	if len(j.Params) != 1 {
		return nil, fmt.Errorf("job has invalid parameters")
	}
	return &UnfreezeRepoJobHandler{repoID: j.Params[0]}, nil
}

// Execute will attempt to Freeze the repos
func (j *UnfreezeRepoJobHandler) Execute(jproc *Processor, manager *core.Manager) error {
	if err := manager.UnfreezeRepo(j.repoID); err != nil {
		log.WithFields(log.Fields{"repo": j.repoID, "error": err}).
			Warning("Failed to unfreeze repository")
		return err
	}

	return nil
}

// Describe returns a human readable description for this job
func (j *UnfreezeRepoJobHandler) Describe() string {
	return fmt.Sprintf("Unfreeze repository '%s'", j.repoID)
}
