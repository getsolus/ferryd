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

import (
	"github.com/getsolus/ferryd/config"
	"github.com/getsolus/ferryd/repo/archive"
	"github.com/getsolus/ferryd/repo/pkgs"
	"github.com/getsolus/ferryd/util"
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

const (
	// DB is the filename of the jobs database
	DB = "repos.db"
	// SQLiteOpts is a list of options for the go-sqlite3 driver
	SQLiteOpts = "?cache=shared"
)

// OpenDB opens a connection to the DB and creates missing tables
func OpenDB() *sqlx.DB {
	// Open the DB
	db, err := sqlx.Open("sqlite3", filepath.Join(config.Current.BaseDir, DB)+SQLiteOpts)
	if err != nil {
		panic(err.Error())
	}
	// See: https://github.com/mattn/go-sqlite3/issues/209
	db.SetMaxOpenConns(1)
	// Create repo tables if missing
	db.MustExec(Schema)
	db.MustExec(pkgs.Schema)
	db.MustExec(archive.Schema)
	// Check that the repos directory exists
	if err = util.CreateDir(config.Current.RepoPath()); err != nil {
		panic(err.Error())
	}
	// Check that the assets directory exists
	if err = util.CreateDir(config.Current.AssetPath()); err != nil {
		panic(err.Error())
	}
	return db
}
