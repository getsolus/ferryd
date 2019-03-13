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

// APIListener sits on a unix socket accepting connections from authenticated
// client, i.e. root or those in the "ferry" group
type APIListener struct {
	srv        *http.Server
	router     *httprouter.Router
	socket     net.Listener
	socketPath string

	// When we first started up.
	timeStarted time.Time

	store *jobs.JobStore // Storage for jobs processor
}

// NewAPIListener will return a newly initialised Server which is currently unbound
func NewAPIListener(store *JobStore) (api *APIListener, err error) {
	router := httprouter.New()
	api = &APIListener{
		srv: &http.Server{
			Handler: router,
		},
		router:      router,
		timeStarted: time.Now().UTC(),
		store:       store,
	}

	// Set up the API bits
	router.GET("/api/v1/status", api.GetStatus)

	// Repo management
	router.GET("/api/v1/create/repo/:id", api.CreateRepo)
	router.GET("/api/v1/remove/repo/:id", api.DeleteRepo)
	router.GET("/api/v1/delta/repo/:id", api.DeltaRepo)
	router.GET("/api/v1/index/repo/:id", api.IndexRepo)

	// Client sends us data
	router.POST("/api/v1/import/:id", api.ImportPackages)
	router.POST("/api/v1/clone/:id", api.CloneRepo)
	router.POST("/api/v1/copy/source/:id", api.CopySource)
	router.POST("/api/v1/pull/:id", api.PullRepo)

	// Removal
	router.POST("/api/v1/remove/source/:id", api.RemoveSource)
	router.POST("/api/v1/trim/packages/:id", api.TrimPackages)
	router.GET("/api/v1/trim/obsoletes/:id", api.TrimObsolete)

	// Reset jobs are special and go straight to the store
	// We can't queue them as a job because we'd be in catch 22..
	router.GET("/api/v1/reset/completed", api.ResetCompleted)
	router.GET("/api/v1/reset/failed", api.ResetFailed)

	// List commands
	router.GET("/api/v1/list/repos", api.GetRepos)
	router.GET("/api/v1/list/pool", api.GetPoolItems)
	return s, nil
}

// Bind will attempt to set up the listener on the unix socket
// prior to serving.
func (api *APIListener) Bind() error {
	var listener net.Listener

	// Set from global CLI flag
	api.socketPath = socketPath

	// Check if we're systemd activated.
	if _, b := os.LookupEnv("LISTEN_FDS"); b {
		listeners, err := activation.Listeners(true)
		if err != nil {
			return err
		}
		if len(listeners) != 1 {
			return errors.New("expected a single unix socket")
		}
		// listener will be sockets[0], now we'll need to follow systemd activation path
		listener = listeners[0]
		// Mustn't delete!
		if unix, ok := listener.(*net.UnixListener); ok {
			unix.SetUnlinkOnClose(false)
		} else {
			return errors.New("expected unix socket")
		}
		systemdEnabled = true
	} else {
		l, e := net.Listen("unix", api.socketPath)
		if e != nil {
			return e
		}
		listener = l
	}

	uid := os.Getuid()
	gid := os.Getgid()
	if !systemdEnabled {
		// Avoid umask issues
		if e = os.Chown(s.socketPath, uid, gid); e != nil {
			return e
		}
		// Fatal if we cannot chmod the socket to be ours only
		if e = os.Chmod(api.socketPath, 0660); e != nil {
			return e
		}
	}
	api.socket = listener
	return nil
}

// Serve will continuously serve on the unix socket until dead
func (api *APIListener) Start() error {
	if api.socket == nil {
		return errors.New("Cannot serve without a bound server socket")
	}

	// Don't treat Shutdown/Close as an error, it's intended by us.
	if e := api.srv.Serve(api.socket); e != http.ErrServerClosed {
		return e
	}
	return nil
}

// Close will shut down and cleanup the socket
func (s *Server) Close() {
	api.srv.Shutdown(nil)

	// We don't technically fully own it if systemd created it
	if !systemdEnabled {
		os.Remove(api.socketPath)
	}
}
