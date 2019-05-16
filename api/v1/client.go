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

// Version is the version of ferryd
const Version = "0.5.0"

// Client is a client for the V1 API
type Client struct {
	client *http.Client
}

// NewClient will return a new ClientV1 for the local unix socket, suitable
// for communicating with the daemon.
func NewClient(address string) *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", address)
				},
				DisableKeepAlives:     false,
				IdleConnTimeout:       30 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 60 * time.Second,
		},
	}
}

// Close will kill any idle connections still in "keep-alive" and ensure we're
// not leaking file descriptors.
func (c *Client) Close() {
	c.client.Transport.(*http.Transport).CloseIdleConnections()
}

func formURI(part string) string {
	return fmt.Sprintf("http://localhost.localdomain:0/%s", part)
}

// GetStatus retrieves the status of the ferryd service
func (c *Client) GetStatus() (status StatusResponse, err error) {
	resp, err := c.client.Get(formURI("api/v1/status"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
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

// GetRepos will grab a list of repos from the daemon
func (c *Client) GetRepos() (list RepoList, err error) {
	resp, err := c.client.Get(formURI("api/v1/repos"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return
	}
	return
}

// DeleteRepo will attempt to remove a repository in the daemon
func (c *Client) DeleteRepo(id string) (gen GenericResponse, err error) {
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

// GetPoolItems will grab a list of pool items from the daemon
func (c *Client) GetPoolItems() (r PoolResponse, err error) {
	resp, err := c.client.Get(formURI("api/v1/pool"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}
	return
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

// DeltaRepo will generate missing metas in a given repo
func (c *Client) DeltaRepo(src string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewDeltaRepoJob(src))
}

// IndexRepo will attempt to index a repository in the daemon
func (c *Client) IndexRepo(src string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewIndexRepoJob(src))
}

// BulkAddPackages will ask ferryd to import the named packages with absolute paths
func (c *Client) BulkAddPackages(src string, paths []string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewBulkAddJob(src, paths))
}

// CloneRepo will ask the backend to clone an existing repository into a new repository
func (c *Client) CloneRepo(src, dest string, cloneAll bool) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewCloneRepoJob(src, dest, cloneAll))
}

// PullRepo will ask the backend to pull from target into repo
func (c *Client) PullRepo(src, dest string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewPullRepoJob(src, dest))
}

// RemoveSource will ask the backend to remove packages by source name
func (c *Client) RemoveSource(src, pkg string, release int) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewRemoveSourceJob(src, pkg, release))
}

// CopySource will ask the backend to copy packages by source name
func (c *Client) CopySource(src, dest, pkg string, release int) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewCopySourceJob(src, dest, pkg, release))
}

// TrimPackages will request that packages in the repo are trimmed to maxKeep
func (c *Client) TrimPackages(src string, maxKeep int) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewTrimPackagesJob(src, maxKeep))
}

// TrimObsolete will request that all packages marked obsolete are removed
func (c *Client) TrimObsolete(src string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewTrimObsoleteJob(src))
}

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
