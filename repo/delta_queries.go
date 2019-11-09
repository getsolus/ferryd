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

// DeltaSchema is the SQLite3 schema for the Delta table
const DeltaSchema = `
CREATE TABLE IF NOT EXISTS deltas (
    id           INTEGER PRIMARY KEY,
    package_name STRING,
    uri          STRING,
    size         INTEGER,
    hash         TEXT,
    from_rel     INTEGER,
    to_rel       INTEGER,
)
`

// Queries for retrieving Deltas
const (
	namedDeltas = "SELECT * FROM deltas WHERE package_name=:package_name"
)

// Query for creating a new Delta
const insertDelta = `
INSERT INTO deltas (
    id, package_name, uri, size, hash, from_release, to_release
) VALUES (
    NULL, :package_name, :uri, :size, :hash, :from_release, :to_release
)
`

// Queries for removing a Delta
const (
	trimDeltas     = "DELETE FROM deltas WHERE package_name=:package_name AND to_release < :to_release"
	obsoleteDeltas = "DELETE FROM deltas WHERE package_name=:package_name"
	removeDelta    = "DELETE FROM deltas WHERE package_name=:package_name AND to_release = :to_release"
)
