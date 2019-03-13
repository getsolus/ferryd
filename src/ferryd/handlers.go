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

package main

import (
	"bytes"
	"encoding/json"
	"ferryd/jobs"
	log "github.com/DataDrake/waterlog"
	"github.com/julienschmidt/httprouter"
	"libferry"
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
func (api *APIListener) sendStockError(err error, w http.ResponseWriter, r *http.Request) {
	response := libferry.Response{
		Error:       true,
		ErrorString: err.Error(),
	}
	log.Errorf("Client communication error for method '%s', message: '%s'\n", getMethodCaller(), err.Error())
	buf := bytes.Buffer{}
	if e2 := json.NewEncoder(&buf).Encode(&response); e2 != nil {
		http.Error(w, e2.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(buf.Bytes())
}

// GetStatus will return the current status of the ferryd instance
func (api *APIListener) GetStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ret := libferry.StatusRequest{
		TimeStarted: s.timeStarted,
		Version:     libferry.Version,
	}

	// Stuff the active jobs in
	jo, err := s.store.ActiveJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret.CurrentJobs = jo

	fj, err := s.store.FailedJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret.FailedJobs = fj

	cj, err := s.store.CompletedJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ret.CompletedJobs = cj

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&ret); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

// GetRepos will attempt to serialise our known repositories into a response
func (api *APIListener) GetRepos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := libferry.RepoListingRequest{}
	repos, err := s.manager.GetRepos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, repo := range repos {
		req.Repository = append(req.Repository, repo.ID)
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

// GetPoolItems will handle responding with the currently known pool items
func (api *APIListener) GetPoolItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := libferry.PoolListingRequest{}
	pools, err := s.manager.GetPoolItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, pool := range pools {
		req.Item = append(req.Item, libferry.PoolItem{
			ID:       pool.Name,
			RefCount: int(pool.RefCount),
		})
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

// CreateRepo will handle remote requests for repository creation
func (api *APIListener) CreateRepo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	log.Infof("Creation of repo '%s' requested\n", id)
	api.store.PushJob(jobs.NewCreateRepoJob(id))
}

// DeleteRepo will handle remote requests for repository deletion
func (api *APIListener) DeleteRepo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	log.Infof("Deletion of repo '%s' requested\n", id)
	api.store.PushJob(jobs.NewDeleteRepoJob(id))
}

// DeltaRepo will handle remote requests for repository deltaing
func (api *APIListener) DeltaRepo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	log.Infof("Delta of repo '%s' requested\n", id)
	api.store.PushJob(jobs.NewDeltaRepoJob(id))
}

// IndexRepo will handle remote requests for repository indexing
func (api *APIListener) IndexRepo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	log.Infof("Index of repo '%s' requested\n", id)
	api.store.PushJob(jobs.NewIndexRepoJob(id))
}

// ImportPackages will bulk-import the packages in the request
func (api *APIListener) ImportPackages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	req := libferry.ImportRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Bulk import of '%d' packages for repo '%s' requested: '%v'\n", len(req.Path), id)

	api.store.PushJob(jobs.NewBulkAddJob(id, req.Path))
}

// CloneRepo will proxy a job to clone an existing repository
func (api *APIListener) CloneRepo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	req := libferry.CloneRepoRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Clone of repo '%s' into '%s' requested, full? '%b'\n", id, req.CloneName, req.CopyAll)

	api.store.PushJob(jobs.NewCloneRepoJob(id, req.CloneName, req.CopyAll))
}

// PullRepo will proxy a job to pull an existing repository
func (api *APIListener) PullRepo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	target := p.ByName("id")

	req := libferry.PullRepoRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Pulll of repo '%s' into '%s' requested\n", req.Source, target)

	api.store.PushJob(jobs.NewPullRepoJob(req.Source, target))
}

// RemoveSource will proxy a job to remove an existing set of packages by source name + relno
func (api *APIListener) RemoveSource(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	target := p.ByName("id")

	req := libferry.RemoveSourceRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Removal of release '%d' of source '%s' in repo '%s' requested", req.Release, req.Source, target)

	api.store.PushJob(jobs.NewRemoveSourceJob(target, req.Source, req.Release))
}

// CopySource will proxy a job to copy a package by source&relno into target
func (api *APIListener) CopySource(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sourceRepo := p.ByName("id")

	req := libferry.CopySourceRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Copy of release '%d' of source '%s' from repo '%s' to '%s' requested\n", req.Release, req.Source, sourceRepo, req.Target)

	api.store.PushJob(jobs.NewCopySourceJob(sourceRepo, req.Target, req.Source, req.Release))
}

// TrimPackages will proxy a job to remove excess fat from a repo
func (api *APIListener) TrimPackages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	target := p.ByName("id")

	req := libferry.TrimPackagesRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("Trim of packages with more than '%d' releases in repo '%s' requested\n", req.MaxKeep, target)

	api.store.PushJob(jobs.NewTrimPackagesJob(target, req.MaxKeep))
}

// TrimObsolete will proxy a job to remove obsolete packages from a repo
func (api *APIListener) TrimObsolete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	log.Infof("Trim of obsoletes in repo '%s' requested\n", id)
	api.store.PushJob(jobs.NewTrimObsoleteJob(id))
}

// ResetCompleted will ask the job store to remove completed jobs. This is blocking.
func (api *APIListener) ResetCompleted(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := s.store.ResetCompleted(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ResetFailed will ask the job store to remove failed jobs. This is blocking.
func (api *APIListener) ResetFailed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := s.store.ResetFailed(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
