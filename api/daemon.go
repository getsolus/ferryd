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

package v1

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

// ModifyDaemon makes a requested change to the daemon
func (l *Listener) ModifyDaemon(ctx *fasthttp.RequestCtx) {
	// Get the "action" query argument
	action := string(ctx.QueryArgs().Peek("action"))
	if len(action) == 0 {
		writeErrorString(ctx, "Action not specified when modifying daemon", http.StatusBadRequest)
		return
	}
	// Pivot by action
	switch action {
	default:
		writeErrorString(ctx, fmt.Sprintf("Action '%s' not implemented for the daemon", action), http.StatusBadRequest)
	}
	return
}
