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

// Restart terminates the running ferryd service and starts a new one
func (c *Client) Restart() (status StatusResponse, err error) {
    req, err := http.NewRequest("PATCH", formURI("api/v1/daemon")
	if err != nil {
		return
	}
    q := req.URL.Query()
    q.Add("action", "restart")
    req.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return
	}
	return
}

// Status retrieves the status of the ferryd service
func (c *Client) Status() (status StatusResponse, err error) {
	resp, err := c.client.Get(formURI("api/v1/daemon/status"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return
	}
	return
}

// Stop terminates the running ferryd service
func (c *Client) Stop() (status StatusResponse, err error) {
    req, err := http.NewRequest("DELETE", formURI("api/v1/daemon")
	if err != nil {
		return
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return
	}
	return
}
