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

// PackageSchema is the SQLite3 schema for the Package table
const PackageSchema = `
CREATE TABLE IF NOT EXISTS packages (
    id      INTEGER PRIMARY KEY,
    name    STRING,
    uri     STRING,
    size    INTEGER,
    hash    TEXT,
    release INTEGER,
    meta    BLOB
)
`

// Queries for retrieving Packages
const (
	namedPackage = "SELECT * FROM packages WHERE name=:name"
)

// Query for creating a new Package
const insertPackage = `
INSERT INTO packages (
    id, name, uri, size, hash, release, meta
) VALUES (
    NULL, :name, :uri, :size, :hash, :release, :meta
)
`

// Queries for updating a Package
const (
	updatePackage = "UPDATE packages SET size=:size, hash=:hash, meta=:meta WHERE id=:id"
)

// Queries for removing a Package
const (
	trimPackages  = "DELETE FROM packages WHERE name=:name AND release < :release"
	trimObsoletes = "DELETE FROM packages WHERE name=:name"
	removePackage = "DELETE FROM packages WHERE id=:id"
)
