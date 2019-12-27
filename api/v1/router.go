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

package v1

import (
	"errors"
	log "github.com/DataDrake/waterlog"
	"github.com/coreos/go-systemd/activation"
	"github.com/fasthttp/router"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/manager"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"os"
	"time"
)

// SocketPath is the path to find the ferryd socket
const SocketPath = "/run/ferryd.sock"

// Listener sits on a unix socket accepting connections from authenticated
// client, i.e. root or those in the "ferry" group
type Listener struct {
	srv    *fasthttp.Server
	router *router.Router
	socket net.Listener
	// If systemd is enabled, we'll talk to it.
	SystemdEnabled bool
	// When we first started up.
	timeStarted time.Time

	store   *jobs.Store      // Storage for jobs processor
	manager *manager.Manager // manager of the repos
}

// NewListener will return a newly initialised Server which is currently unbound
func NewListener(store *jobs.Store, mgr *manager.Manager) (api *Listener, err error) {
	r := router.New()
	api = &Listener{
		srv: &fasthttp.Server{
			Handler: r.Handler,
		},
		router:         r,
		SystemdEnabled: false,
		timeStarted:    time.Now().UTC(),
		store:          store,
		manager:        mgr,
	}

	// Set up the API bits
	// Daemon Management
	r.GET("/api/v1/status", api.Status)
	r.PATCH("/api/v1/daemon", api.ModifyDaemon) // restart only, for now
	r.DELETE("/api/v1/daemon", api.StopDaemon)

	// Repo management
	r.GET("/api/v1/repos", api.Repos)           // Summaries of all repos
	r.POST("/api/v1/repos/:id", api.CreateRepo) // Clone, Create, Import
	// r.GET("/api/v1/repos/:id", api.GetRepo) // Summary of repo
	r.PATCH("/api/v1/repos/:id", api.ModifyRepo) // ?action={check, delta, index, rescan, trim-packages, trim-obsoletes}
	r.DELETE("/api/v1/repos/:id", api.RemoveRepo)

	r.PATCH("/api/v1/repos/:left/cherrypick/:right", api.CherryPickRepo)
	r.GET("/api/v1/repos/:left/compare/:right", api.CompareRepo)
	r.PATCH("/api/v1/repos/:left/sync/:right", api.SyncRepo)

	// Job Management
	r.DELETE("/api/v1/jobs", api.ResetJobs) // ?status={completed,failed,queued}
	r.GET("/api/v1/jobs/:id", api.GetJob)
	//r.DELETE("/api/v1/jobs/:id", api.CancelJob)

	return api, nil
}

// Bind will attempt to set up the listener on the unix socket
// prior to serving.
func (api *Listener) Bind() error {
	var listener net.Listener

	// Check if we're systemd activated.
	if v, b := os.LookupEnv("LISTEN_FDS"); b {
		log.Debugf("LISTEN_FDS: %v\n", v)
		listeners, err := activation.Listeners()
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
		api.SystemdEnabled = true
	} else {
		l, e := net.Listen("unix", SocketPath)
		if e != nil {
			return e
		}
		listener = l
	}

	uid := os.Getuid()
	gid := os.Getgid()
	if !api.SystemdEnabled {
		// Avoid umask issues
		if e := os.Chown(SocketPath, uid, gid); e != nil {
			return e
		}
		// Fatal if we cannot chmod the socket to be ours only
		if e := os.Chmod(SocketPath, 0660); e != nil {
			return e
		}
	}
	api.socket = listener
	return nil
}

// Start will continuously serve on the unix socket until dead
func (api *Listener) Start() error {
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
func (api *Listener) Close() {
	api.srv.Shutdown()

	// We don't technically fully own it if systemd created it
	if !api.SystemdEnabled {
		os.Remove(SocketPath)
	}
}
