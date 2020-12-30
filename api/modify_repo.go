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
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

func (c *Client) modifyRepo(id, action string) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("PATCH", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	// Set the query parameters
	q := req.URL.Query()
	q.Add("action", action)
	req.URL.RawQuery = q.Encode()
	// execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Read and decode the Job ID for the newly created job
	jobID, err := readID(resp)
	if err != nil {
		return
	}
	// wait for job to complete
	j, err = c.waitJob(jobID)
	return
}

// ModifyRepo performs a modification to an existing repo
func (l *Listener) ModifyRepo(ctx *fasthttp.RequestCtx) {
	// Get the "left" query parameter as the repo name
	id := ctx.UserValue("left").(string)
	if len(id) == 0 {
		writeErrorString(ctx, "ID required when modifying repo", http.StatusBadRequest)
		return
	}
	// Get the "action" query parameter
	action := string(ctx.QueryArgs().Peek("action"))
	if len(action) == 0 {
		writeErrorString(ctx, "Action required when modifying repo", http.StatusBadRequest)
		return
	}
	// Pivot by the requested action
	var err error
	var jobID int
	switch action {
	case "check":
		jobID, err = l.manager.Check(id)
	case "delta":
		jobID, err = l.manager.Delta(id)
	case "index":
		jobID, err = l.manager.Index(id)
	case "rescan":
		jobID, err = l.manager.Rescan(id)
	case "trim-packages":
		// Get the "max" query parameter
		max := string(ctx.QueryArgs().Peek("max"))
		if len(max) == 0 {
			writeErrorString(ctx, "Max required when trimming packages", http.StatusBadRequest)
			return
		}
		m, err := strconv.Atoi(max)
		if err != nil {
			writeErrorString(ctx, "Max must be an integer", http.StatusBadRequest)
			return
		}
		jobID, err = l.manager.TrimPackages(id, m)
	case "trim-obsoletes":
		jobID, err = l.manager.TrimObsoletes(id)
	default:
		writeErrorString(ctx, fmt.Sprintf("Invalid action '%s' when modifying repo", action), http.StatusBadRequest)
		return
	}
	// Check for any errors
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	// Send back the ID of the created job
	writeID(ctx, jobID)
}

// Check will compare a repo on disk with the DB
func (c *Client) Check(id string) (*jobs.Job, error) {
	return c.modifyRepo(id, "check")
}

// Delta will generate missing metas in a given repo
func (c *Client) Delta(id string) (j *jobs.Job, err error) {
	return c.modifyRepo(id, "delta")
}

// Index will attempt to index a repository in the daemon
func (c *Client) Index(id string) (j *jobs.Job, err error) {
	return c.modifyRepo(id, "index")
}

// Rescan will ask ferryd to re-import a repository from disk
func (c *Client) Rescan(id string) (j *jobs.Job, err error) {
	return c.modifyRepo(id, "rescan")
}

// TrimPackages will request that packages in the repo are trimmed to maxKeep
func (c *Client) TrimPackages(id string, maxKeep int) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("PATCH", formURI("api/v1/repos/"+id), nil)
	if err != nil {
		return
	}
	// Set the query parameters
	q := req.URL.Query()
	q.Add("action", "trim-packages")
	q.Add("max", strconv.Itoa(maxKeep))
	req.URL.RawQuery = q.Encode()
	// execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Read and decode the Job ID for the newly created job
	jobID, err := readID(resp)
	if err != nil {
		return
	}
	// wait for job to complete
	j, err = c.waitJob(jobID)
	return
}

// TrimObsoletes will request that all packages marked obsolete are removed
func (c *Client) TrimObsoletes(id string) (j *jobs.Job, err error) {
	return c.modifyRepo(id, "trim-obsoletes")
}
