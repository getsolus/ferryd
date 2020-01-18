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

package v1

import (
	"encoding/json"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"github.com/valyala/fasthttp"
	"net/http"
)

// Create will attempt to create a repository in the daemon
func (c *Client) Create(id string, instant bool) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("POST", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	// Set the query parameters
	i := "false"
	if instant {
		i = "true"
	}
	q := req.URL.Query()
	q.Add("instant", i)
	req.URL.RawQuery = q.Encode()
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Check for failure
	if resp.StatusCode != http.StatusOK {
		err = readError(resp.Body)
		return
	}
	// Decode the body as a Job
	err = json.NewDecoder(resp.Body).Decode(j)
	return
}

// CreateRepo will handle remote requests for repository creation
func (l *Listener) CreateRepo(ctx *fasthttp.RequestCtx) {
	// Get the query parameters
	id := ctx.UserValue("id").(string)
	imp := ctx.QueryArgs().GetBool("import")
	instant := ctx.QueryArgs().GetBool("instant")
	// Request the repo creation
	var jobID int
	var err error
	if imp {
		jobID, err = l.manager.Import(id, instant)
	} else {
		src := string(ctx.QueryArgs().Peek("clone"))
		if len(src) == 0 {
			jobID, err = l.manager.Create(id, instant)
		} else {
			jobID, err = l.manager.Clone(id, src)
		}
	}
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	// Write the ID for the client
	writeID(ctx, jobID)
}

// Import will ask ferryd to import a repository from disk
func (c *Client) Import(id string, instant bool) (r *repo.Summary, j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("POST", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	// Set the query parameters
	i := "false"
	if instant {
		i = "true"
	}
	q := req.URL.Query()
	q.Add("import", "true")
	q.Add("instant", i)
	req.URL.RawQuery = q.Encode()
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Check for failure
	if resp.StatusCode != http.StatusOK {
		err = readError(resp.Body)
		return
	}
	// Decode the body as a Job
	err = json.NewDecoder(resp.Body).Decode(j)
	return
}

// Remove will attempt to remove a repository in the daemon
func (c *Client) Remove(id string) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("DELETE", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Check for failure
	if resp.StatusCode != http.StatusOK {
		err = readError(resp.Body)
		return
	}
	// Decode the result as a Job
	err = json.NewDecoder(resp.Body).Decode(j)
	return
}

// RemoveRepo will handle remote requests for repository deletion
func (l *Listener) RemoveRepo(ctx *fasthttp.RequestCtx) {
	// Get the query parameters
	id := ctx.UserValue("id").(string)
	// Request the repo creation
	jobID, err := l.manager.Remove(id)
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	// Write the ID for the client
	writeID(ctx, jobID)
}
