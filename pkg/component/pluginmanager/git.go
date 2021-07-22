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
	"strings"

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

// ListGitCommitIds is used to list git commit ids
func ListGitCommitIds(ctx context.Context, log logger.Logger, repo string) ([]string, error) {
	data, err := command.Execute(ctx, log, "", "git", []string{
		"ls-remote", repo,
	}, nil)
	if err != nil {
		return nil, &ErrGitList{
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	var commitIds []string
	for _, line := range strings.Split(string(data), "\n") {
		f := strings.Fields(line)
		if len(f) != 2 {
			continue
		}
		commitIds = append(commitIds, f[0])
	}
	return commitIds, nil
}

// ListGitTags is used to list the git tags of specified repository
func ListGitTags(ctx context.Context, log logger.Logger, repo string) ([]string, error) {
	data, err := command.Execute(ctx, log, "", "git", []string{
		"ls-remote", "--tags", "--refs", "--sort", "version:refname", repo,
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
	return tags, nil
}
