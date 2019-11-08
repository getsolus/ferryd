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

// JobStatus indicates the status of a Job in the DB
type JobStatus int

const (
	// New indicates a brand-new job
	New JobStatus = 0
	// Running indicates an executing job
	Running = 1
	// Failed indicates a job that finished in failure
	Failed = 2
	// Cancelled indicates a job that was cancelled
	Cancelled = 3
	// Completed indicates a job that successfully finished
	Completed = 4
)

var statusMap = map[JobStatus]string{
	New:       "new",
	Running:   "running",
	Failed:    "failed",
	Cancelled: "cancelled",
	Completed: "completed",
}
