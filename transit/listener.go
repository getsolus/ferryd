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

package transit

import (
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/ferryd/manager"
	"github.com/radu-munteanu/fsnotify"
	"os"
	"path/filepath"
	"strings"
)

// Listener is a process that creates jobs in response to new TRAM files
type Listener struct {
	base    string
	watcher *fsnotify.Watcher
	manager *manager.Manager
	stop    chan bool
	done    chan bool
}

// NewListener creates and sets up a new TransitListener
func NewListener(base []string, mgr *manager.Manager) (tl *Listener, err error) {
	// Create a new listener
	tl = &Listener{
		base:    filepath.Join(base...),
		manager: mgr,
		stop:    make(chan bool),
		done:    make(chan bool),
	}
	// Make sure the transit dir exists, creating it if missing
	if _, err = os.Stat(tl.base); os.IsNotExist(err) {
		if err = os.Mkdir(tl.base, 0755); err != nil {
			return
		}
	}
	// Create a watcher
	tl.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	// Monitor the incoming dir
	err = tl.watcher.Add(tl.base)
	return
}

// Start creates a gorouting than will wait for events on the incoming directory
// and process incoming .tram files
func (tl *Listener) Start() {
	go func() {
		for {
			select {
			case event := <-tl.watcher.Events:
				// Not interested in subdirs
				if filepath.Dir(event.Name) != tl.base {
					continue
				}
				// Filtering on Update events
				if event.Op&fsnotify.Update == fsnotify.Update {
					if strings.HasSuffix(event.Name, ManifestSuffix) {
						tl.processManifest(filepath.Base(event.Name))
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
func (tl *Listener) Stop() bool {
	tl.stop <- true
	return <-tl.done
}

// processManifest is invoked when a .tram file is closed in our incoming
// directory. We'll now push it for further processing
func (tl *Listener) processManifest(name string) {
	// Check that the path still exists
	fullpath := filepath.Join(tl.base, name)
	st, err := os.Stat(fullpath)
	if err != nil {
		return
	}
	// Make sure it is a regular file
	if !st.Mode().IsRegular() {
		return
	}
	// Transit the new packages
	log.Infof("Received transit manifest upload: '%s'\n", name)
	tl.manager.TransitPackage(fullpath)
}
