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

// JobSchema is the SQLite3 schema for the Jobs table
const JobSchema = `
CREATE TABLE IF NOT EXISTS jobs (
    id       INTEGER PRIMARY KEY,
    type     INTEGER,
    src      STRING,
    dst      STRING,
    pkg      TEXT,
    max      INTEGER,
    created  DATETIME,
    started  DATETIME,
    finished DATETIME,
    status   INTEGER,
    message  TEXT,
    results  BLOB
)
`

// Queries for retrieving Jobs of a particular status
const (
	newJobs       = "SELECT * FROM jobs WHERE status=0"
	runningJobs   = "SELECT * FROM jobs WHERE status=1"
	failedJobs    = "SELECT * FROM jobs WHERE status=2"
	cancelledJobs = "SELECT * FROM jobs WHERE status=3"
	completedJobs = "SELECT * FROM jobs WHERE status=4"
)

// Query for creating a new Job
const insertJob = `
INSERT INTO jobs (
    id, type,
    src, dst, pkg, max,
    created, started, finished, status, message, results
) VALUES (
    NULL, :type,
    :src, :dst, :pkg, :max,
    :created, NULL, NULL, :status, NULL, NULL
)
`

const (
	getJob  = "SELECT * FROM jobs WHERE id=?"
	nextJob = "SELECT * FROM jobs WHERE status=0 ORDER BY id LIMIT 1"
)

// Queries for updating the status of a job
const (
	markRunning  = "UPDATE jobs SET status=:status, started=:started WHERE id=:id"
	markFinished = "UPDATE jobs SET status=:status, finished=:finished, message=:message WHERE id=:id"
)

// Queries for Cleaning up the Job queue
const (
	clearQueuedJobs    = "DELETE FROM jobs WHERE status=0"
	clearRunningJobs   = "UPDATE jobs SET status=0 WHERE status=1"
	clearFailedJobs    = "DELETE FROM jobs WHERE status=2"
	clearCancelledJobs = "DELETE FROM jobs WHERE status=3"
	clearCompletedJobs = "DELETE FROM jobs WHERE status=4"
)
