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

package api

import (
	"errors"
	"ferryd/core"
	"ferryd/jobs"
	log "github.com/DataDrake/waterlog"
	"github.com/coreos/go-systemd/activation"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"os"
	"time"
)

// Default socket path we expect to use
const defaultSocketPath = "/run/ferryd.sock"

// SocketPath is the path to find the ferryd socket
var SocketPath string

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

	store   *jobs.JobStore // Storage for jobs processor
	manager *core.Manager  // manager of the repos
}

// NewListener will return a newly initialised Server which is currently unbound
func NewListener(store *jobs.JobStore, manager *core.Manager) (api *Listener, err error) {
	r := router.New()
	api = &Listener{
		srv: &fasthttp.Server{
			Handler: r.Handler,
		},
		router:         r,
		SystemdEnabled: false,
		timeStarted:    time.Now().UTC(),
		store:          store,
		manager:        manager,
	}

	// Set up the API bits
	r.GET("/api/v1/status", api.GetStatus)

	// Repo management
	r.GET("/api/v1/create/repo/:id", api.CreateRepo)
	r.GET("/api/v1/remove/repo/:id", api.DeleteRepo)
	r.GET("/api/v1/delta/repo/:id", api.DeltaRepo)
	r.GET("/api/v1/index/repo/:id", api.IndexRepo)

	// Client sends us data
	r.POST("/api/v1/import/:id", api.ImportPackages)
	r.POST("/api/v1/clone/:id", api.CloneRepo)
	r.POST("/api/v1/copy/source/:id", api.CopySource)
	r.POST("/api/v1/pull/:id", api.PullRepo)

	// Removal
	r.POST("/api/v1/remove/source/:id", api.RemoveSource)
	r.POST("/api/v1/trim/packages/:id", api.TrimPackages)
	r.GET("/api/v1/trim/obsoletes/:id", api.TrimObsolete)

	// Reset jobs are special and go straight to the store
	// We can't queue them as a job because we'd be in catch 22..
	r.GET("/api/v1/reset/completed", api.ResetCompleted)
	r.GET("/api/v1/reset/failed", api.ResetFailed)

	// List commands
	r.GET("/api/v1/list/repos", api.GetRepos)
	r.GET("/api/v1/list/pool", api.GetPoolItems)
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
