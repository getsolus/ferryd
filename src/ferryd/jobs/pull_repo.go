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

	log "github.com/sirupsen/logrus"

	"github.com/getsolus/ferryd/src/ferryd/core"
)

// PullRepoJobHandler is responsible for cloning an existing repository
type PullRepoJobHandler struct {
	sourceID string
	targetID string
}

// NewPullRepoJob will return a job suitable for adding to the job processor
func NewPullRepoJob(sourceID, targetID string) *JobEntry {
	return &JobEntry{
		sequential: true,
		Type:       PullRepo,
		Params:     []string{sourceID, targetID},
	}
}

// NewPullRepoJobHandler will create a job handler for the input job and ensure it validates
func NewPullRepoJobHandler(j *JobEntry) (*PullRepoJobHandler, error) {
	if len(j.Params) != 2 {
		return nil, fmt.Errorf("job has invalid parameters")
	}
	return &PullRepoJobHandler{
		sourceID: j.Params[0],
		targetID: j.Params[1],
	}, nil
}

// Execute will attempt to pull the repos
func (j *PullRepoJobHandler) Execute(jproc *Processor, manager *core.Manager) error {
	changedNames, err := manager.PullRepo(j.sourceID, j.targetID)
	if err != nil {
		log.WithFields(log.Fields{
			"source": j.sourceID,
			"target": j.targetID,
			"error":  err,
		}).Warning("Failed to pull repository")
		return err
	}

	log.WithFields(log.Fields{
		"source": j.sourceID,
		"target": j.targetID,
	}).Info("Pulled repository")

	if !deltasEnabled {
		return nil
	}

	// Create delta job in this repository on the changed names
	// Don't cause indexing because it causes noise
	for _, pkg := range changedNames {
		jproc.PushJob(NewDeltaIndexJob(j.targetID, pkg))
	}

	return nil
}

// Describe returns a human readable description for this job
func (j *PullRepoJobHandler) Describe() string {
	return fmt.Sprintf("Pull repository '%s' into '%s'", j.sourceID, j.targetID)
}
