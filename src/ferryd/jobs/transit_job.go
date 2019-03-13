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
	"os"
	"path/filepath"
)

// TransitJobHandler is responsible for accepting new upload payloads in the repository
type TransitJobHandler Job

// NewTransitJob will return a job suitable for adding to the job processor
func NewTransitJob(path string) *Job {
	return &Job{
		Type:    TransitProcess,
        Sources: []string{path},
	}
}

// NewTransitJobHandler will create a job handler for the input job and ensure it validates
func NewTransitJobHandler(j *Job) (handler *TransitJobHandler, err error) {
	if len(j.Sources) == 0 {
		err = fmt.Errorf("job is missing the path to the manifest")
        return
	}
	if len(j.Sources[0]) == 0 {
		err = fmt.Errorf("job is missing the path to the manifest")
        return
	}
	h := TransitJobHandler(*j)
    handler = &h
    return
}

// Execute will process incoming .tram files for potential repo inclusion
func (j *TransitJobHandler) Execute(jproc *Processor, manager *core.Manager) error {
	tram, err := core.NewTransitManifest(j.Sources[0])
	if err != nil {
		return err
	}

	if err = tram.ValidatePayload(); err != nil {
		return err
	}

	// Sanity.
	repo := tram.Manifest.Target
	if _, err := manager.GetRepo(repo); err != nil {
		return err
	}

	// Now try to merge into the repo
	pkgs := tram.GetPaths()
	if err = manager.AddPackages(repo, pkgs, true); err != nil {
		return err
	}

	log.Infof("Successfully processed manifest '%v' upload to repo '%s'\n", tram.ID(), repo)

	// Append the manifest path because now we'll want to delete these
	pkgs = append(pkgs, j.path)

	for _, p := range pkgs {
		if !core.PathExists(p) {
			continue
		}
		if err := os.Remove(p); err != nil {
			log.Errorf("Failed to remove manifest file '%s', reason: '%s'\n", p, err.Error())
		}
	}

	// At this point we should actually have valid pool entries so
	// we'll grab their names, and schedule that they be re-deltad.
	// It might be the case no delta is possible, but we'll let the
	// DeltaJobHandler decide on that.
	for _, pkg := range pkgs {
		pkgID := filepath.Base(pkg)
		p, ent := manager.GetPoolEntry(pkgID)
		if ent != nil {
			return err
		}
		jproc.PushJob(NewDeltaIndexJob(repo, p.Name))
	}

	return nil
}

// Describe returns a human readable description for this job
func (j *TransitJobHandler) Describe() string {
	return fmt.Sprintf("Process manifest '%s'", j.Sources[0])
}

// IsSerial returns true if a job should not be run alongside other jobs
func (J *TransitJobHandler) IsSerial() bool {
    return true
}
