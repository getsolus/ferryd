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
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

// Restart terminates the running ferryd service and starts a new one
func (c *Client) Restart() error {
	// Create the request
	req, err := http.NewRequest("PATCH", formURI("api/v1/daemon"), nil)
	if err != nil {
		return err
	}
	// Set the query parameters
	q := req.URL.Query()
	q.Add("action", "restart")
	req.URL.RawQuery = q.Encode()
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Read the response
	return readError(resp.Body)
}

// ModifyDaemon makes a requested change to the daemon
func (l *Listener) ModifyDaemon(ctx *fasthttp.RequestCtx) {
	// Get the "action" query argument
	action := string(ctx.QueryArgs().Peek("action"))
	if len(action) == 0 {
		writeErrorString(ctx, "Action not specified when modifying daemon", http.StatusBadRequest)
		return
	}
	// Pivot by action
	switch action {
	case "restart":
		writeErrorString(ctx, "Restart of Daemon not yet implemented", http.StatusBadRequest)
	default:
		writeErrorString(ctx, fmt.Sprintf("Action '%s' not implemented for the daemon", action), http.StatusBadRequest)
	}
	return
}

// Stop terminates the running ferryd service
func (c *Client) Stop() error {
	// Create the request
	req, err := http.NewRequest("DELETE", formURI("api/v1/daemon"), nil)
	if err != nil {
		return err
	}
	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Read the response
	return readError(resp.Body)
}

// StopDaemon handles the stopping of the daemon
func (l *Listener) StopDaemon(ctx *fasthttp.RequestCtx) {
	writeErrorString(ctx, "Stopping of Daemon not yet implemented", http.StatusBadRequest)
	return
}
