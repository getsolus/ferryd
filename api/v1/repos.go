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
	"bytes"
	"encoding/json"
	"github.com/getsolus/ferryd/repo"
	"github.com/valyala/fasthttp"
	"net/http"
)

// Repos will grab a list of repos from the daemon
func (c *Client) Repos() (f repo.FullSummary, err error) {
	// Send the request
	resp, err := c.client.Get(formURI("api/v1/repos"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Decode the body as a full repo summary
	if err = json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return
	}
	return
}

// Repos will attempt to serialise our known repositories into a response
func (l *Listener) Repos(ctx *fasthttp.RequestCtx) {
	// Request a full repo summary
	f, err := l.manager.Repos()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
	}
	// Encode as JSON in the response
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&f); err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}
