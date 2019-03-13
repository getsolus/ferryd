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
	"errors"
	"ferryd/core"
	"ferryd/jobs"
	log "github.com/DataDrake/waterlog"
	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"github.com/julienschmidt/httprouter"
	"github.com/radu-munteanu/fsnotify"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// Server sits on a unix socket accepting connections from authenticated
// client, i.e. root or those in the "ferry" group
type Server struct {
	running bool
	api     *ApiListener     // the HTTP socket handler
	manager *core.Manager    // heart of the story
	store   *jobs.JobStore   // Storage for jobs processor
	jproc   *jobs.Processor  // Allow scheduling jobs
	tl      *TransitListener //Listener for TRAM files

	// We store a global lock file ..
	lockFile *LockFile
	lockPath string
}

// NewServer will return a newly initialised Server which is currently unbound
func NewServer() (*Server, error) {
	// Before we can actually bind the socket, we must lock the file
	api.lockPath = filepath.Join(baseDir, LockFilePath)
	lfile, err := NewLockFile(api.lockPath)
	api.lockFile = lfile

	if err != nil {
		return nil, err
	}

	// Try to lock our lockfile now
	if err := api.lockFile.Lock(); err != nil {
		return nil, err
	}

	return &Server{
		running: false,
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
	s.jproc = jobs.NewProcessor(s.manager, s.store, backgroundJobCount)

	// Set up watching the manager's incoming directory
	if tl, err := NewTransitListener(s.manager.IncomingPath, s.store); err != nil {
		return err
	}

	// api
	api, err = NewAPIListener(s.store)
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
	s.jproc.Begin()
	s.tl.Start()
	err = s.api.Start()
	if err != nil {
		return err
	}

	if systemdEnabled {
		daemon.SdNotify(false, "READY=1")
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
	s.api.Stop()
	s.tl.Stop()
	s.jproc.Close()
	s.store.Close()
	s.manager.Close()
	s.running = false

	// We don't technically fully own it if systemd created it
	if !systemdEnabled {
		os.Remove(s.socketPath)
	}
}
