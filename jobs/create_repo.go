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
	"github.com/getsolus/ferryd/repo"
)

// CreateRepoJobHandler is responsible for creating new repositories and should only
// ever be used in sequential queues.
type CreateRepoJobHandler Job

// NewCreateRepoJob will return a job suitable for adding to the job processor
func NewCreateRepoJob(id string) *Job {
	return &Job{
		Type:    CreateRepo,
		DstRepo: id,
	}
}

// NewCreateRepoJobHandler will create a job handler for the input job and ensure it validates
func NewCreateRepoJobHandler(j *Job, running bool) (handler *CreateRepoJobHandler, errs []error) {
	if len(j.DstRepo) == 0 {
		errs = append(errs, fmt.Errorf("job is missing a destination repo"))
	}
	if len(errs) == 0 && !running {
		log.Infof("Creation of repo '%s' requested\n", j.DstRepo)
	}
	h := CreateRepoJobHandler(*j)
	handler = &h
	return
}

// Execute will construct a new repository if possible
func (j *CreateRepoJobHandler) Execute(_ *Store, manager *repo.Manager) error {
	if err := manager.CreateRepo(j.DstRepo); err != nil {
		return err
	}
	log.Goodf("Successfully created repo '%s'\n", j.DstRepo)
	return nil
}

// Describe returns a human readable description for this job
func (j *CreateRepoJobHandler) Describe() string {
	return fmt.Sprintf("Create repository '%s'", j.DstRepo)
}
