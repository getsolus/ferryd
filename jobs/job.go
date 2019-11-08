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
	"database/sql"
	"database/sql/driver"
	"time"
)

// JobSchema is the SQLite3 schema for the Jobs table
const JobSchema = `
CREATE TABLE IF NOT EXISTS jobs (
    id         INTEGER PRIMARY KEY,
    type       INTEGER,
    src_repo   STRING,
    dst_repo   STRING,
    sources    TEXT,
    release    INTEGER,
    max_keep   INTEGER,
    mode       INTEGER,
    created    DATETIME,
    started    DATETIME,
    finished   DATETIME,
    status     INTEGER,
    message    TEXT
)
`

// Job Status
const (
	New       = 0
	Running   = 1
	Failed    = 2
	Cancelled = 3
	Completed = 4
)

// NullTime allows for NULL timestamps
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// Job is an entry in the Job Table
type Job struct {
	ID   int
	Type JobType
	// Job-specific arguments
	SrcRepo     string `db:"src_repo"`
	DstRepo     string `db:"dst_repo"`
	Sources     string `db:"sources"`
	SourcesList []string
	Release     int
	MaxKeep     int `db:"max_keep"`
	Mode        int
	// Job tracking
	Created  NullTime
	Started  NullTime
	Finished NullTime
	Status   int
	Message  sql.NullString
}

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
    src_repo, dst_repo, sources, release, max_keep, mode,
    created, started, finished, status, message
) VALUES (
    NULL, :type,
    :src_repo, :dst_repo, :sources, :release, :max_keep, :mode,
    :created, NULL, NULL, :status, NULL
)
`

const (
	getJob = "SELECT * FROM jobs WHERE id=:id"
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

// QueuedSince will let us know how long this task has been queued
func (j *Job) QueuedSince() time.Duration {
	return time.Now().UTC().Sub(j.Created.Time)
}

// Executed will work out how long ago it stopped executing
func (j *Job) Executed() time.Duration {
	return time.Now().UTC().Sub(j.Finished.Time)
}

// ExecutionTime will return the time it took to execute a job
func (j *Job) ExecutionTime() time.Duration {
	return j.Finished.Time.Sub(j.Started.Time)
}

// QueuedTime will return the total time that the job was queued for
func (j *Job) QueuedTime() time.Duration {
	return j.Started.Time.Sub(j.Created.Time)
}

// TotalTime will return the total time a job took to complete from queuing
func (j *Job) TotalTime() time.Duration {
	return j.QueuedTime() + j.ExecutionTime()
}

// List is a list of Jobs
type List []*Job

// Len returns the length of the list
func (l List) Len() int {
	return len(l)
}

// Less compares tow jobs by creation time
func (l List) Less(i, j int) bool {
	return l[i].Created.Time.Before(l[j].Created.Time)
}

// Swap switches two Jobs for sorting
func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
