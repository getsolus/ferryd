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

package daemon

import (
	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/coreos/go-systemd/daemon"
	"github.com/getsolus/ferryd/api/v1"
	"github.com/getsolus/ferryd/config"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/manager"
	"os"
	"os/signal"
	"syscall"
)

// Server sits on a unix socket accepting connections from authenticated
// client, i.e. root or those in the "ferry" group
type Server struct {
	api     *v1.Listener      // the HTTP socket handler
	manager *manager.Manager  // heart of the story
	tl      *manager.Listener // Listener for TRAM files

	// We store a global lock file ..
	lockFile *LockFile
	running  bool
}

// NewServer will return a newly initialised Server which is currently unbound
func NewServer() *Server {
	return &Server{
		lockFile: NewLockFile(config.Current.LockFile),
	}
}

// killHandler will ensure we cleanly tear down on a ctrl+c/sigint
func (s *Server) killHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Infoln("ferryd shutting down")
		s.Close()
		// Stop any mainLoop defers here
		os.Exit(1)
	}()
}

// Bind will attempt to set up the listener on the unix socket
// prior to serving.
func (s *Server) Bind() error {
	e := s.lockFile.Lock()
	if e != nil {
		return e
	}
	// Set up Job Store
	store, e := jobs.NewStore()
	if e != nil {
		return e
	}
	// Set up Repo Manager
	s.manager = manager.NewManager(store)
	// Set up Transit Listener
	tl, err := manager.NewListener(config.Current.TransitPath(), s.manager)
	if err != nil {
		return err
	}
	s.tl = tl
	// Set up the API
	api, err := v1.NewListener(store, s.manager)
	if err != nil {
		return err
	}
	s.api = api
	// Bind the API Server to its socket
	return s.api.Bind()
}

// Serve will continuously serve on the unix socket until dead
func (s *Server) Serve() error {
	// Set up waterlog
	log.SetOutput(os.Stderr)
	log.SetLevel(level.Debug)
	log.SetFormat(format.Un)
	s.running = true
	s.killHandler()
	// Serve the job queue
	s.tl.Start()
	if s.api.SystemdEnabled {
		ok, err := daemon.SdNotify(false, daemon.SdNotifyReady)
		if err != nil {
			log.Errorf("Failed to notify systemd, reason: '%s'\n", err.Error())
			return err
		}
		if !ok {
			log.Warnln("SdNotify failed due to missing environment variable")
		} else {
			log.Goodln("SdNotify successful")
		}
	}
	// Start the API server
	return s.api.Start()
}

// Close will shut down and cleanup the socket
func (s *Server) Close() {
	if !s.running {
		return
	}
	// Shutdown the services
	s.api.Close()
	s.tl.Stop()
	s.manager.Close()
	// Cleanup the lock file
	s.lockFile.Close()
	// Mark as stopped
	s.running = false
}
