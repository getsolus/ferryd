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

package api

import (
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"strings"
)

type JobRequest struct {
	Type string
	// Type-specific arguments
	SrcRepo string
	DstRepo string
	Sources []string
	Release int
	MaxKeep int
	Mode    int
}

func (j *JobRequest) Convert() (job *jobs.Job, errs []error) {
	// Get the JobType
	t := jobs.Mapping[j.Type]
	if t == jobs.Invalid {
		errs = append(errs, fmt.Errorf("Invalid job type '%s'", j.Type))
		return
	}
	// Create the Job instance
	job = &jobs.Job{
		Type:        t,
		SrcRepo:     j.SrcRepo,
		DstRepo:     j.DstRepo,
		Sources:     strings.Join(j.Sources, ";"),
		SourcesList: j.Sources,
		Release:     j.Release,
		MaxKeep:     j.MaxKeep,
		Mode:        j.Mode,
	}
	// Indirectly Validate Job
	_, errs = jobs.NewJobHandler(job, false)
	return
}
