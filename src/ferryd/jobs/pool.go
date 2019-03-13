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
	"runtime"
)

// A Pool is responsible for the main dispatch and bulking of jobs
// to ensure they're handled in the most optimal fashion.
type Pool struct {
	closed  bool
	njobs   int
	workers []*Worker
}

// NewPool will return a new Pool with the specified number
// of jobs. Note that "njobs" only refers to the number of *background jobs*,
// the majority of operations will run sequentially
func NewPool(store *JobStore, manager *core.Manager, njobs int) *Pool {
	// If we set to -1, we'll automatically set to half of the system core count
	// because we use xz -T 2 (so twice the number of threads ..)
	if njobs < 0 {
		njobs = runtime.NumCPU() / 2
	}

	log.Infof("Set runtime job limit: %d\n", njobs)

	ret := &Pool{
		closed: false,
		njobs:  njobs,
	}

	// Construct worker pool
	for i := 0; i < njobs; i++ {
		ret.workers = append(ret.workers, NewWorker(store, manager))
	}
	return ret
}

// Close an existing Pool, waiting for all jobs to complete
func (j *Pool) Close() {
	if j.closed {
		return
	}
	j.closed = true

	// Close all of our workers
	for _, j := range j.workers {
		j.Stop()
	}
}

// Begin will start the main job pool in parallel
func (j *Pool) Begin() {
	if j.closed {
		return
	}
	for _, j := range j.workers {
		go j.Start()
	}
}
