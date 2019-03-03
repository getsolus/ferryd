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
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

// JobFetcher will be provided by either the Async or Sequential claim functions
type JobFetcher func() (*JobEntry, error)

// JobReaper will be provided by either the Async or Sequential retire functions
type JobReaper func(j *JobEntry) error

// MinWait is the minimum amount of time between retries for a worker
const MinWait = time.Second * 2

// MaxJitter sets the upper limit on the random jitter used for retry times
const MaxJitter int64 = 512

// A Worker is used to execute some portion of the incoming workload, and will
// keep polling for the correct job type to process
type Worker struct {
	sequential bool
	exit       chan int
	timer      *time.Timer
	wg         *sync.WaitGroup
	manager    *core.Manager
	store      *JobStore
	processor  *Processor

	fetcher JobFetcher // Fetch a new job
	reaper  JobReaper  // Purge an old job
}

// newWorker is an internal method to initialise a worker for usage
func newWorker(processor *Processor, sequential bool) *Worker {
	if processor.store == nil {
		panic("Constructed a Worker without a valid JobStore!")
	}
	if processor.wg == nil {
		panic("Constructed a Worker without a valid WaitGroup!")
	}

	w := &Worker{
		sequential: sequential,
		wg:         processor.wg,
		exit:       make(chan int, 1),
		timer:      nil, // Init this when we start up
		manager:    processor.manager,
		store:      processor.store,
		processor:  processor,
	}

	// Set up appropriate functions for dealing with jobs
	if sequential {
		w.fetcher = w.store.ClaimSequentialJob
		w.reaper = w.store.RetireSequentialJob
	} else {
		w.fetcher = w.store.ClaimAsyncJob
		w.reaper = w.store.RetireAsyncJob
	}

	return w
}

// NewWorkerAsync will return an asynchronous processing worker which will only
// pull from the store's async job queue
func NewWorkerAsync(processor *Processor) *Worker {
	return newWorker(processor, false)
}

// NewWorkerSequential will return a sequential worker operating on the main
// sequential job loop
func NewWorkerSequential(processor *Processor) *Worker {
	return newWorker(processor, true)
}

// Stop will demand that all new requests are no longer processed
func (w *Worker) Stop() {
	w.exit <- 1
	if w.timer != nil {
		w.timer.Stop()
	}
}

// Start will begin the main execution of this worker, and will continuously
// poll for new jobs with an increasing increment (with a ceiling limit)
func (w *Worker) Start() {
	defer w.wg.Done()

	// Let's get our timer initialised
	w.setTime()

	for {
		select {
		case <-w.exit:
			// Bail now, we've been told to go home
			return

		case <-w.timer.C:
			// Try to grab a job
			job, err := w.fetcher()

			// Report the error
			if err != nil {
				if err != ErrEmptyQueue {
					log.WithFields(log.Fields{
						"error": err,
						"async": !w.sequential,
					}).Error("Failed to grab a work queue item")
				}
				w.setTime()
				continue
			}

			// Got a job, now process it
			w.processJob(job)

			// Now we mark end time so we can calculate how long it took
			job.Timing.End = time.Now().UTC()

			// Mark the job as dealt with
			err = w.reaper(job)

			// Report failure in retiring the job
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"id":    job.GetID(),
					"type":  job.Type,
					"async": !w.sequential,
				}).Error("Error in retiring job")
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
	}
	w.timer.Reset(delay)
}

// processJob will actually examine the given job and figure out how
// to execute it. Each Worker can only execute a single job at a time
func (w *Worker) processJob(job *JobEntry) {
	handler, err := NewJobHandler(job)

	fields := log.Fields{
		"id":    job.GetID(),
		"type":  job.Type,
		"async": !w.sequential,
	}

	if err != nil {
		fields["error"] = err
		job.failure = err
		log.WithFields(fields).Error("No known job handler, cannot continue with job")
		return
	}

	// Safely have a handler now
	job.description = handler.Describe()
	fields["description"] = job.description

	// Try to execute it, report the error
	if err := handler.Execute(w.processor, w.manager); err != nil {
		fields["error"] = err
		job.failure = err
		log.WithFields(fields).Error("Job failed with error")
		return
	}

	// Succeeded
	log.WithFields(fields).Info("Job completed successfully")
}
