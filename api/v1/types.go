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

package v1

import (
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"time"
)

// GenericResponse is a catchall for API responses
type GenericResponse struct {
	// Errors is a list of errors encountered in this transaction
	Errors []string
}

// StatusResponse is a response from the 'status' endpoint
type StatusResponse struct {
	GenericResponse
	// TimeStarted is the time when the daemon was last started
	TimeStarted time.Time
	// Version is the version of the daemon
	Version string
	// CurrentJobs is a list of running and queued jobs
	CurrentJobs jobs.List
	// FailedJobs is a list of failed jobs
	FailedJobs jobs.List
	// CompletedJobs is a list of completed jobs
	CompletedJobs jobs.List
}

// Uptime will determine the uptime of the daemon
func (s StatusResponse) Uptime() time.Duration {
	return time.Now().UTC().Sub(s.TimeStarted)
}

// RepoList is a response from the 'repos' endpoint
type RepoList struct {
	GenericResponse
	// Repos is list of all the repos
	Repos []repo.Repo
}

// A PoolItem simply has an ID and a refcount, allowing us to examine our
// local storage efficiency.
type PoolItem struct {
	ID       string `json:"id"`
	RefCount int    `json:"refCount"`
}

// A PoolResponse is a listing of all pool items
type PoolResponse struct {
	GenericResponse
	Items []PoolItem
}
