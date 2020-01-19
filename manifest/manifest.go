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

package manifest

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/getsolus/ferryd/core"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	// Suffix is the extension that a valid transit manifest must have
	Suffix = ".tram"
)

var (
	// ErrInvalidHeader will be returned when the [manifest] section is malformed
	ErrInvalidHeader = errors.New("Manifest contains an invalid header")

	// ErrMissingTarget will be returned when the target is not present
	ErrMissingTarget = errors.New("Manifest contains no target")

	// ErrMissingPayload will be returned when the [[file]]s are missing
	ErrMissingPayload = errors.New("Manifest does not contain a payload")

	// ErrInvalidPayload will be returned when the payload is in some way invalid
	ErrInvalidPayload = errors.New("Manifest contains an invalid payload")

	// ErrIllegalUpload is returned when someone is a spanner and tries uploading an unsupported file
	ErrIllegalUpload = errors.New("The manifest file is NOT an eopkg")
)

// A Header is required in all .tram uploads to ensure that both
// the sender and recipient are talking in the same fashion.
type Header struct {
	// Versioning to protect against future format changes
	Version string `toml:"version"`

	// The repo that the uploader is intending to upload *to*
	Target string `toml:"target"`
}

// A Manifest is provided by build servers to validate the upload of
// packages into the incoming directory.
//
// This is to ensure all uploads are intentional, complete and verifiable.
type Manifest struct {

	// Every .tram file has a [manifest] header - this will never change and is
	// version agnostic.
	Head Header `toml:"manifest"`

	// A list of files that accompanied this .tram upload
	File []File `toml:"file"`

	Path string // Privately held path to the file
	dir  string // Where the .tram was loaded from
	id   string // Effectively our basename
}

// ID will return the unique ID for the transit manifest file
func (t *Manifest) ID() string {
	return t.id
}

// GetPaths will return the package paths as a slice of strings
func (t *Manifest) GetPaths() []string {
	var ret []string
	for i := range t.File {
		f := &t.File[i]
		ret = append(ret, filepath.Join(t.dir, f.Path))
	}
	return ret
}

// File provides simple verification data for each file in the
// uploaded payload.
type File struct {

	// Relative filename, i.e. nano-2.7.5-68-1-x86_64.eopkg
	Path string `toml:"path"`

	// Cryptographic checksum to allow integrity checks post-upload/pre-merge
	Sha256 string `toml:"sha256"`
}

// NewManifest will attempt to load the transit manifest from the
// named path and perform *basic* validation.
func NewManifest(path string) (*Manifest, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	ret := &Manifest{
		Path: abs,
		dir:  filepath.Dir(abs),
		id:   filepath.Base(abs),
	}

	blob, err := ioutil.ReadFile(ret.Path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(blob), ret); err != nil {
		return nil, err
	}

	ret.Head.Target = strings.TrimSpace(ret.Head.Target)
	ret.Head.Version = strings.TrimSpace(ret.Head.Version)

	if ret.Head.Version != "1.0" {
		return nil, ErrInvalidHeader
	}

	if len(ret.Head.Target) < 1 {
		return nil, ErrMissingTarget
	}

	if len(ret.File) < 1 {
		return nil, ErrMissingPayload
	}

	for i := range ret.File {
		f := &ret.File[i]
		f.Path = strings.TrimSpace(f.Path)
		f.Sha256 = strings.TrimSpace(f.Sha256)

		if len(f.Path) < 1 || len(f.Sha256) < 1 {
			return nil, ErrInvalidPayload
		}

		if !strings.HasSuffix(f.Path, ".eopkg") {
			return nil, ErrIllegalUpload
		}
	}

	return ret, nil
}

// Verify will verify the files listed in the manifest locally, ensuring
// that they actually exist, and that the hashes match to prevent any corrupted
// uploads being inadvertently imported
func (t *Manifest) Verify() error {
	for i := range t.File {
		f := &t.File[i]
		path := filepath.Join(t.dir, f.Path)
		sha, err := core.FileSHA256Sum(path)
		if err != nil {
			return err
		}
		if sha != f.Sha256 {
			return fmt.Errorf("Invalid SHA256 for '%s'. Local: '%s'", f.Path, sha)
		}
	}
	return nil
}
