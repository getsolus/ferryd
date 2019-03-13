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

package main

import (
	"ferryd/core"
	"ferryd/jobs"
	log "github.com/DataDrake/waterlog"
	"github.com/radu-munteanu/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

// TransitListener is a process that creates jobs in response to new TRAM files
type TransitListener struct {
	base    string
	watcher *fsnotify.Watcher
	store   *JobStore
	stop    chan bool
	done    chan bool
}

// NewTransitListener creates and sets up a new TransitListener
func NewTransitListener(base string, store *JobStore) (tl *TransitListener, err error) {
	tl = &TtransitListener{
		base:  base,
		store: store,
		stop:  make(chan bool),
		done:  make(chan bool),
	}
	tl.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	// Monitor the incoming dir
	if err = tl.watcher.Add(tl.base); err != nil {
		return err
	}
	return
}

// Start creates a gorouting than will wait for events on the incoming directory
// and process incoming .tram files
func (tl *TransitListener) Start() {
	go func() {
		defer s.watchGroup.Done()
		for {
			select {
			case event := <-s.watcher.Events:
				// Not interested in subdirs
				if filepath.Dir(event.Name) != tl.base {
					continue
				}
				if event.Op&fsnotify.Update == fsnotify.Update {
					if strings.HasSuffix(event.Name, core.TransitManifestSuffix) {
						s.processTransitManifest(filepath.Base(event.Name))
					}
				}
			case <-tl.stop:
				tl.done <- true
				return
			}
		}
	}()
}

// Stop will force the fsnotify code to shut down
func (tl *TransitListener) Stop() bool {
	s.stop <- true
	return <-s.done
}

// processTransitManifest is invoked when a .tram file is closed in our incoming
// directory. We'll now push it for further processing
func (tl *TransitListener) processTransitManifest(name string) {
	fullpath := filepath.Join(tl.base, name)

	st, err := os.Stat(fullpath)
	if err != nil {
		return
	}

	if !st.Mode().IsRegular() {
		return
	}

	log.Infof("Received transit manifest upload: '%s'\n", name)
	s.store.PushJob(jobs.NewTransitJob(fullpath))
}
