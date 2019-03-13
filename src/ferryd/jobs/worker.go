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
	"ferryd/core"
	log "github.com/DataDrake/waterlog"
	"math/rand"
	"time"
)

// MinWait is the minimum amount of time between retries for a worker
const MinWait = time.Second * 2

// MaxJitter sets the upper limit on the random jitter used for retry times
const MaxJitter int64 = 512

// A Worker is used to execute some portion of the incoming workload, and will
// keep polling for the correct job type to process
type Worker struct {
	timer   *time.Timer
	manager *core.Manager
	store   *JobStore
	stop    chan bool
	done    chan bool
}

// NewWorker creates a new worker for the pool
func NewWorker(store *JobStore, manager *core.Manager) *Worker {
	if store == nil {
		panic("Constructed a Worker without a valid JobStore!")
	}
	if manager == nil {
		panic("Constructed a Worker without a valid Manager!")
	}

	return &Worker{
		timer:   nil, // Init this when we start up
		manager: manager,
		store:   store,
		stop:    make(chan bool),
		done:    make(chan bool),
	}
}

// Stop will demand that all new requests are no longer processed
func (w *Worker) Stop() {
	if w.timer == nil {
		return
	}
	w.stop <- true
	w.timer.Stop()
	<-w.done
}

// Start will begin the main execution of this worker, and will continuously
// poll for new jobs with an increasing increment (with a ceiling limit)
func (w *Worker) Start() {
	// Let's get our timer initialised
	w.setTime()

	for {
		select {
		case <-w.stop:
			// Bail now, we've been told to go home
			return

		case <-w.timer.C:
			// Try to grab a job
			job, err := w.store.Claim()

			// Report the error
			if err != nil {
				if err != ErrNoJobReady {
					log.Errorf("Failed to grab a work queue item, reason: '%s'\n", err.Error())
				}
				w.setTime()
				continue
			}

			// Got a job, now process it
			w.processJob(job)

			// Now we mark end time so we can calculate how long it took
			job.Finished = time.Now().UTC()

			// Mark the job as dealt with
			err = w.store.Retire(job)

			// Report failure in retiring the job
			if err != nil {
				log.Error("Error in retiring job '%v' of type '%v', reason: '%s'\n", job.ID, job.Type, err.Error())
			}

			// We had a job, so we must reset the timeout period
			w.setTime()
		}
	}
}

// setTime will update the timer resetting it to MinWait + some random jitter to help with contention
func (w *Worker) setTime() {
	delay := MinWait + (time.Millisecond * time.Duration(rand.Int63n(MaxJitter)))
	if w.timer == nil {
		w.timer = time.NewTimer(delay)
	} else {
		w.timer.Reset(delay)
	}
}

// processJob will actually examine the given job and figure out how
// to execute it. Each Worker can only execute a single job at a time
func (w *Worker) processJob(job *Job) {
	handler, err := NewJobHandler(job)

	if err != nil {
		job.Status = Failed
		job.Message = err.Error()
		log.Errorf("No known job handler, cannot continue with job '%v'\n", job.ID)
		return
	}

	// Safely have a handler now
	job.Message = handler.Describe()

	// Try to execute it, report the error
	if err := handler.Execute(w.store, w.manager); err != nil {
		job.Status = Failed
		job.Message = err.Error()
		log.Error("Job '%v' failed with error: '%s'\n", job.ID, err.Error())
		return
	}

	// Succeeded
	log.Infof("Job '%v' completed successfully", job.ID)
}
