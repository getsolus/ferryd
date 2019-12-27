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
	"github.com/getsolus/ferryd/repo"
	"github.com/valyala/fasthttp"
	"net/http"
)

// CherryPick will ask the backend to sync a single package from one repo to another
func (c *Client) CherryPick(left, right, pkg string) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("PATCH", formURI("api/v1/repos/"+left+"/cherrypick/"+right), nil)
	if err != nil {
		return
	}
	// Set the query parameters
	q := req.URL.Query()
	q.Add("package", pkg)
	req.URL.RawQuery = q.Encode()
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Read back the Job ID
	jobID, err := readID(resp)
	if err != nil {
		return
	}
	// wait for job to complete
	j, err = c.waitJob(jobID)
	return
}

// CherryPickRepo will ask the backend to sync a single package from one repo to another
func (l *Listener) CherryPickRepo(ctx *fasthttp.RequestCtx) {
	// Get the repo IDs
	left := ctx.UserValue("left").(string)
	right := ctx.UserValue("right").(string)
	// Get the package name
	pkg := string(ctx.QueryArgs().Peek("package"))
	if len(pkg) == 0 {
		writeErrorString(ctx, "Package name required when cherry-picking", http.StatusBadRequest)
		return
	}
	// Request the cherry pick
	jobID, err := l.manager.CherryPick(left, right, pkg)
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	// write Job ID to the request
	writeID(ctx, jobID)
}

// Compare will ask the backend to compare one repo to another
func (c *Client) Compare(left, right string) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("GET", formURI("api/v1/repos/"+left+"/compare/"+right), nil)
	if err != nil {
		return
	}
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Read back the job ID
	jobID, err := readID(resp)
	if err != nil {
		return
	}
	// wait for job to complete
	j, err = c.waitJob(jobID)
	return
}

// CompareRepo will ask the backend to compare one repo to another
func (l *Listener) CompareRepo(ctx *fasthttp.RequestCtx) {
	// get the repo names
	left := ctx.UserValue("left").(string)
	right := ctx.UserValue("right").(string)
	// Request the comparison
	jobID, err := l.manager.Compare(left, right)
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	// write Job ID to the request
	writeID(ctx, jobID)
}

// Sync will ask the backend to sync one repo to another
func (c *Client) Sync(src, dst string) (j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("PATCH", formURI("api/v1/repos/"+src+"/sync/"+dst), nil)
	if err != nil {
		return
	}
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Read back the Job ID
	jobID, err := readID(resp)
	if err != nil {
		return
	}
	// wait for job to complete
	j, err = c.waitJob(jobID)
	return
}

// SyncRepo will ask the backend to sync one repo to another
func (l *Listener) SyncRepo(ctx *fasthttp.RequestCtx) {
	// Get the repo names
	left := ctx.UserValue("left").(string)
	right := ctx.UserValue("right").(string)
	// Request a Sync
	jobID, err := l.manager.Sync(left, right)
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	// write Job ID to the request
	writeID(ctx, jobID)
}

// Clone will ask the backend to clone an existing repository into a new repository
func (c *Client) Clone(src, dest string) (s *repo.Summary, j *jobs.Job, err error) {
	// Create a new request
	req, err := http.NewRequest("POST", formURI("api/v1/repos/"+src), nil)
	if err != nil {
		return
	}
	// Set the query parameters
	q := req.URL.Query()
	q.Add("clone", src)
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
