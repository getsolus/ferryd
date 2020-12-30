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
	// Print header
	fmt.Fprintf(out, "Failed jobs: (%d tracked)\n\n", len(s.Failed))
	if len(s.Failed) == 0 {
		return
	}
	// Sort newest to oldest
	sort.Sort(sort.Reverse(s.Failed))
	// Setup for writing as a table
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{
		"Status",
		"Completed",
		"Duration",
		"Description",
		"Error",
	})
	table.SetBorder(false)
	// Print the 10 most recent failures
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
	// Print header
	fmt.Fprintf(out, "Completed jobs: (%d tracked)\n\n", len(s.Completed))
	if len(s.Completed) == 0 {
		return
	}
	// Sort from newest to oldest
	sort.Sort(sort.Reverse(s.Completed))
	// Setup for writing as a table
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
	// Print the 10 most recent completed jobs
	for i, j := range s.Completed {
		if i >= 9 {
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
	// Print the header
	fmt.Fprintf(out, "Queued jobs: (%d tracked)\n\n", len(s.Current))
	if len(s.Current) == 0 {
		return
	}
	// Sort from newest to oldest
	sort.Sort(sort.Reverse(s.Current))
	// Setup to print as a table
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{
		"Status",
		"Created",
		"Waiting For",
		"Description",
	})
	table.SetBorder(false)
	// Print the 10 most recently queued jobs
	i := 0
	for _, j := range s.Current {
		if i >= 10 {
			break
		}
		i++
		if j.Status == jobs.Running {
			table.Append([]string{
				"running",
				j.Created.Time.Format(time.RFC3339),
				j.QueuedSince().String(),
				j.Describe(),
			})
		} else {
			table.Append([]string{
				"queued",
				j.Created.Time.Format(time.RFC3339),
				j.QueuedSince().String(),
				j.Describe(),
			})
		}
	}
	table.Render()
}

// Print writes out a StatusResponse
func (s StatusResponse) Print(out io.Writer) {
	// Print daemon statistics
	fmt.Fprintf(out, " - Daemon uptime: %v\n", s.Uptime())
	fmt.Fprintf(out, " - Daemon version: %v\n\n", s.Version)
	// Print jobs
	s.printFailed(out)
	println()
	s.printCurrent(out)
	println()
	s.printCompleted(out)
}

// Status retrieves the status of the ferryd service
func (c *Client) Status() (status StatusResponse, err error) {
	// Send the request
	resp, err := c.client.Get(formURI("api/v1/status"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Decode the body as a StatusResponse
	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return
	}
	return
}

// Status will return the current status of the ferryd instance
func (l *Listener) Status(ctx *fasthttp.RequestCtx) {
	// Create the new response
	ret := StatusResponse{
		TimeStarted: l.timeStarted,
		Version:     Version,
	}
	// Add the active jobs
	jo, err := l.store.Active()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ret.Current = jo
	// Add the failed jobs
	fj, err := l.store.Failed()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ret.Failed = fj
	// Add the completed jobs
	cj, err := l.store.Completed()
	if err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ret.Completed = cj
	// Encode the StatusResponse as JSON in the body
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&ret); err != nil {
		writeError(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.SetBody(buf.Bytes())
}
