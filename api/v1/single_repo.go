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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"net"
	"net/http"
	//"runtime"
	//"strings"
	"time"
)


// CheckRepo will compare a repo on disk with the DB
func (c *Client) CheckRepo(id string) (gen GenericResponse, err error) {
	req, err := http.NewRequest("POST", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
    q := req.URL.Query()
    q.Add("action", check)
    req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&gen); err != nil {
		return
	}
	return
}

// CreateRepo will attempt to create a repository in the daemon
func (c *Client) CreateRepo(id string) (gen GenericResponse, err error) {
	req, err := http.NewRequest("POST", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&gen); err != nil {
		return
	}
	return
}

// CreateRepo will handle remote requests for repository creation
func (l *Listener) CreateRepo(ctx *fasthttp.RequestCtx) {
    id := ctx.UserValue("id").(string)
    l.store.Push(jobs.NewCreateRepoJob(id))
}


// DeltaRepo will generate missing metas in a given repo
func (c *Client) DeltaRepo(id string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewDeltaRepoJob(id))
}

// Import will ask ferryd to import a repository from disk
func (c *Client) ImportRepo(id string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewImportRepoJob(id))
}

// IndexRepo will attempt to index a repository in the daemon
func (c *Client) IndexRepo(id string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewIndexRepoJob(id))
}

// RemoveRepo will attempt to remove a repository in the daemon
func (c *Client) RemoveRepo(id string) (gen GenericResponse, err error) {
	req, err := http.NewRequest("DELETE", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&gen); err != nil {
		return
	}
	return
}

// RemoveRepo will handle remote requests for repository deletion
func (l *Listener) RemoveRepo(ctx *fasthttp.RequestCtx) {
    id := ctx.UserValue("id").(string)
    l.store.Push(jobs.NewRemoveRepoJob(id))
}

// Rescan will ask ferryd to re-import a repository from disk
func (c *Client) RescanRepo(id string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewRescanRepoJob(id))
}

// TrimPackages will request that packages in the repo are trimmed to maxKeep
func (c *Client) TrimPackages(id string, maxKeep int) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewTrimPackagesJob(id, maxKeep))
}

// TrimObsolete will request that all packages marked obsolete are removed
func (c *Client) TrimObsolete(id string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewTrimObsoleteJob(id))
}
