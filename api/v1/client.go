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
