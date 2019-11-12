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
	"encoding/json"
	"fmt"
	"github.com/getsolus/ferryd/jobs"
	"github.com/olekukonko/tablewriter"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"sort"
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
	Current jobs.List
	// FailedJobs is a list of failed jobs
	Failed jobs.List
	// CompletedJobs is a list of completed jobs
	Completed jobs.List
}

// Uptime will determine the uptime of the daemon
func (s StatusResponse) Uptime() time.Duration {
	return time.Now().UTC().Sub(s.TimeStarted)
}

// Print out all the failed jobs
func (s StatusResponse) printFailed(out io.Writer) {
	fmt.Fprintf(out, "Failed jobs: (%d tracked)\n\n", len(s.Failed))
	if len(s.Failed) == 0 {
		return
	}
	sort.Sort(sort.Reverse(s.Failed))
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{
		"Status",
		"Completed",
		"Duration",
		"Description",
		"Error",
	})
	table.SetBorder(false)

	i := 0

	for _, j := range s.Failed {
		if i >= 10 {
			break
		}
		i++
		table.Append([]string{
			"failed",
			j.Finished.Time.Format(time.RFC3339),
			j.ExecutionTime().String(),
			j.Describe(),
			j.Message.String,
		})
	}
	table.Render()
}

// Print out all the completed jobs
func (s StatusResponse) printCompleted(out io.Writer) {
	fmt.Fprintf(out, "Completed jobs: (%d tracked)\n\n", len(s.Completed))
	if len(s.Completed) == 0 {
		return
	}
	sort.Sort(sort.Reverse(s.Completed))
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{
		"Status",
		"Completed",
		"Run Time",
		"Duration",
		"Description",
		"Message",
	})
	table.SetBorder(false)

	i := 0

	for _, j := range s.Completed {
		if i >= 10 {
			break
		}
		i++
		table.Append([]string{
			"completed",
			j.Finished.Time.Format(time.RFC3339),
			j.ExecutionTime().String(),
			j.TotalTime().String(),
			j.Describe(),
			j.Message.String,
		})
	}
	table.Render()
}

// Print out all the queued jobs
func (s StatusResponse) printCurrent(out io.Writer) {
	fmt.Fprintf(out, "Quened jobs: (%d tracked)\n\n", len(s.Current))
	if len(s.Completed) == 0 {
		return
	}
	sort.Sort(sort.Reverse(s.Completed))
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{
		"Status",
		"Queued",
		"Waiting",
		"Description",
	})
	table.SetBorder(false)

	i := 0

	for _, j := range s.Failed {
		if i >= 10 {
			break
		}
		i++
		status := "queued"
		if j.Status == jobs.Running {
			status = "running"
		}
		table.Append([]string{
			status,
			j.Finished.Time.Format(time.RFC3339),
			j.QueuedTime().String(),
			j.Describe(),
		})
	}
	table.Render()
}

// Print writes out a StatusResponse
func (s StatusResponse) Print(out io.Writer) {
	fmt.Fprintf(out, " - Daemon uptime: %v\n", s.Uptime())
	fmt.Fprintf(out, " - Daemon version: %v\n", s.Version)

	s.printFailed(out)
	s.printCurrent(out)
	s.printCompleted(out)
}

// Status retrieves the status of the ferryd service
func (c *Client) Status() (status StatusResponse, err error) {
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

// Status will return the current status of the ferryd instance
func (l *Listener) Status(ctx *fasthttp.RequestCtx) {
	ret := StatusResponse{
		TimeStarted: l.timeStarted,
		Version:     Version,
	}

	// Stuff the active jobs in
	jo, err := l.store.Active()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ret.Current = jo

	fj, err := l.store.Failed()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ret.Failed = fj

	cj, err := l.store.Completed()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ret.Completed = cj

	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&ret); err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}
