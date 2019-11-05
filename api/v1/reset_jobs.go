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
