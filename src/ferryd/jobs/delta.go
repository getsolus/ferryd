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
	"fmt"
	log "github.com/DataDrake/waterlog"
	"libeopkg"
	"os"
	"sort"
)

// DeltaJobHandler is responsible for indexing repositories and should only
// ever be used in async queues. Deltas may take some time to produce and
// shouldn't be allowed to block the sequential processing queue.
type DeltaJobHandler Job

// NewDeltaJob will return a job suitable for adding to the job processor
func NewDeltaJob(repoID, packageID string) *Job {
	return &Job{
		Type:    Delta,
		SrcRepo: repoID,
		Sources: []string{packageID},
	}
}

// NewDeltaIndexJob will return a new job for creating delta packages as well
// as scheduling an index operation when complete.
func NewDeltaIndexJob(repoID, packageID string) *Job {
	return &Job{
		Type:    DeltaIndex,
		SrcRepo: repoID,
		Sources: []string{packageID},
	}
}

// NewDeltaJobHandler will create a job handler for the input job and ensure it validates
func NewDeltaJobHandler(j *Job) (handler *DeltaJobHandler, err error) {
	if len(j.SrcRepo) == 0 {
		err = fmt.Errorf("job is missing source repo")
		return
	}
	if len(j.Sources) == 0 {
		err = fmt.Errorf("job is missing a source package")
		return
	}
	h := DeltaJobHandler(*j)
	handler = &h
	return
}

// executeInternal is the common code shared in the delta jobs, and is
// split out to save duplication.
func (j *DeltaJobHandler) executeInternal(manager *core.Manager) (nDeltas int, err error) {
	pkgs, err := manager.GetPackages(j.SrcRepo, j.Sources[0])
	if err != nil {
		return
	}

	// Need at least 2 packages for a delta op.
	if len(pkgs) < 2 {
		log.Debugf("No delta is possible for package '%s' in repo '%s'\n", j.Sources[0], j.SrcRepo)
		return
	}

	sort.Sort(libeopkg.PackageSet(pkgs))
	tip := pkgs[len(pkgs)-1]

	// Process all potential deltas
	for i := 0; i < len(pkgs)-1; i++ {
		old := pkgs[i]

		deltaID := libeopkg.ComputeDeltaName(old, tip)

		// Don't need to report that it failed, we know this from history
		if manager.GetDeltaFailed(deltaID) {
			continue
		}

		hasDelta, e := manager.HasDelta(j.SrcRepo, j.Sources[0], deltaID)
		if e != nil {
			err = e
			return
		}

		// Package has this delta already? Continue.
		if hasDelta {
			continue
		}

		mapping := &core.DeltaInformation{
			FromID:      old.GetID(),
			ToID:        tip.GetID(),
			FromRelease: old.GetRelease(),
			ToRelease:   tip.GetRelease(),
		}

		// Before we go off creating it - does the delta package exist already?
		// If so, just re-ref it for usage within the new repo
		entry, e := manager.GetPoolEntry(deltaID)
		if entry != nil && e == nil {
			if e := manager.RefDelta(j.SrcRepo, deltaID); e != nil {
				log.Errorf("Failed to ref existing delta, id: %v\n", deltaID)
				err = e
				return
			}
			log.Debugf("Reused existing delta, id: %v\n", deltaID)
			continue
		}

		deltaPath, e := manager.CreateDelta(j.SrcRepo, old, tip)
		if e != nil {
			if err == libeopkg.ErrDeltaPointless {
				// Non-fatal, ask the manager to record this delta as a no-go
				log.Infof("Delta not possible, marked permanently: %v\n", deltaID)
				if e = manager.MarkDeltaFailed(deltaID, mapping); e != nil {
					log.Errorf("Failed to mark delta failure, id: %v\n", deltaID)
					err = e
					return
				}
				continue
			} else if err == libeopkg.ErrMismatchedDelta {
				log.Errorln("Package delta candidates do not match")
				continue
			} else {
				// Genuinely an issue now
				log.Errorf("Error in delta production, message: '%s'\n", e.Error())
				return
			}
		}

		nDeltas++

		// Produced a delta!
		log.Infof("Successfully producing delta package, path: '%s'\n", deltaPath)

		// Let's get it included now.
		if err = j.includeDelta(manager, mapping, deltaPath); err != nil {
			log.Error("Failed to include delta package, reason: '%s'\n", err.Error())
			return
		}
	}
	return
}

// includeDelta will wrap up the basic functionality to get a delta package
// imported into a target repository.
func (j *DeltaJobHandler) includeDelta(manager *core.Manager, mapping *core.DeltaInformation, deltaPath string) error {
	// Try to insert the delta
	if err := manager.AddDelta(j.SrcRepo, deltaPath, mapping); err != nil {
		return err
	}

	// Delete the deltaPath if the add is successful
	return os.Remove(deltaPath)
}

// Execute will delta the target package within the target repository.
func (j *DeltaJobHandler) Execute(_ *JobStore, manager *core.Manager) error {
	nDeltas, err := j.executeInternal(manager)
	if err != nil {
		return err
	}
	if j.Type != DeltaIndex {
		return nil
	}
	// Ask that our repository now be reindexed because we've added deltas but
	// only if we've successfully produced some delta packages
	if nDeltas < 1 {
		return nil
	}

	if err := manager.Index(j.SrcRepo); err != nil {
		log.Errorf("Failed to index repository '%s', reason: '%s'\n", j.SrcRepo, err.Error())
		return err
	}

	return nil
}

// Describe returns a human readable description for this job
func (j *DeltaJobHandler) Describe() string {
	if j.Type == DeltaIndex {
		return fmt.Sprintf("Delta package '%s' on '%s', then re-index", j.Sources[0], j.SrcRepo)
	}
	return fmt.Sprintf("Delta package '%s' on '%s'", j.Sources[0], j.SrcRepo)
}

// IsSerial returns true if a job should not be run alongside other jobs
func (j *DeltaJobHandler) IsSerial() bool {
	return false
}
