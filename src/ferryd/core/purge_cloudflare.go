//
// Copyright © 2017-2023 Solus Project
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

package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

const FerrydDir = "/etc/ferryd"

func purgeCloudflare(target string) error {
	// Check if we have an api token
	apiKey := os.Getenv("CLOUDFLARE_API_TOKEN")
	if len(apiKey) == 0 {
		return nil
	}

	zoneId := os.Getenv("CLOUDFLARE_ZONE_ID")
	if len(zoneId) == 0 {
		return nil
	}

	fileList := filepath.Join(FerrydDir, target)
	if _, err := os.Stat(fileList); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	content, err := os.ReadFile(fileList)
	if err != nil {
		return err
	}

	// We only want non-empty lines
	files := slices.DeleteFunc(strings.Split(string(content), "\n"), func(e string) bool {
		return e == ""
	})

	api := cloudflare.NewClient(option.WithAPIToken(apiKey))

	_, err = api.Cache.Purge(context.Background(), cache.CachePurgeParams{
		ZoneID: cloudflare.F(zoneId),
		Body:   cache.CachePurgeParamsBody{Files: cloudflare.F(any(files))},
	})
	if err != nil {
		return err
	}

	return nil
}
