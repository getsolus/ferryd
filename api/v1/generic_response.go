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

package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/DataDrake/waterlog"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"os"
)

// GenericResponse is a catchall for API responses
type GenericResponse struct {
	// Errors is a list of errors encountered in this transaction
	Errors []string
}

// Print prints out a list of errors, one by one
func (resp GenericResponse) Print(out io.Writer) {
	for _, e := range resp.Errors {
		fmt.Fprintf(os.Stderr, "\t%s\n", e)
	}
}

// sendStockErrors is a utility to send a standard response to the ferry
// client that embeds the error message from ourside.
func (l *Listener) sendStockErrors(errs []error, ctx *fasthttp.RequestCtx) {
	errors := make([]string, len(errs))
	for i, err := range errs {
		errors[i] = err.Error()
		log.Errorf("Client communication error for method '%s', message: '%s'\n", getMethodCaller(), errors[i])
	}
	response := GenericResponse{
		Errors: errors,
	}
	buf := bytes.Buffer{}
	if e2 := json.NewEncoder(&buf).Encode(&response); e2 != nil {
		ctx.Error(e2.Error(), http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusBadRequest)
	ctx.SetBody(buf.Bytes())
}
