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

package core

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/getsolus/ferryd/util"
	"github.com/getsolus/libeopkg"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// LinkOrCopyFile is a helper which will initially try to hard link,
// however if we hit an error (because we tried a cross-filesystem hardlink)
// we'll try to copy instead.
func LinkOrCopyFile(source, dest string, forceCopy bool) error {
	if forceCopy {
		return util.CopyFile(source, dest)
	}
	if os.Link(source, dest) == nil {
		return nil
	}
	return util.CopyFile(source, dest)
}

// RemovePackageParents will try to remove the leading components of
// a package file, only if they are empty.
func RemovePackageParents(path string) error {
	sourceDir := filepath.Dir(path)      // i.e. libr/libreoffice
	letterDir := filepath.Dir(sourceDir) // i.e. libr/

	removalPaths := []string{
		sourceDir,
		letterDir,
	}

	for _, p := range removalPaths {
		contents, err := ioutil.ReadDir(p)
		if err != nil {
			return err
		}
		if len(contents) != 0 {
			continue
		}
		if err = os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

func hashFile(path string, h hash.Hash) (sum string, err error) {
	mfile, err := os.Open(path)
	if err != nil {
		return
	}
	defer mfile.Close()
	// Pump from memory into hash for zero-copy sha1sum
	_, err = io.Copy(h, mfile)
	if err != nil {
		return
	}
	sum = hex.EncodeToString(h.Sum(nil))
	return
}

// FileSHA1Sum is a quick wrapper to grab the sha1sum for the given file
func FileSHA1Sum(path string) (string, error) {
	return hashFile(path, sha1.New())
}

// FileSHA256Sum is a quick wrapper to grab the sha256sum for the given file
func FileSHA256Sum(path string) (string, error) {
	return hashFile(path, sha256.New())
}

// WriteSHA1Sum will take the sha1sum of the input path and then dump it to
// the given output path
func WriteSHA1Sum(inpPath, outPath string) error {
	hash, err := FileSHA1Sum(inpPath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outPath, []byte(hash), 00644)
}

// WriteSHA256Sum will take the sha256sum of the input path and then dump it to
// the given output path
func WriteSHA256Sum(inpPath, outPath string) error {
	hash, err := FileSHA256Sum(inpPath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outPath, []byte(hash), 00644)
}

// ProduceDelta will attempt to batch the delta production between the
// two listed file paths and then copy it into the final targetPath
func ProduceDelta(tmpDir, oldPackage, newPackage, targetPath string) error {
	del, err := libeopkg.NewDeltaProducer(tmpDir, oldPackage, newPackage)
	if err != nil {
		return err
	}
	defer del.Close()
	path, err := del.Create()
	if err != nil {
		return err
	}

	// Always nuke the tmpfile
	defer os.Remove(path)

	return LinkOrCopyFile(path, targetPath, false)
}
