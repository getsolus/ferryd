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

package jobs

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"ferryd/core"
	"fmt"
	"libferry"
)

// JobEntry is an entry in the JobQueue
type JobEntry struct {
	id         []byte // Unique ID for this job
	sequential bool   // Private to the job implementation
	Type       JobType
	Claimed    bool
	Params     []string
	Timing     libferry.TimingInformation // Store all timing information

	// Not serialised, set by the worker on claim
	description string

	// Not serialised, stored by the worker if the job fails
	failure error
}


// Serialize uses Gob encoding to convert a JobEntry to a byte slice
func (j *JobEntry) Serialize() (result []byte, err error) {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	err = enc.Encode(j)
	if err != nil {
		return
	}
	result = buff.Bytes()
	return
}

// Deserialize use Gob decoding to convert a byte slice to a JobEntry
func Deserialize(serial []byte) (*JobEntry, error) {
	ret := &JobEntry{}
	buff := bytes.NewBuffer(serial)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(ret)
	return ret, err
}

// GetID gets the true numerical ID for this job entry
func (j *JobEntry) GetID() string {
	return fmt.Sprintf("%v", binary.BigEndian.Uint64(j.id))
}
