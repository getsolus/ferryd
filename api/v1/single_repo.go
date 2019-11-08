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
	"encoding/json"
	"github.com/getsolus/ferryd/jobs"
	"github.com/valyala/fasthttp"
	"net/http"
)

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
	//id := ctx.UserValue("id").(string)
	//jobID, err := l.manager.CreateRepo(id)
	writeErrorString(ctx, "Not yet implemented", http.StatusInternalServerError)
}

// ImportRepo will ask ferryd to import a repository from disk
func (c *Client) ImportRepo(id string) (j *jobs.Job, err error) {
	err = errors.New("Not yet implemented")
	return
}

// RemoveRepo will attempt to remove a repository in the daemon
func (c *Client) RemoveRepo(id string) (j *jobs.Job, err error) {
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
	// id := ctx.UserValue("id").(string)
	// l.store.Push(jobs.NewRemoveRepoJob(id))
	writeErrorString(ctx, "Not yet implemented", http.StatusInternalServerError)
}
