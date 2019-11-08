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
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

// GetJob will request current information about a Job
func (c *Client) GetJob(id int) (j *jobs.Job, err error) {
	resp, err := c.client.Get(fmt.Sprintf("/api/v1/jobs/%d", id))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = readError(resp.Body)
		return
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(j)
	return
}

// GetJob handles request the current information about a Job
func (l *Listener) GetJob(ctx *fasthttp.RequestCtx) {
	idString := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idString)
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	var job *jobs.Job
	if job, err := l.store.GetJob(id); err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(ctx.Response.BodyWriter())
	err = enc.Encode(job)
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
	}
}

func (c *Client) resetJobs(status string) error {
	req, err := http.NewRequest("DELETE", formURI("api/v1/jobs"), nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("status", status)
	req.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return readError(resp.Body)
	}
	return nil
}

// ResetFailed asks the daemon to reset failed jobs
func (c *Client) ResetFailed() error {
	return c.resetJobs("failed")
}

// ResetCompleted asks the daemon to reset completed jobs
func (c *Client) ResetCompleted() error {
	return c.resetJobs("completed")
}

// ResetQueued asks the daemon to reset queued jobs
func (c *Client) ResetQueued() error {
	return c.resetJobs("queued")
}

// ResetJobs will ask the job store to remove jobs of a certain status
func (l *Listener) ResetJobs(ctx *fasthttp.RequestCtx) {
	status := string(ctx.QueryArgs().Peek("status"))
	if len(status) == 0 {
		writeErrorString(ctx, "Job status required when resetting jobs", http.StatusBadRequest)
		return
	}
	var err error
	switch status {
	case "completed":
		err = l.store.ResetCompleted()
	case "failed":
		err = l.store.ResetFailed()
	case "queued":
		err = l.store.ResetQueued()
	default:
		writeErrorString(ctx, fmt.Sprintf("Invalid job status '%s' when resetting jobs", status), http.StatusBadRequest)
		return
	}
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
	}
}
