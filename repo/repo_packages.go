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

// RepoPackage is an entry in the RepoPackages Table
type RepoPackage struct {
	RepoID    int `db:"repo_id"`
	PackageID int `db:"package_id"`
}

// RepoPackagesSchema is the SQLite3 schema for the RepoPackages table
const RepoPackagesSchema = `
CREATE TABLE IF NOT EXISTS repo_packages (
    repo_id    INTEGER,
    package_id INTEGER
)
`

// Queries for retrieving RepoPackages of a particular status
const (
	repoPackages = "SELECT * FROM repo_packages WHERE repo_id=:repo_id"
	packageRepos = "SELECT * FROM repo_packages WHERE package_id=:package_id"
)

// Query for creating a new RepoPackages
const insertRepoPackages = `
INSERT INTO repo_packages (
    repo_id, package_id
) VALUES (
    repo_id, package_id
)
`

// Queries for removing a repo_package
const (
	deleteRepoPackage     = "DELETE * FROM repo_packages WHERE package_id=:package_id AND repo_id=:repo_id"
	deleteAllPackagesRepo = "DELETE * FROM repo_packages WHERE repo_id=:repo_id"
	deleteAllReposPackage = "DELETE * FROM repo_packages WHERE package_id=:package_id"
)
