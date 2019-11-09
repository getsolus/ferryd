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

/*********************/
/* PACKAGE FUNCTIONS */
/*********************/

// TransitPackage processes an incoming package manifest, adding the package to "instant transit" repos
func (m *Manager) TransitPackage(pkg string) (int, error) {
	return -1, errors.New("not yet implemented")
	/* TODO: Finish this
	if len(pkg) == 0 {
		return -1, errors.New("job is missing a package")
	}
	j := &jobs.Job{
		Type: jobs.TransitPackage,
		Pkg:  pkg,
	}
	id, err := m.store.Push(j)
	return int(id), err
	*/
}

// TransitPackageExecute carries out a TransitPackage job
func (m *Manager) TransitPackageExecute(j *jobs.Job) error {
	return errors.New("Function not inplemented")
}
