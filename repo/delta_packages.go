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

// DeltaPackagesSchema is the SQLite3 schema for the DeltaPackages table
const DeltaPackagesSchema = `
CREATE TABLE IF NOT EXISTS delta_packages (
    delta_id    INTEGER,
    package_id INTEGER
)
`

// DeltaPackage is an entry in the DeltaPackages Table
type DeltaPackage struct {
	DeltaID   int `db:"delta_id"`
	PackageID int `db:"package_id"`
}

// Queries for retrieving DeltaPackages of a particular status
const (
	deltaPackages = "SELECT * FROM delta_packages WHERE delta_id=:delta_id"
	packageDeltas = "SELECT * FROM delta_packages WHERE package_id=:package_id"
)

// Query for creating a new DeltaPackages
const insertDeltaPackages = `
INSERT INTO delta_packages (
    delta_id, package_id
) VALUES (
    delta_id, package_id
)
`

// Queries for removing a delta_package
const (
	deleteDeltaPackage     = "DELETE * FROM delta_packages WHERE package_id=:package_id AND delta_id=:delta_id"
	deleteAllDeltasPackage = "DELETE * FROM delta_packages WHERE package_id=:package_id"
)
