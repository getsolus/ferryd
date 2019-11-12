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

package repo

// RepoSchema is the SQLite3 schema for the Repo table
const RepoSchema = `
CREATE TABLE IF NOT EXISTS repos (
    id              INTEGER PRIMARY KEY,
    name            STRING,
    instant_transit BOOLEAN
)
`

// Queries for retrieving Repo of a particular status
const (
	primaryRepo  = "SELECT * FROM repos WHERE id=1"
	allRepos     = "SELECT * FROM repos"
	instantRepos = "SELECT * FROM repos WHERE instant_transit=TRUE"
)

// Query for creating a new Repo
const insertRepo = `
INSERT INTO repos (
    id, name, instant_transit
) VALUES (
    NULL, :name, :instant_transit
)
`

// Queries for updating a repo
const (
	updateRepo = "UPDATE repos SET name=:name, instant_repo=:instant_repo WHERE id=:id"
)

// Queries for removing a repo
const (
	removeRepo = "DELETE FROM jobs WHERE id=:id"
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
