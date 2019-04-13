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
	"ferryd/api"
	"ferryd/core"
	"ferryd/jobs"
	log "github.com/DataDrake/waterlog"
	"github.com/coreos/go-systemd/daemon"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Server sits on a unix socket accepting connections from authenticated
// client, i.e. root or those in the "ferry" group
type Server struct {
	running bool
	api     *api.Listener    // the HTTP socket handler
	manager *core.Manager    // heart of the story
	store   *jobs.JobStore   // Storage for jobs processor
	pool    *jobs.Pool       // Allow scheduling jobs
	tl      *TransitListener //Listener for TRAM files

	// We store a global lock file ..
	lockFile *LockFile
	lockPath string
}

// NewServer will return a newly initialised Server which is currently unbound
func NewServer() (*Server, error) {
	// Before we can actually bind the socket, we must lock the file
	lockPath := filepath.Join(baseDir, LockFilePath)
	lfile, err := NewLockFile(lockPath)

	if err != nil {
		return nil, err
	}

	// Try to lock our lockfile now
	if err := lfile.Lock(); err != nil {
		return nil, err
	}

	return &Server{
		running:  false,
		lockPath: lockPath,
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
	// manager
	m, e := core.NewManager(baseDir)
	if e != nil {
		return e
	}
	s.manager = m

	// jobstore
	st, e := jobs.NewStore(baseDir)
	if e != nil {
		return e
	}
	s.store = st

	// processor
	s.pool = jobs.NewPool(s.store, s.manager, backgroundJobCount)

	// Set up watching the manager's incoming directory
	tl, err := NewTransitListener(s.manager.IncomingPath, s.store)
	if err != nil {
		return err
	}
	s.tl = tl

	// api
	api, err := api.NewListener(s.store, s.manager)
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
	s.pool.Begin()
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
	s.pool.Close()
	s.store.Close()
	s.manager.Close()
	s.running = false
}
