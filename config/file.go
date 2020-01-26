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

package config

import (
	"encoding/json"
	log "github.com/DataDrake/waterlog"
	"os"
	"path/filepath"
)

const (
	// AssetSuffix for index generation and obsoletion
	AssetSuffix = "assets"
	// DeltaSuffix for delta creation
	DeltaSuffix = "deltas"
	// RepoSuffix for repo storage
	RepoSuffix = "repos"
	// TransitSuffix for incoming packages
	TransitSuffix = "transit"
)

const (
	// DefaultBaseDir for all persistent data
	DefaultBaseDir = "/var/lib/ferryd"
	// DefaultBuildDir for all temporary artifacts
	DefaultBuildDir = "/tmp/ferryd"
	// DefaultLockFile for a running daemon
	DefaultLockFile = "/var/lib/ferryd/ferryd.lock"
	// DefaultSocket for a running daemon
	DefaultSocket = "/run/ferryd.sock"
)

// File contains the file configuration for ferryd
type File struct {
	// BaseDir for all persistent data
	BaseDir string
	// BuildDir for all temporary artifacts
	BuildDir  string
	basePath  []string
	buildPath []string
	// LockFile for the Daemon
	LockFile string
	// Socket for the Daemon
	Socket string
}

// Current is the configuration of the system as it was when the daemon started
var Current *File

// Load reads in a ferryd configuration and validates it
func Load() error {
	// Open File
	cFile, err := os.Open("/etc/ferryd/ferryd.conf")
	if err != nil {
		return err
	}
	defer cFile.Close()
	// Parse JSON
	Current = &File{}
	dec := json.NewDecoder(cFile)
	if err = dec.Decode(Current); err != nil {
		return err
	}
	// Validate Base Directory
	if len(Current.BaseDir) == 0 {
		Current.BaseDir = DefaultBaseDir
		log.Warnf("No BaseDir specified. Using default: %s\n", DefaultBaseDir)
	}
	Current.basePath = filepath.SplitList(Current.BaseDir)
	// Validate Build Directory
	if len(Current.BuildDir) == 0 {
		Current.BuildDir = DefaultBuildDir
		log.Warnf("No BuildDir specified. Using default: %s\n", DefaultBuildDir)
	}
	Current.buildPath = filepath.SplitList(Current.BuildDir)
	// Validate Lock File
	if len(Current.LockFile) == 0 {
		Current.LockFile = DefaultLockFile
		log.Warnf("No Lockfile specified. Using default: %s\n", DefaultLockFile)
	}
	// Validate Socket
	if len(Current.Socket) == 0 {
		Current.Socket = DefaultSocket
		log.Warnf("No Socket specified. Using default: %s\n", DefaultSocket)
	}
	return nil
}

func init() {
	if err := Load(); err != nil {
		panic(err.Error())
	}
}

// AssetPath for index generation and obsoletion
func (f *File) AssetPath() []string {
	return append(f.basePath, AssetSuffix)
}

// DeltaPath for delta creation
func (f *File) DeltaPath() []string {
	return append(f.buildPath, DeltaSuffix)
}

// RepoPath for repo storage
func (f *File) RepoPath() []string {
	return append(f.basePath, RepoSuffix)
}

// TransitPath for incoming packages
func (f *File) TransitPath() []string {
	return append(f.basePath, TransitSuffix)
}
