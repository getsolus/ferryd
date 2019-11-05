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

// CherryPick will ask the backend to sync a single package from one repo to another
func (c *Client) CherryPick(src, dest, pkg string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewCherryPickJob(src, dest, pkg))
}

// CloneRepo will ask the backend to clone an existing repository into a new repository
func (c *Client) CloneRepo(src, dest string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewCloneRepoJob(src, dest))
}

// SyncRepo will ask the backend to sync from one repo to another
func (c *Client) PullRepo(src, dest string) (gen GenericResponse, err error) {
	return c.createJob(jobs.NewSyncRepoJob(src, dest))
}
