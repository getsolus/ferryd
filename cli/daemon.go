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

package cli

import (
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
	log "github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/getsolus/ferryd/config"
	"github.com/getsolus/ferryd/daemon"
	log2 "log"
	"os"
	"path/filepath"
)

// Daemon fulfills the "daemon" sub-command
var Daemon = &cmd.CMD{
	Name:  "daemon",
	Alias: "up",
	Short: "Start a new ferryd daemon",
	Args:  &DaemonArgs{},
	Run:   DaemonRun,
}

// DaemonArgs are the arguments to the "daemon" sub-command
type DaemonArgs struct{}

// DaemonRun executes the "daemon" sub-command
func DaemonRun(r *cmd.RootCMD, c *cmd.CMD) {
	//flags := r.Flags.(*GlobalFlags)
	//args  := c.Args.(*DaemonArgs)
	// Set up the logger
	log.SetFormat(format.Partial)
	log.SetFlags(log2.Ltime | log2.Ldate | log2.LUTC)
	log.SetLevel(level.Debug)
	// Make sure BaseDir exists
	if _, err := os.Stat(config.Current.BaseDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Base directory does not exist: %s\n", config.Current.BaseDir)
		os.Exit(1)
	}
	// Open the log file
	logPath := filepath.Join(config.Current.LogPath()...)
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 00644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %s %v\n", logPath, err)
		os.Exit(1)
	}
	defer logFile.Close()
	// Set output of the logger to the log file
	log.SetOutput(logFile)
	// Need to get a lock file before we can even grab the log file
	log.Infoln("Initialising server")
	srv := daemon.NewServer()
	defer srv.Close()
	// Bind to the socket
	log.Infoln("Binding to API socket")
	if err := srv.Bind(); err != nil {
		log.Errorf("Error in binding server socket '%s', message: '%s'\n", config.Current.Socket, err.Error())
		return
	}
	// Start serving
	log.Infoln("Starting API endpoints")
	if err := srv.Serve(); err != nil {
		log.Errorf("Error in serving on socket '%s', message: '%s'\n", config.Current.Socket, err.Error())
		return
	}
}
