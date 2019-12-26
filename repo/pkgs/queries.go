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

package pkgs

// PackageSchema is the SQLite3 schema for the Package table
const PackageSchema = `
CREATE TABLE IF NOT EXISTS packages (
    repo_id    INTEGER,
    release_id INTEGER,
    UNIQUE(repo_id,release_id)
)
`

// RepoReleases gets all the releases for a repo
const RepoReleases = `
WITH ids AS (
    SELECT release_id FROM packages
    WHERE repo_id=(SELECT repo_id FROM repos WHERE name=?)
)
SELECT id, package, uri, size, hash, release, from_release, meta FROM releases
INNER JOIN ids ON ids.release_id = releases.id
`

// Insert Query for creating a new Package entry
const Insert = `
INSERT INTO packages (
    repo_id, release_id
) VALUES (
    :repo_id, :release_id
)
`

const (
	// Remove deletes a specific package entry with a repo_id and a release_id
	Remove = "DELETE * FROM packages WHERE repo_id=:repo_id AND release_id=:release_id"
	// RemoveByRepo all package entries for a given repo
	RemoveByRepo = "DELETE * FROM packages WHERE repo_id=:repo_id"
	// RemoveByRelease all package entries for a fiven release
	RemoveByRelease = "DELETE * FROM packages WHERE release_id=:release_id"
)
