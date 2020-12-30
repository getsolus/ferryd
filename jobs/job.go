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

package jobs

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// Job is an entry in the Job Table
type Job struct {
	ID   int     `db:"id"`
	Type JobType `db:"type"`
	// Job-specific arguments
	Src string `db:"src"`
	Dst string `db:"dst"`
	Pkg string `db:"pkg"`
	Max int    `db:"max"`
	// Job tracking
	Created  sql.NullTime   `db:"created"`
	Started  sql.NullTime   `db:"started"`
	Finished sql.NullTime   `db:"finished"`
	Status   JobStatus      `db:"status"`
	Message  sql.NullString `db:"message"`
	Results  []byte         `db:"results"`
}

// RunningSince will return the job has been running
func (j *Job) RunningSince() time.Duration {
	return time.Now().UTC().Sub(j.Started.Time)
}

// RunTime will return the time it took to execute a job (used when a job fails)
func (j *Job) RunTime() time.Duration {
	return j.Finished.Time.Sub(j.Started.Time)
}

// QueuedSince will let us know how long this task has been queued
func (j *Job) QueuedSince() time.Duration {
	return time.Now().UTC().Sub(j.Created.Time)
}

// QueuedTime will return the total time that the job was queued for
func (j *Job) QueuedTime() time.Duration {
	return j.Started.Time.Sub(j.Created.Time)
}

// TotalTime will return the total time a job took to complete from queuing
func (j *Job) TotalTime() time.Duration {
	return j.Finished.Time.Sub(j.Created.Time)
}

// Describe generates a short description of what a package does
func (j *Job) Describe() string {
	switch j.Type {
	case Check:
		return fmt.Sprintf("Comparing DB with Disk for repository '%s'", j.Src)
	case CherryPick:
		return fmt.Sprintf("Cherry-picking '%s' from '%s' to '%s'", j.Pkg, j.Src, j.Dst)
	case Clone:
		return fmt.Sprintf("Cloning repo from '%s' to '%s'", j.Src, j.Dst)
	case Compare:
		return fmt.Sprintf("Comparing '%s' with '%s'", j.Src, j.Dst)
	case Create:
		return fmt.Sprintf("Creating new repo '%s'", j.Dst)
	case Delta:
		return fmt.Sprintf("Generating Deltas for repo '%s'", j.Dst)
	case DeltaPackage:
		return fmt.Sprintf("Generating Deltas for '%s' in repo '%s'", j.Pkg, j.Dst)
	case Import:
		return fmt.Sprintf("Importing existing repo '%s'", j.Src)
	case Index:
		return fmt.Sprintf("Generating Index for repo '%s'", j.Dst)
	case Remove:
		return fmt.Sprintf("Removing repo '%s' from DB", j.Src)
	case Rescan:
		return fmt.Sprintf("Re-scanning repo '%s'", j.Src)
	case Sync:
		return fmt.Sprintf("Syncing from '%s' to '%s'", j.Src, j.Dst)
	case TrimObsoletes:
		return fmt.Sprintf("Trimming Obsolete packages from repo '%s'", j.Src)
	case TrimPackages:
		return fmt.Sprintf("Trimming old releases (max: %d) in repo '%s'", j.Max, j.Src)
	case TransitPackage:
		return fmt.Sprintf("Transiting new package '%s' to '%s'", j.Pkg, j.Src)
	default:
		return "Unsupported Job Type"
	}
}

// Print writes out a human-readable version of a job
func (j *Job) Print() {
	fmt.Printf("ID:   %d\n", j.ID)
	fmt.Printf("Type: %s\n", typeMap[j.Type])
	fmt.Println("Arguments:")
	none := true
	if len(j.Src) > 0 {
		fmt.Printf("\tSource:  %s\n", j.Src)
		none = false
	}
	if len(j.Dst) > 0 {
		fmt.Printf("\tDest:    %s\n", j.Dst)
		none = false
	}
	if len(j.Pkg) > 0 {
		fmt.Printf("\tPackage: %s\n", j.Pkg)
		none = false
	}
	if j.Max != 0 {
		fmt.Printf("\tMax:     %d\n", j.Max)
		none = false
	}
	if none {
		fmt.Println("\tNone.")
	}
	fmt.Println("Times:")
	if j.Created.Valid && !j.Created.Time.IsZero() {
		fmt.Printf("\tCreated:  %s\n", j.Created.Time.Format(time.RFC3339))
	}
	if j.Started.Valid && !j.Started.Time.IsZero() {
		fmt.Printf("\tStarted:  %s\n", j.Started.Time.Format(time.RFC3339))
		fmt.Printf("\t\tQueued:   %s\n", j.QueuedTime().String())
	}
	if j.Finished.Valid && !j.Finished.Time.IsZero() {
		fmt.Printf("\tFinished: %s\n", j.Finished.Time.Format(time.RFC3339))
		fmt.Printf("\t\tRuntime: %s\n", j.RunTime().String())
	}
	if j.Status > Running {
		fmt.Printf("\tTotal:      %s\n", j.TotalTime().String())
	}
	fmt.Printf("Status: %s\n", statusMap[j.Status])
	if j.Message.Valid && len(j.Message.String) > 0 {
		fmt.Printf("Last Message: %s\n", j.Message.String)
	}
	if l := len(j.Results); l > 0 {
		fmt.Printf("Results:     %dB\n", l)
	}
}

// Create adds this job to the DB
func (j *Job) Create(tx *sqlx.Tx) error {
	// Create the record
	res, err := tx.NamedExec(Insert, j)
	if err != nil {
		return err
	}
	// Get the ID of the new record
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	j.ID = int(id)
	return nil
}

// Save updates this job in the DB
func (j *Job) Save(tx *sqlx.Tx) error {
	_, err := tx.NamedExec(Update, j)
	return err
}
