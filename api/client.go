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
	"context"
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"net"
	"net/http"
	"time"
)

// Version is the version of ferryd
const Version = "1.0.0"

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

// waitJob retries periodically to read back a job.
func (c *Client) waitJob(id int) (j *jobs.Job, err error) {
	start := time.Now()
	for {
		// wait
		time.Sleep(time.Second)
		// request
		j, err = c.GetJob(id)
		if err != nil {
			return
		}
		// Stop if job is finished
		if j.Status > jobs.Running {
			return
		}
		fmt.Printf("Elapsed Time: %s\n", time.Now().Sub(start).String())
	}
}

// waitDiff retries periodically to read back a job with a diff as the result
func (c *Client) waitDiff(id int) (d *repo.Diff, j *jobs.Job, err error) {
	j, err = c.waitJob(id)
	if err != nil {
		return
	}
	if err = d.UnmarshalBinary(j.Results); err != nil {
		err = fmt.Errorf("error while decoding diff: %v", err)
	}
	return
}
