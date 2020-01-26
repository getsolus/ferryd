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

package daemon

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

var (
	// ErrDeadLockFile is returned when an dead lockfile was encountered
	ErrDeadLockFile = errors.New("Dead lockfile")

	// ErrOwnedLockFile is returned when the lockfile is already owned by
	// another active process.
	ErrOwnedLockFile = errors.New("File is locked")
)

// A LockFile encapsulates locking functionality
type LockFile struct {
	path string   // Path of the lockfile
	fd   *os.File // Actual file being locked
}

// NewLockFile will return a new lockfile for the given path
func NewLockFile(path string) *LockFile {
	return &LockFile{
		path: path,
	}
}

// Lock will attempt to lock the file, or return an error if this fails
func (l *LockFile) Lock() error {
	exists := true
	// Stat the file
	fi, err := os.Stat(l.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = nil
		exists = false
	}
	pid := -1
	if exists {
		// Check for empty
		if fi.Size() > 0 {
			// Try to open the lock file
			w, err := os.OpenFile(l.path, os.O_RDONLY, 00644)
			if err != nil {
				return err
			}
			l.fd = w
			// Read its contents
			if _, err = fmt.Fscanf(l.fd, "%d", &pid); err != nil {
				return ErrDeadLockFile
			}
			// check if the PID matches the current process
			if pid != os.Getpid() {
				// check if the the PID is still running
				// Process is still active
				// Unix this always works
				p, _ := os.FindProcess(pid)
				if err2 := p.Signal(syscall.Signal(0)); err2 == nil {
					return ErrOwnedLockFile
				}
				// Close file to force overwrite
				l.fd.Close()
				l.fd = nil
			}
		} else {
			return ErrDeadLockFile
		}
	}
	// Create a new lock file
	if l.fd == nil {
		// Try to open the lock file
		w, err := os.OpenFile(l.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 00644)
		if err != nil {
			return err
		}
		l.fd = w
		// write our PID to the file
		pid = os.Getpid()
		fmt.Fprintf(l.fd, "%d", pid)
		err = l.fd.Sync()
	}
	return err
}

// Close will dispose of the lock file
func (l *LockFile) Close() error {
	l.fd.Close()
	return os.Remove(l.path)
}
