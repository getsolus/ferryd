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

package libeopkg

import (
	"fmt"
	"os"
	"os/exec"
)

// DISCLAIMER: This stuff is just supporting the existing eopkg stuff.
// We know it's not ideal. When sol comes we'll have a much improved
// package format with a hash-indexed self deduplicating archive to
// mitigate delta issues, reduce sizes, and ensure verification at all stages

// ComputeDeltaName will determine the target name for the delta eopkg
func ComputeDeltaName(oldPackage, newPackage *MetaPackage) string {
	return fmt.Sprintf("%s-%d-%d-%s-%s.delta.eopkg",
		newPackage.Name,
		oldPackage.GetRelease(),
		newPackage.GetRelease(),
		newPackage.DistributionRelease,
		newPackage.Architecture)
}

// IsDeltaPossible will compare the two input packages and determine if it
// is possible for a delta to be considered. Note that we do not compare the
// distribution _name_ because Solus already had to do a rename once, and that
// broke delta updates. Let's not do that again. eopkg should in reality determine
// delta applicability based on repo origin + upgrade path, not names
func IsDeltaPossible(oldPackage, newPackage *MetaPackage) bool {
	return oldPackage.GetRelease() < newPackage.GetRelease() &&
		oldPackage.Name == newPackage.Name &&
		oldPackage.DistributionRelease == newPackage.DistributionRelease &&
		oldPackage.Architecture == newPackage.Architecture
}

// XzFile is a simple wrapper around the xz utility to compress the input
// file. This will be performed in place and leave a ".xz" suffixed file in
// place
// Keep original determines whether we'll keep the original file
func XzFile(inputPath string, keepOriginal bool) error {
	cmd := []string{
		"xz",
		"-6",
		"-T", "8",
		inputPath,
	}
	if keepOriginal {
		cmd = append(cmd, "-k")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stderr = os.Stderr
	return c.Run()
}

// UnxzFile will decompress the input XZ file and leave a new file in place
// without the .xz suffix
func UnxzFile(inputPath string, keepOriginal bool) error {
	cmd := []string{
		"unxz",
		"-T", "8",
		inputPath,
	}
	if keepOriginal {
		cmd = append(cmd, "-k")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stderr = os.Stderr
	return c.Run()
}

// ZstdFile is a simple wrapper around the zstd utility to compress the input
// file. This will be performed in place and leave a ".zst" suffixed file in
// place
// Keep original determines whether we'll keep the original file
func ZstdFile(inputPath string, keepOriginal bool) error {
	cmd := []string{
		"zstd",
		"-3",
		"-T8",
		inputPath,
	}
	if keepOriginal {
		cmd = append(cmd, "-k")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stderr = os.Stderr
	return c.Run()
}

// UnzstdFile will decompress the input zstd file and leave a new file in place
// without the .zst suffix
func UnzstdFile(inputPath string, keepOriginal bool) error {
	cmd := []string{
		"unzstd",
		"-T8",
		inputPath,
	}
	if keepOriginal {
		cmd = append(cmd, "-k")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stderr = os.Stderr
	return c.Run()
}
