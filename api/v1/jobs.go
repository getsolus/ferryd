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

func (c *Client) resetJobs(status string) (gen GenericResponse, err error) {
	req, err := http.NewRequest("DELETE", formURI("api/v1/jobs"), nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("status", status)
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

// ResetFailed asks the daemon to reset failed jobs
func (c *Client) ResetFailed() (gen GenericResponse, err error) {
	return c.resetJobs("failed")
}

// ResetCompleted asks the daemon to reset completed jobs
func (c *Client) ResetCompleted() (gen GenericResponse, err error) {
	return c.resetJobs("completed")
}

// ResetQueued asks the daemon to reset queued jobs
func (c *Client) ResetQueued() (gen GenericResponse, err error) {
	return c.resetJobs("queued")
}

func (c *Client) createJob(j *jobs.Job) (gen GenericResponse, err error) {
    raw, err := json.Marshal(j)
    if err != nil {
        return
    }
    req, err := http.NewRequest("POST", formURI("api/v1/jobs"), bytes.NewBuffer(raw))
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

// CreateJob will proxy a job to remove an existing set of packages by source name + relno
func (l *Listener) CreateJob(ctx *fasthttp.RequestCtx) {
    job := &jobs.Job{}
    if err := json.Unmarshal(ctx.Request.Body(), job); err != nil {
        l.sendStockErrors([]error{err}, ctx)
        return
    }
    l.store.Push(job)
}

// ResetJobs will ask the job store to remove jobs of a certain status
func (l *Listener) ResetJobs(ctx *fasthttp.RequestCtx) {
    status := string(ctx.QueryArgs().Peek("status"))
    if len(status) == 0 {
        err := "Job status required when resetting jobs"
        log.Errorln(err)
        ctx.Error(err, http.StatusBadRequest)
        return
    }
    switch status {
    case "completed":
        if err := l.store.ResetCompleted(); err != nil {
            log.Errorln(err)
            ctx.Error(err.Error(), http.StatusInternalServerError)
        }
    case "failed":
        if err := l.store.ResetFailed(); err != nil {
            log.Errorln(err)
            ctx.Error(err.Error(), http.StatusInternalServerError)
        }
    case "queued":
        if err := l.store.ResetQueued(); err != nil {
            log.Errorln(err)
            ctx.Error(err.Error(), http.StatusInternalServerError)
        }
    default:
        err := fmt.Sprintf("Invalid job status '%s' when resetting jobs", status)
        log.Errorln(err)
        ctx.Error(err, http.StatusBadRequest)
    }
}
