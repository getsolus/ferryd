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
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
	"time"
)
// StatusResponse is a response from the 'status' endpoint
type StatusResponse struct {
	GenericResponse
	// TimeStarted is the time when the daemon was last started
	TimeStarted time.Time
	// Version is the version of the daemon
	Version string
	// CurrentJobs is a list of running and queued jobs
	CurrentJobs jobs.List
	// FailedJobs is a list of failed jobs
	FailedJobs jobs.List
	// CompletedJobs is a list of completed jobs
	CompletedJobs jobs.List
}

// Uptime will determine the uptime of the daemon
func (s StatusResponse) Uptime() time.Duration {
	return time.Now().UTC().Sub(s.TimeStarted)
}

// Print writes out a StatusResponse
func (resp StatusResponse) Print(out io.Writer) {

}

// Status will return the current status of the ferryd instance
func (l *Listener) Status(ctx *fasthttp.RequestCtx) {
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
