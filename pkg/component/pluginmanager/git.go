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
	"os"
	"path"
	"path/filepath"
	"strings"

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

// GitRepository describes a local git repository
type GitRepository struct {
	uri       string
	commit    string
	workspace string
}

func GetGitRepository(ctx context.Context, uri string, commitId string, log logger.Logger) (*GitRepository, error) {
	workspace, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}

	repo := GitRepository{
		uri:       uri,
		commit:    commitId,
		workspace: workspace,
	}

	// clone via the git command instead of using for example "go-git" so that the authentication is not our problem
	_, err = command.Execute(ctx, log, workspace, "git", []string{
		"clone", uri, repo.GetLocalDir(),
	}, nil)
	if err != nil {
		return nil, err
	}

	_, err = command.Execute(ctx, log, repo.GetLocalDir(), "git", []string{
		"reset", "--hard", commitId,
	}, nil)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

// GetLocalDir is used to get local dir of repo
func (r *GitRepository) GetLocalDir() string {
	dir := path.Base(r.uri) + "-" + r.commit
	return filepath.Join(r.workspace, dir)
}

// Clear is used to clear the workspace
func (r *GitRepository) Clear() error {
	return os.RemoveAll(r.workspace)
}
