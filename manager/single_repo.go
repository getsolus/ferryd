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
)

/*************************/
/* SINGLE REPO FUNCTIONS */
/*************************/

// Check compares an existing repo on Disk with its DB
func (m *Manager) Check(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.Check,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// CheckExecute carries out a Check job
func (m *Manager) CheckExecute(j *jobs.Job) error {
	if len(j.Src) == 0 {
		return errors.New("job is missing a source repo")
	}
	return errors.New("Function not implemented")
}

// Create sets up a new repo
func (m *Manager) Create(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Create,
		Dst:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// CreateExecute carries out a Create job
func (m *Manager) CreateExecute(j *jobs.Job) error {
	if len(j.Dst) == 0 {
		return errors.New("job is missing a destination repo")
	}
	return errors.New("Function not implemented")
}

// Delta generates missing package deltas for an entire repo
func (m *Manager) Delta(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Delta,
		Dst:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// DeltaExecute carries out a Delta job
func (m *Manager) DeltaExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// DeltaPackage generates missing package deltas for a single package
func (m *Manager) DeltaPackage(dst, pkg string) (int, error) {
	if len(dst) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	if len(pkg) == 0 {
		return -1, errors.New("job is missing a package name")
	}
	j := &jobs.Job{
		Type: jobs.DeltaPackage,
		Dst:  dst,
		Pkg:  pkg,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// DeltaPackageExecute carries out a DeltaPackage job
func (m *Manager) DeltaPackageExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Import adds an existing repo to the database
func (m *Manager) Import(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.Import,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// ImportExecute carries out an Import job
func (m *Manager) ImportExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Index generates a new package index
func (m *Manager) Index(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a destination repo")
	}
	j := &jobs.Job{
		Type: jobs.Index,
		Dst:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// IndexExecute carries out an Index job
func (m *Manager) IndexExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Remove deletes a repo from the DB
func (m *Manager) Remove(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.Remove,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// RemoveExecute carries out a Remove job
func (m *Manager) RemoveExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// Rescan rebuild the database for an existing repo
func (m *Manager) Rescan(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.Rescan,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// RescanExecute carries out a Rescan job
func (m *Manager) RescanExecute(j *jobs.Job) error {
	return errors.New("Function not implemented")
}

// TrimPackages removes old package releases and their deltas
func (m *Manager) TrimPackages(name string, max int) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	if max < 1 {
		return -1, errors.New("max releases must be at least 1")
	}
	j := &jobs.Job{
		Type: jobs.TrimPackages,
		Src:  name,
		Max:  max,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// TrimPackagesExecute carries out a TrimPackages job
func (m *Manager) TrimPackagesExecute(j *jobs.Job) error {
	return errors.New("Function not inplemented")
}

// TrimObsoletes removes obsolete packages and their deltas
func (m *Manager) TrimObsoletes(name string) (int, error) {
	if len(name) == 0 {
		return -1, errors.New("job is missing a source repo")
	}
	j := &jobs.Job{
		Type: jobs.TrimObsoletes,
		Src:  name,
	}
	id, err := m.store.Push(j)
	return int(id), err
}

// TrimObsoletesExecute carries out the TrimObsoletes job
func (m *Manager) TrimObsoletesExecute(j *jobs.Job) error {
	return errors.New("Function not inplemented")
}
