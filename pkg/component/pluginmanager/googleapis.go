// Copyright 2021 storyicon@foxmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginmanager

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

// GoogleAPIRelease defines the release of google api
type GoogleAPIRelease struct {
	workspace string
	commit    string
}

// GetDir is used to get the dir of google api
func (p *GoogleAPIRelease) GetDir() string {
	return filepath.Join(p.workspace, "googleapis-"+p.commit)
}

// Clear is used to clear the workspace
func (p *GoogleAPIRelease) Clear() error {
	return os.RemoveAll(p.workspace)
}

// GetGoogleAPIRelease is used to get the release of google api
func GetGoogleAPIRelease(ctx context.Context, commitId string) (*GoogleAPIRelease, error) {
	filename := fmt.Sprintf("%s.zip", commitId)
	uri := fmt.Sprintf("https://github.com/googleapis/googleapis/archive/%s", filename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, &ErrHTTPDownload{
			Url: uri,
			Err: err,
		}
	}
	workspace, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &ErrHTTPDownload{
			Url: uri,
			Err: err,
		}
	}
	zipFilePath := filepath.Join(workspace, filename)
	if err := downloadFile(resp, zipFilePath); err != nil {
		return nil, &ErrHTTPDownload{
			Url:  uri,
			Err:  err,
			Code: resp.StatusCode,
		}
	}
	zip := archiver.NewZip()
	if err := zip.Unarchive(zipFilePath, workspace); err != nil {
		return nil, err
	}
	return &GoogleAPIRelease{
		commit:    commitId,
		workspace: workspace,
	}, nil
}
