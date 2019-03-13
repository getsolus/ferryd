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
	"errors"
	"github.com/jmoiron/sqlx"
	"libferry"
)

var (
	// ErrNoJobReady is returned when there are no available jobs or the next job is blocked by a running job
	ErrNoJobReady = errors.New("No jobs ready to run")
)

// JobStore handles the storage and manipulation of incomplete jobs
type JobStore struct {
	db   *sqlx.DB
	next chan *Job
	stop chan bool
	done chan bool
}

// NewStore creates a fully initialized JobStore and sets up Bolt Buckets as needed
func NewStore(path string) (*JobStore, error) {
	// Open the database if we can
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	s := &JobStore{
		db:   db,
		next: make(chan *Job),
		stop: make(chan bool),
		done: make(chan bool),
	}

	return s, nil
}

// Close will clean up our private job database
func (s *JobStore) Close() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}

// UnclaimRunning will find all claimed jobs and unclaim them again
func (s *JobStore) UnclaimRunning() error {
	return nil
}

// Push inserts a new Job into the queue
func (s *JobStore) Push(j *Job) error {
	return nil
}

// Claim gets the first available job, if one exists and is not blocked by running jobs
func (s *JobStore) Claim() (*Job, error) {
	return nil, nil
}

// Retire marks a job as completed and updates the DB record
func (s *JobStore) Retire(j *Job) error {
	return nil
}

// Active will attempt to return a list of active jobs within
// the scheduler suitable for consumption by the CLI client
func (s *JobStore) Active() (libferry.JobSet, error) {
	return nil, nil
}

// Completed will return all successfully completed jobs still stored
func (s *JobStore) Completed() (libferry.JobSet, error) {
	return nil, nil
}

// Failed will return all failed jobs that are still stored
func (s *JobStore) Failed() (libferry.JobSet, error) {
	return nil, nil
}

// ResetCompleted will remove all completion records from our store and reset the pointer
func (s *JobStore) ResetCompleted() error {
	return nil
}

// ResetFailed will remove all fail records from our store and reset the pointer
func (s *JobStore) ResetFailed() error {
	return nil
}
