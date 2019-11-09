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

package daemon

import (
	log "github.com/DataDrake/waterlog"
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
	running bool
	api     *v1.Listener     // the HTTP socket handler
	manager *manager.Manager // heart of the story
	store   *jobs.Store      // Storage for jobs processor
	tl      *TransitListener //Listener for TRAM files

	// We store a global lock file ..
	lockFile *LockFile
	lockPath string
}

// NewServer will return a newly initialised Server which is currently unbound
func NewServer() (*Server, error) {
	// Before we can actually bind the socket, we must lock the file
	lfile, err := NewLockFile(config.Current.LockFile)

	if err != nil {
		return nil, err
	}

	// Try to lock our lockfile now
	if err := lfile.Lock(); err != nil {
		return nil, err
	}

	return &Server{
		running:  false,
		lockPath: config.Current.LockFile,
		lockFile: lfile,
	}, nil
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
	// Load config from file
	config.Load()

	// Store
	st, e := jobs.NewStore()
	if e != nil {
		return e
	}
	s.store = st

	// manager
	s.manager = manager.NewManager(st)

	// Set up watching the manager's incoming directory
	tl, err := NewTransitListener(config.Current.TransitPath(), s.manager)
	if err != nil {
		return err
	}
	s.tl = tl

	// api
	api, err := v1.NewListener(s.store, s.manager)
	if err != nil {
		return err
	}
	s.api = api
	return s.api.Bind()
}

// Serve will continuously serve on the unix socket until dead
func (s *Server) Serve() error {
	s.running = true
	s.killHandler()
	defer func() {
		s.running = false
	}()
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
	err := s.api.Start()
	if err != nil {
		return err
	}

	return nil
}

// Close will shut down and cleanup the socket
func (s *Server) Close() {
	if !s.running {
		return
	}
	if s.lockFile != nil {
		s.lockFile.Unlock()
		s.lockFile.Clean()
		s.lockFile = nil
	}
	s.api.Close()
	s.tl.Stop()
	s.store.Close()
	s.manager.Close()
	s.running = false
}
