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

package repo

// Schema is the SQLite3 schema for the Repo table
const Schema = `
CREATE TABLE IF NOT EXISTS repos (
    id              INTEGER PRIMARY KEY,
    name            STRING,
    instant_transit BOOLEAN
)
`

// Queries for retrieving Repos
const (
	// GetSingle retrieves a repo by name
	GetSingle = "SELECT * FROM repos WHERE name=?"
	// GetAll retrieves all repos
	GetAll = "SELECT * FROM repos"
)

// GetSize gives a total size in bytes of the repo
const GetSize = `
WITH ids AS (
    SELECT release_id FROM packages
    WHERE repo_id=?
)
SELECT sum(size) FROM releases
INNER JOIN ids ON ids.release_id = releases.id
`

// PackageCount gets the number of packages in a repo
const PackageCount = `
WITH ids AS (
    SELECT release_id FROM packages
    WHERE repo_id=?
)
SELECT count(*) FROM releases
INNER JOIN ids ON ids.release_id = releases.id
WHERE from_release IS NULL
`

// DeltaCount gets the number of deltas in a repo
const DeltaCount = `
WITH ids AS (
    SELECT release_id FROM packages
    WHERE repo_id=?
)
SELECT count(*) FROM releases
INNER JOIN ids ON ids.release_id = releases.id
WHERE from_release IS NOT NULL
`

// Insert is a Query for creating a new Repo
const Insert = `
INSERT INTO repos (
    id, name, instant_transit
) VALUES (
    NULL, :name, :instant_transit
)
`

// Queries for removing a repo
const (
	Remove = "DELETE FROM repos WHERE id=:id"
)

const reposPackages = `
WITH ids AS (
    SELECT package_id FROM repos
    INNER JOIN repo_packages
    ON repos.id = repo_packages.repo_id
    WHERE repos.name = :name
)
SELECT id, name, uri, size, hash, release, meta FROM packages
INNER JOIN ids
ON packages.id = ids.package_id
`

const sharedPackages = `
WITH ids AS(
    SELECT package_id FROM repos
    INNER JOIN repo_packages
    ON repos.id = repo_packages.repo_id
    WHERE repos.name=:name1 OR repos.name=:name2
    GROUP BY package_id
    HAVING count(*) > 1
)
SELECT id, name, uri, size, hash, release, meta FROM packages
INNER JOIN ids
ON packages.id = ids.package_id;
`

const uniquePackages = `
WITH ids AS(
    SELECT package_id FROM repos
    INNER JOIN repo_packages
    ON repos.id = repo_packages.repo_id
    WHERE repos.name=:name1 OR repos.name=:name2
    GROUP BY package_id
    HAVING count(*) = 1 AND repos.name=:name1
)
SELECT id, name, uri, size, hash, release, meta FROM packages
INNER JOIN ids
ON packages.id = ids.package_id
`
