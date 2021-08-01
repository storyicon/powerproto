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
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
	"github.com/storyicon/powerproto/pkg/util"

	"github.com/storyicon/powerproto/pkg/util/command"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// ErrGitList defines the git list error
type ErrGitList struct {
	*command.ErrCommandExec
}

// GetGitLatestCommitId is used to get the latest commit id
func GetGitLatestCommitId(ctx context.Context, log logger.Logger, repo string) (string, error) {
	data, err := command.Execute(ctx, log, "", "git", []string{
		"ls-remote", repo, "HEAD",
	}, nil)
	if err != nil {
		return "", &ErrGitList{
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	f := strings.Fields(string(data))
	if len(f) != 2 {
		return "", &ErrGitList{
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	return f[0], nil
}

// ListGitTags is used to list the git tags of specified repository
func ListGitTags(ctx context.Context, log logger.Logger, repo string) ([]string, error) {
	data, err := command.Execute(ctx, log, "", "git", []string{
		"ls-remote", "--tags", "--refs", repo,
	}, nil)
	if err != nil {
		return nil, &ErrGitList{
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	var tags []string
	for _, line := range strings.Split(string(data), "\n") {
		f := strings.Fields(line)
		if len(f) != 2 {
			continue
		}
		if strings.HasPrefix(f[1], "refs/tags/") {
			tags = append(tags, strings.TrimPrefix(f[1], "refs/tags/"))
		}
	}
	malformed, wellFormed := util.SortSemanticVersion(tags)
	return append(malformed, wellFormed...), nil
}

// GithubArchive is github archive
type GithubArchive struct {
	uri       string
	commit    string
	workspace string
}

// GetGithubArchive is used to download github archive
func GetGithubArchive(ctx context.Context, uri string, commitId string) (*GithubArchive, error) {
	filename := fmt.Sprintf("%s.zip", commitId)
	addr := fmt.Sprintf("%s/archive/%s", uri, filename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return nil, &ErrHTTPDownload{
			Url: addr,
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
			Url: addr,
			Err: err,
		}
	}
	zipFilePath := filepath.Join(workspace, filename)
	if err := downloadFile(resp, zipFilePath); err != nil {
		return nil, &ErrHTTPDownload{
			Url:  addr,
			Err:  err,
			Code: resp.StatusCode,
		}
	}
	zip := archiver.NewZip()
	if err := zip.Unarchive(zipFilePath, workspace); err != nil {
		return nil, err
	}
	return &GithubArchive{
		uri:       uri,
		commit:    commitId,
		workspace: workspace,
	}, nil
}

// GetLocalDir is used to get local dir of archive
func (c *GithubArchive) GetLocalDir() string {
	dir := path.Base(c.uri) + "-" + c.commit
	return filepath.Join(c.workspace, dir)
}

// Clear is used to clear the workspace
func (c *GithubArchive) Clear() error {
	return os.RemoveAll(c.workspace)
}
