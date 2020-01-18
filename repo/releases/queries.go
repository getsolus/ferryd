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

package releases

// Schema is the SQLite3 schema for the Releases table
const Schema = `
CREATE TABLE IF NOT EXISTS releases (
    id           INTEGER PRIMARY KEY,
    package      STRING,
    uri          STRING,
    size         INTEGER,
    hash         TEXT,
    release      INTEGER,
    from_release INTEGER,
    meta         BLOB
)
`

// Queries for retrieving Releases
const packageReleases = "SELECT * FROM releases WHERE name=:name"

// Insert Query for creating a new Release
const Insert = `
INSERT INTO releases (
    id, package, uri, size, hash, release, from_release, meta
) VALUES (
    NULL, :package, :uri, :size, :hash, :release, :from_release, :meta
)
`

// Update Query for updating a Release
const Update = "UPDATE releases SET size=:size, hash=:hash, meta=:meta WHERE id=:id"

// Queries for removing Releases
const (
	trimPackages  = "DELETE FROM releases WHERE name=:name AND release < :release"
	trimObsoletes = "DELETE FROM releases WHERE name=:name"
	removePackage = "DELETE FROM releases WHERE id=:id"
)
