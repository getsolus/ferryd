//
// Copyright © 2017-2020 Solus Project
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
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/ferryd/config"
	"github.com/jmoiron/sqlx"
	"path/filepath"
	"sync"
	"time"
)

var (
	// ErrNoJobReady is returned when there are no available jobs or the next job is blocked by a running job
	ErrNoJobReady = errors.New("No jobs ready to run")
)

const (
	// DB is the filename of the jobs database
	DB = "jobs.db"
	// SQLiteOpts is a list of options for the go-sqlite3 driver
	SQLiteOpts = "?cache=shared"
)

// Store handles the storage and manipulation of incomplete jobs
type Store struct {
	sync.Mutex

	db   *sqlx.DB
	next *Job
	stop chan bool
	done chan bool
}

// NewStore creates a fully initialized Store and sets up Bolt Buckets as needed
func NewStore() (s *Store, err error) {
	// Open the database if we can
	db, err := sqlx.Open("sqlite3", filepath.Join(config.Current.BaseDir, DB)+SQLiteOpts)
	if err != nil {
		return nil, err
	}
	// See: https://github.com/mattn/go-sqlite3/issues/209
	db.SetMaxOpenConns(1)
	// Create "jobs" table if missing
	db.MustExec(JobSchema)
	s = &Store{
		db:   db,
		next: nil,
		stop: make(chan bool),
		done: make(chan bool),
	}
	// reset running jobs and return
	err = s.UnclaimRunning()
	return
}

// Close will clean up our private job database
func (s *Store) Close() error {
	if s.db != nil {
		err := s.db.Close()
		s.db = nil
		return err
	}
	return nil
}

// GetJob retrieves a Job from the DB
func (s *Store) GetJob(id int) (j *Job, err error) {
	j = &Job{}
	err = s.db.Get(&j, getJob, id)
	return
}

// UnclaimRunning will find all claimed jobs and unclaim them again
func (s *Store) UnclaimRunning() (err error) {
	s.Lock()
	if _, err = s.db.Exec(clearRunningJobs); err != nil {
		err = fmt.Errorf("Failed to unclaim running jobs, reason: '%s'", err.Error())
	}
	s.Unlock()
	return
}

// Push inserts a new Job into the queue
func (s *Store) Push(j *Job) (id int, err error) {
	s.Lock()
	// Set Job parameters
	j.Status = New
	j.Created.Time = time.Now().UTC()
	j.Created.Valid = true
	// Start a DB transaction
	tx, err := s.db.Beginx()
	if err != nil {
		goto UNLOCK
	}
	// Insert the New Job
	if err = j.Create(tx); err != nil {
		tx.Rollback()
		goto UNLOCK
	}
	id = j.ID
	// Complete the transaction
	err = tx.Commit()
UNLOCK:
	s.Unlock()
	return id, err
}

func (s *Store) findNewJob() {
	var next Job
	if err := s.db.Get(&next, nextJob); err != nil {
		return
	}
	s.next = &next
}

// Claim gets the first available job, if one exists and is not blocked by running jobs
func (s *Store) Claim() (j *Job, err error) {
	var tx *sqlx.Tx
	s.Lock()
	if s.next == nil {
		err = ErrNoJobReady
		goto UNLOCK
	}
	// claim the next job
	s.next.Status = Running
	s.next.Started.Time = time.Now().UTC()
	s.next.Started.Valid = true
	// Start a DB transaction
	tx, err = s.db.Beginx()
	if err != nil {
		goto UNLOCK
	}
	// Save the status change
	if err = s.next.Save(tx); err != nil {
		tx.Rollback()
		goto UNLOCK
	}
	// Finish the transaction
	if err = tx.Commit(); err != nil {
		goto UNLOCK
	}
	// find the next replacement job
	j, s.next = s.next, nil
UNLOCK:
	s.findNewJob()
	s.Unlock()
	return
}

// Retire marks a job as completed and updates the DB record
func (s *Store) Retire(j *Job) error {
	s.Lock()
	// Start a DB transaction
	tx, err := s.db.Beginx()
	if err != nil {
		goto UNLOCK
	}
	// Mark as finished
	j.Finished.Time = time.Now().UTC()
	j.Finished.Valid = true
	if err = j.Save(tx); err != nil {
		tx.Rollback()
		goto UNLOCK
	}
	// Finish the transaction
	err = tx.Commit()
UNLOCK:
	s.Unlock()
	return err
}

// Active will attempt to return a list of active jobs within
// the scheduler suitable for consumption by the CLI client
func (s *Store) Active() (list List, err error) {
	var list2 List
	// Get all new jobs
	if err = s.db.Select(&list, newJobs); err != nil {
		log.Errorf("Failed to read new jobs, reason: '%s'", err.Error())
		return
	}
	// Get All running jobs
	if err = s.db.Select(&list2, runningJobs); err != nil {
		log.Errorf("Failed to read active jobs, reason: '%s'", err.Error())
		return
	}
	// Append them together
	list = append(list, list2...)
	return
}

// Completed will return all successfully completed jobs still stored
func (s *Store) Completed() (list List, err error) {
	if err = s.db.Select(&list, completedJobs); err != nil {
		err = fmt.Errorf("Failed to read completed jobs, reason: '%s'", err.Error())
	}
	return
}

// Failed will return all failed jobs that are still stored
func (s *Store) Failed() (list List, err error) {
	if err = s.db.Select(&list, failedJobs); err != nil {
		err = fmt.Errorf("Failed to read failed jobs, reason: '%s'", err.Error())
	}
	return
}

// ResetCompleted will remove all completion records from our store and reset the pointer
func (s *Store) ResetCompleted() (err error) {
	s.Lock()
	if _, err = s.db.Exec(clearCompletedJobs); err != nil {
		err = fmt.Errorf("Failed to clear completed jobs, reason: '%s'", err.Error())
	}
	s.Unlock()
	return
}

// ResetFailed will remove all fail records from our store and reset the pointer
func (s *Store) ResetFailed() (err error) {
	s.Lock()
	if _, err = s.db.Exec(clearFailedJobs); err != nil {
		err = fmt.Errorf("Failed to clear failed jobs, reason: '%s'", err.Error())
	}
	s.Unlock()
	return
}

// ResetQueued will remove all unexecuted records from our store and reset the pointer
func (s *Store) ResetQueued() (err error) {
	s.Lock()
	if _, err = s.db.Exec(clearQueuedJobs); err != nil {
		err = fmt.Errorf("Failed to clear queued jobs, reason: '%s'", err.Error())
	}
	s.Unlock()
	return
}
