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
	"bytes"
	"encoding/json"
	log "github.com/DataDrake/waterlog"
	"github.com/getsolus/ferryd/client"
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

// sendStockError is a utility to send a standard response to the ferry
// client that embeds the error message from ourside.
func (l *Listener) sendStockError(err error, ctx *fasthttp.RequestCtx) {
	response := client.Response{
		Error:       true,
		ErrorString: err.Error(),
	}
	log.Errorf("Client communication error for method '%s', message: '%s'\n", getMethodCaller(), err.Error())
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
	ret := client.StatusRequest{
		TimeStarted: l.timeStarted,
		Version:     client.Version,
	}

	// Stuff the active jobs in
	jo, err := l.store.Active()
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ret.CurrentJobs = jo

	fj, err := l.store.Failed()
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ret.FailedJobs = fj

	cj, err := l.store.Completed()
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ret.CompletedJobs = cj

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&ret); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}

// GetRepos will attempt to serialise our known repositories into a response
func (l *Listener) GetRepos(ctx *fasthttp.RequestCtx) {
	req := client.RepoListingRequest{}
	repos, err := l.manager.GetRepos()
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	for _, repo := range repos {
		req.Repository = append(req.Repository, repo.ID)
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}

// GetPoolItems will handle responding with the currently known pool items
func (l *Listener) GetPoolItems(ctx *fasthttp.RequestCtx) {
	req := client.PoolListingRequest{}
	pools, err := l.manager.GetPoolItems()
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	for _, pool := range pools {
		req.Item = append(req.Item, client.PoolItem{
			ID:       pool.Name,
			RefCount: int(pool.RefCount),
		})
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}

// CreateRepo will handle remote requests for repository creation
func (l *Listener) CreateRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	log.Infof("Creation of repo '%s' requested\n", id)
	l.store.Push(jobs.NewCreateRepoJob(id))
}

// DeleteRepo will handle remote requests for repository deletion
func (l *Listener) DeleteRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	log.Infof("Deletion of repo '%s' requested\n", id)
	l.store.Push(jobs.NewDeleteRepoJob(id))
}

// DeltaRepo will handle remote requests for repository deltaing
func (l *Listener) DeltaRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	log.Infof("Delta of repo '%s' requested\n", id)
	l.store.Push(jobs.NewDeltaRepoJob(id))
}

// IndexRepo will handle remote requests for repository indexing
func (l *Listener) IndexRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	log.Infof("Index of repo '%s' requested\n", id)
	l.store.Push(jobs.NewIndexRepoJob(id))
}

// ImportPackages will bulk-import the packages in the request
func (l *Listener) ImportPackages(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)

	req := client.ImportRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Bulk import of '%d' packages for repo '%s' requested: '%v'\n", len(req.Path), id, req.Path)

	l.store.Push(jobs.NewBulkAddJob(id, req.Path))
}

// CloneRepo will proxy a job to clone an existing repository
func (l *Listener) CloneRepo(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)

	req := client.CloneRepoRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Clone of repo '%s' into '%s' requested, full? '%t'\n", id, req.CloneName, req.CopyAll)

	l.store.Push(jobs.NewCloneRepoJob(id, req.CloneName, req.CopyAll))
}

// PullRepo will proxy a job to pull an existing repository
func (l *Listener) PullRepo(ctx *fasthttp.RequestCtx) {
	target := ctx.UserValue("id").(string)

	req := client.PullRepoRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Pull of repo '%s' into '%s' requested\n", req.Source, target)

	l.store.Push(jobs.NewPullRepoJob(req.Source, target))
}

// RemoveSource will proxy a job to remove an existing set of packages by source name + relno
func (l *Listener) RemoveSource(ctx *fasthttp.RequestCtx) {
	target := ctx.UserValue("id").(string)

	req := client.RemoveSourceRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Removal of release '%d' of source '%s' in repo '%s' requested", req.Release, req.Source, target)

	l.store.Push(jobs.NewRemoveSourceJob(target, req.Source, req.Release))
}

// CopySource will proxy a job to copy a package by source&relno into target
func (l *Listener) CopySource(ctx *fasthttp.RequestCtx) {
	sourceRepo := ctx.UserValue("id").(string)

	req := client.CopySourceRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Copy of release '%d' of source '%s' from repo '%s' to '%s' requested\n", req.Release, req.Source, sourceRepo, req.Target)

	l.store.Push(jobs.NewCopySourceJob(sourceRepo, req.Target, req.Source, req.Release))
}

// TrimPackages will proxy a job to remove excess fat from a repo
func (l *Listener) TrimPackages(ctx *fasthttp.RequestCtx) {
	target := ctx.UserValue("id").(string)

	req := client.TrimPackagesRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Trim of packages with more than '%d' releases in repo '%s' requested\n", req.MaxKeep, target)

	l.store.Push(jobs.NewTrimPackagesJob(target, req.MaxKeep))
}

// TrimObsolete will proxy a job to remove obsolete packages from a repo
func (l *Listener) TrimObsolete(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("id").(string)
	log.Infof("Trim of obsoletes in repo '%s' requested\n", id)
	l.store.Push(jobs.NewTrimObsoleteJob(id))
}

// ResetCompleted will ask the job store to remove completed jobs. This is blocking.
func (l *Listener) ResetCompleted(ctx *fasthttp.RequestCtx) {
	if err := l.store.ResetCompleted(); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
	}
}

// ResetFailed will ask the job store to remove failed jobs. This is blocking.
func (l *Listener) ResetFailed(ctx *fasthttp.RequestCtx) {
	if err := l.store.ResetFailed(); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
	}
}
