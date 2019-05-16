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
	"encoding/json"
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/ferryd/jobs"
	"github.com/valyala/fasthttp"
	"net/http"
	"runtime"
)

// getMethodOrigin helps us determine the caller so that we can print
// an appropriate method name into the log without tons of boilerplate
func getMethodCaller() string {
	n, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	if details := runtime.FuncForPC(n); details != nil {
		return details.Name()
	}
	return ""
}

// sendStockErrors is a utility to send a standard response to the ferry
// client that embeds the error message from ourside.
func (l *Listener) sendStockErrors(errs []error, ctx *fasthttp.RequestCtx) {
	errors := make([]string, len(errs))
	for i, err := range errs {
		errors[i] = err.Error()
		log.Errorf("Client communication error for method '%s', message: '%s'\n", getMethodCaller(), errors[i])
	}
	response := GenericResponse{
		Errors: errors,
	}
	buf := bytes.Buffer{}
	if e2 := json.NewEncoder(&buf).Encode(&response); e2 != nil {
		ctx.Error(e2.Error(), http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusBadRequest)
	ctx.SetBody(buf.Bytes())
}

// GetStatus will return the current status of the ferryd instance
func (l *Listener) GetStatus(ctx *fasthttp.RequestCtx) {
	ret := StatusResponse{
		TimeStarted: l.timeStarted,
		Version:     Version,
	}

	// Stuff the active jobs in
	jo, err := l.store.Active()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ret.Errors = append(ret.Errors, err.Error())
	}
	ret.CurrentJobs = jo

	fj, err := l.store.Failed()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ret.Errors = append(ret.Errors, err.Error())
	}
	ret.FailedJobs = fj

	cj, err := l.store.Completed()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ret.Errors = append(ret.Errors, err.Error())
	}
	ret.CompletedJobs = cj

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&ret); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}

// GetRepos will attempt to serialise our known repositories into a response
func (l *Listener) GetRepos(ctx *fasthttp.RequestCtx) {
	resp := RepoList{}
	//TODO: re-enable repos
	_, err := l.manager.GetRepos()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		resp.Errors = append(resp.Errors, err.Error())
	}
	//TODO: Uncomment this
	//resp.Repos = repos
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&resp); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}

// GetPoolItems will handle responding with the currently known pool items
func (l *Listener) GetPoolItems(ctx *fasthttp.RequestCtx) {
	resp := PoolResponse{}
	pools, err := l.manager.GetPoolItems()
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		resp.Errors = append(resp.Errors, err.Error())
		return
	}
	for _, pool := range pools {
		resp.Items = append(resp.Items, PoolItem{
			ID:       pool.Name,
			RefCount: int(pool.RefCount),
		})
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&resp); err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}

// CreateRepo will handle remote requests for repository creation
func (l *Listener) CreateRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	l.store.Push(jobs.NewCreateRepoJob(id))
}

// DeleteRepo will handle remote requests for repository deletion
func (l *Listener) DeleteRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	l.store.Push(jobs.NewDeleteRepoJob(id))
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
