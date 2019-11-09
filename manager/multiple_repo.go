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

package manager

import (
	"errors"
	"github.com/getsolus/ferryd/jobs"
	"github.com/getsolus/ferryd/repo"
)

/***************************/
/* MULTIPLE REPO FUNCTIONS */
/***************************/

// CherryPick syncs a single package from one repo to another
func (m *Manager) CherryPick(src, dest, pkg string) (int, error) {
	if len(src) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if len(src) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	if len(pkg) == 0 {
		return -1, errors.New("job is missing a package name")
	}
	j := &jobs.Job{
		Type: jobs.CherryPick,
		Src:  src,
		Dst:  dest,
		Pkg:  pkg,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// CherryPickExecute carries out a CherryPick Job
func (m *Manager) CherryPickExecute(j *jobs.Job) error {
	return errors.New("Function not inplemented")
}

// Clone creates a new repo as a copy of and existing repo
func (m *Manager) Clone(src, dst string) (int, error) {
	if len(src) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if len(dst) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Clone,
		Src:  src,
		Dst:  dst,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// CloneExecute carries out a clone job
func (m *Manager) CloneExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Compare reports on the differences between two repos
func (m *Manager) Compare(left, right string) (int, error) {
	if len(left) == 0 {
		return -1, errors.New("job is missing a left repo")
	}
	if len(right) == 0 {
		return -1, errors.New("job is missing a right repo")
	}
	j := &jobs.Job{
		Type: jobs.Compare,
		Src:  left,
		Dst:  right,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// CompareExecute carries out a Sync Job
func (m *Manager) CompareExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Sync compares two repos and makes changes so that "new" matches "old"
func (m *Manager) Sync(src, dst string) (int, error) {
	if len(src) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if len(dst) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Sync,
		Src:  src,
		Dst:  dst,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// SyncExecute carries out a Sync job
func (m *Manager) SyncExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Repos provides a summary of all available repos
func (m *Manager) Repos() (l repo.FullSummary, err error) {
	err = errors.New("Function not implemented")
	return
}
