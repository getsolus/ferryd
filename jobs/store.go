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
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/jmoiron/sqlx"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	// ErrNoJobReady is returned when there are no available jobs or the next job is blocked by a running job
	ErrNoJobReady = errors.New("No jobs ready to run")
)

const (
	// JobsDB is the filename of the jobs database
	JobsDB = "jobs.db"
	// SQLiteOpts is a list of options for the go-sqlite3 driver
	SQLiteOpts = "?cache=shared"
)

// JobStore handles the storage and manipulation of incomplete jobs
type JobStore struct {
	db    *sqlx.DB
	next  *Job
	stop  chan bool
	done  chan bool
	wLock sync.Mutex
}

// NewStore creates a fully initialized JobStore and sets up Bolt Buckets as needed
func NewStore(path string) (*JobStore, error) {
	// Open the database if we can
	db, err := sqlx.Open("sqlite3", filepath.Join(path, JobsDB)+SQLiteOpts)
	if err != nil {
		return nil, err
	}
	// See: https://github.com/mattn/go-sqlite3/issues/209
	db.SetMaxOpenConns(1)

	// Create "jobs" table if missing
	db.MustExec(JobSchema)

	s := &JobStore{
		db:   db,
		next: nil,
		stop: make(chan bool),
		done: make(chan bool),
	}
	// reset running jobs and return
	return s, s.UnclaimRunning()
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
	s.wLock.Lock()
	_, err := s.db.Exec(clearRunningJobs)
	if err != nil {
		err = fmt.Errorf("Failed to unclaim running jobs, reason: '%s'", err.Error())
	}
	s.wLock.Unlock()
	return err
}

// Push inserts a new Job into the queue
func (s *JobStore) Push(j *Job) error {
	s.wLock.Lock()
	j.Status = New
	j.Created.Scan(time.Now().UTC())
	_, err := s.db.NamedExec(insertJob, j)
	if err != nil {
		err = fmt.Errorf("Failed to add new job, reason: '%s'", err.Error())
		log.Errorln(err.Error())
	}
	s.wLock.Unlock()
	return err
}

func (s *JobStore) findNewJob() {
	// get the currently runnign jobs
	var active []Job
	err := s.db.Select(&active, runningJobs)
	if err != nil {
		return
	}
	// Check for serial jobs that are blocking
	for _, j := range active {
		if !IsParallel[j.Type] {
			return
		}
	}
	// Otherwise, get the next available job
	var next Job
	err = s.db.Get(&next, nextJob)
	if err != nil {
		return
	}
	// Check if we are blocked by parallel jobs
	if !IsParallel[next.Type] && len(active) > 0 {
		return
	}
	s.next = &next
}

// Claim gets the first available job, if one exists and is not blocked by running jobs
func (s *JobStore) Claim() (j *Job, err error) {
	s.wLock.Lock()
	if s.next == nil {
		err = ErrNoJobReady
		goto UNLOCK
	}
	// claim the next job
	s.next.Status = Running
	s.next.Started.Scan(time.Now().UTC())
	_, err = s.db.NamedExec(markRunning, s.next)
	if err != nil {
		goto UNLOCK
	}
	// find the next replacement job
	j, s.next = s.next, nil
	if j != nil {
		j.SourcesList = strings.Split(j.Sources, ";")
	}
UNLOCK:
	s.findNewJob()
	s.wLock.Unlock()
	return
}

// Retire marks a job as completed and updates the DB record
func (s *JobStore) Retire(j *Job) error {
	s.wLock.Lock()
	j.Finished.Scan(time.Now().UTC())
	_, err := s.db.NamedExec(markFinished, j)
	if err != nil {
		err = fmt.Errorf("Failed to retire job, reason: '%s'", err.Error())
	}
	s.wLock.Unlock()
	return err
}

// Active will attempt to return a list of active jobs within
// the scheduler suitable for consumption by the CLI client
func (s *JobStore) Active() (List, error) {
	var list List
	var list2 List
	err := s.db.Select(&list, newJobs)
	if err != nil {
		err = fmt.Errorf("Failed to read new jobs, reason: '%s'", err.Error())
		log.Errorln(err.Error())
		return nil, err
	}
	err = s.db.Select(&list2, runningJobs)
	if err != nil {
		err = fmt.Errorf("Failed to read active jobs, reason: '%s'", err.Error())
		log.Errorln(err.Error())
		return nil, err
	}
	list = append(list, list2...)
	return list, err
}

// Completed will return all successfully completed jobs still stored
func (s *JobStore) Completed() (List, error) {
	var list List
	err := s.db.Select(&list, completedJobs)
	if err != nil {
		err = fmt.Errorf("Failed to read completed jobs, reason: '%s'", err.Error())
	}
	return list, err
}

// Failed will return all failed jobs that are still stored
func (s *JobStore) Failed() (List, error) {
	var list List
	err := s.db.Select(&list, failedJobs)
	if err != nil {
		err = fmt.Errorf("Failed to read failed jobs, reason: '%s'", err.Error())
	}
	return list, err
}

// ResetCompleted will remove all completion records from our store and reset the pointer
func (s *JobStore) ResetCompleted() error {
	s.wLock.Lock()
	_, err := s.db.Exec(clearCompletedJobs)
	if err != nil {
		err = fmt.Errorf("Failed to clear completed jobs, reason: '%s'", err.Error())
	}
	s.wLock.Unlock()
	return err
}

// ResetFailed will remove all fail records from our store and reset the pointer
func (s *JobStore) ResetFailed() error {
	s.wLock.Lock()
	_, err := s.db.Exec(clearFailedJobs)
	if err != nil {
		err = fmt.Errorf("Failed to clear failed jobs, reason: '%s'", err.Error())
	}
	s.wLock.Unlock()
	return err
}

// ResetQueued will remove all unexecuted records from our store and reset the pointer
func (s *JobStore) ResetQueued() error {
	s.wLock.Lock()
	_, err := s.db.Exec(clearQueuedJobs)
	if err != nil {
		err = fmt.Errorf("Failed to clear queued jobs, reason: '%s'", err.Error())
	}
	s.wLock.Unlock()
	return err
}
