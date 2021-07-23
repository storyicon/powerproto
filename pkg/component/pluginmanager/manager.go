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
	"errors"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

var defaultExecuteTimeout = time.Second * 60

// PluginManager is used to manage plugins
type PluginManager interface {
	// GetPluginLatestVersion is used to get the latest version of plugin
	GetPluginLatestVersion(ctx context.Context, path string) (string, error)
	// ListPluginVersions is used to list the versions of plugin
	ListPluginVersions(ctx context.Context, path string) ([]string, error)
	// IsPluginInstalled is used to check whether the plugin is installed
	IsPluginInstalled(ctx context.Context, path string, version string) (bool, string, error)
	// InstallPlugin is used to install plugin
	InstallPlugin(ctx context.Context, path string, version string) (local string, err error)

	// GetGitRepoLatestVersion is used to get the latest version of google apis
	GetGitRepoLatestVersion(ctx context.Context, uri string) (string, error)
	// InstallGitRepo is used to install google apis
	InstallGitRepo(ctx context.Context, uri string, commitId string) (local string, err error)
	// ListGitRepoVersions is used to list protoc version
	ListGitRepoVersions(ctx context.Context, uri string) ([]string, error)
	// IsGitRepoInstalled is used to check whether the protoc is installed
	IsGitRepoInstalled(ctx context.Context, uri string, commitId string) (bool, string, error)
	// GitRepoPath  returns the googleapis path
	GitRepoPath(ctx context.Context, commitId string) string

	// GetProtocLatestVersion is used to get the latest version of protoc
	GetProtocLatestVersion(ctx context.Context) (string, error)
	// ListProtocVersions is used to list protoc version
	ListProtocVersions(ctx context.Context) ([]string, error)
	// IsProtocInstalled is used to check whether the protoc is installed
	IsProtocInstalled(ctx context.Context, version string) (bool, string, error)
	// InstallProtoc is used to install protoc of specified version
	InstallProtoc(ctx context.Context, version string) (local string, err error)
	// IncludePath returns the default include path
	IncludePath(ctx context.Context) string
}

// Config defines the config of PluginManager
type Config struct {
	StorageDir string `json:"storage"`
}

// NewConfig is used to create config
func NewConfig() *Config {
	return &Config{
		StorageDir: consts.GetHomeDir(),
	}
}

// NewPluginManager is used to create PluginManager
func NewPluginManager(cfg *Config, log logger.Logger) (PluginManager, error) {
	return NewBasicPluginManager(cfg.StorageDir, log)
}

// BasicPluginManager is the basic implement of PluginManager
type BasicPluginManager struct {
	logger.Logger
	storageDir   string
	versions     map[string][]string
	versionsLock sync.RWMutex
}

// NewBasicPluginManager is used to create basic PluginManager
func NewBasicPluginManager(storageDir string, log logger.Logger) (*BasicPluginManager, error) {
	return &BasicPluginManager{
		Logger:     log.NewLogger("pluginmanager"),
		storageDir: storageDir,
		versions:   map[string][]string{},
	}, nil
}

// GetPluginLatestVersion is used to get the latest version of plugin
func (b *BasicPluginManager) GetPluginLatestVersion(ctx context.Context, path string) (string, error) {
	versions, err := b.ListPluginVersions(ctx, path)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", errors.New("no version list")
	}
	return versions[len(versions)-1], nil
}

// ListPluginVersions is used to list the versions of plugin
func (b *BasicPluginManager) ListPluginVersions(ctx context.Context, path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultExecuteTimeout)
	defer cancel()

	b.versionsLock.RLock()
	versions, ok := b.versions[path]
	b.versionsLock.RUnlock()
	if ok {
		return versions, nil
	}
	versions, err := ListsGoPackageVersionsAmbiguously(ctx, b.Logger, path)
	if err != nil {
		return nil, err
	}
	b.versionsLock.Lock()
	b.versions[path] = versions
	b.versionsLock.Unlock()
	return versions, nil
}

// IsPluginInstalled is used to check whether the plugin is installed
func (b *BasicPluginManager) IsPluginInstalled(ctx context.Context, path string, version string) (bool, string, error) {
	return IsPluginInstalled(ctx, b.storageDir, path, version)
}

// InstallPlugin is used to install plugin
func (b *BasicPluginManager) InstallPlugin(ctx context.Context, path string, version string) (local string, err error) {
	return InstallPluginUsingGo(ctx, b.Logger, b.storageDir, path, version)
}

// GetGitRepoLatestVersion is used to get the latest version of google apis
func (b *BasicPluginManager) GetGitRepoLatestVersion(ctx context.Context, url string) (string, error) {
	versions, err := b.ListGitRepoVersions(ctx, url)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", errors.New("no version list")
	}
	return versions[len(versions)-1], nil
}

// InstallGitRepo is used to install google apis
func (b *BasicPluginManager) InstallGitRepo(ctx context.Context, uri string, commitId string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultExecuteTimeout)
	defer cancel()
	exists, local, err := b.IsGitRepoInstalled(ctx, uri, commitId)
	if err != nil {
		return "", err
	}
	if exists {
		return local, nil
	}
	release, err := GetGithubArchive(ctx, uri, commitId)
	if err != nil {
		return "", err
	}
	defer release.Clear()

	codePath, err := PathForGitReposCode(b.storageDir, uri, commitId)
	if err != nil {
		return "", err
	}
	if err := util.CopyDirectory(release.GetLocalDir(), codePath); err != nil {
		return "", err
	}
	return local, nil
}

// ListGitRepoVersions is used to list protoc version
func (b *BasicPluginManager) ListGitRepoVersions(ctx context.Context, uri string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultExecuteTimeout)
	defer cancel()

	b.versionsLock.RLock()
	versions, ok := b.versions[uri]
	b.versionsLock.RUnlock()
	if ok {
		return versions, nil
	}
	versions, err := ListGitCommitIds(ctx, b.Logger, uri)
	if err != nil {
		return nil, err
	}
	b.versionsLock.Lock()
	b.versions[uri] = versions
	b.versionsLock.Unlock()
	return versions, nil
}

// IsGitRepoInstalled is used to check whether the protoc is installed
func (b *BasicPluginManager) IsGitRepoInstalled(ctx context.Context, uri string, commitId string) (bool, string, error) {
	codePath, err := PathForGitReposCode(b.storageDir, uri, commitId)
	if err != nil {
		return false, "", err
	}
	exists, err := util.IsDirExists(codePath)
	return exists, PathForGitRepos(b.storageDir, commitId), err
}

// GitRepoPath returns the googleapis path
func (b *BasicPluginManager) GitRepoPath(ctx context.Context, commitId string) string {
	return PathForGitRepos(b.storageDir, commitId)
}

// IsProtocInstalled is used to check whether the protoc is installed
func (b *BasicPluginManager) IsProtocInstalled(ctx context.Context, version string) (bool, string, error) {
	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}
	return IsProtocInstalled(ctx, b.storageDir, version)
}

// GetProtocLatestVersion is used to geet the latest version of protoc
func (b *BasicPluginManager) GetProtocLatestVersion(ctx context.Context) (string, error) {
	versions, err := b.ListProtocVersions(ctx)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", errors.New("no version list")
	}
	return versions[len(versions)-1], nil
}

// ListProtocVersions is used to list protoc version
func (b *BasicPluginManager) ListProtocVersions(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultExecuteTimeout)
	defer cancel()

	b.versionsLock.RLock()
	versions, ok := b.versions["protoc"]
	b.versionsLock.RUnlock()
	if ok {
		return versions, nil
	}
	versions, err := ListGitTags(ctx, b.Logger, consts.ProtobufRepository)
	if err != nil {
		return nil, err
	}
	b.versionsLock.Lock()
	b.versions["protoc"] = versions
	b.versionsLock.Unlock()
	return versions, nil
}

// InstallProtoc is used to install protoc of specified version
func (b *BasicPluginManager) InstallProtoc(ctx context.Context, version string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultExecuteTimeout)
	defer cancel()

	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}
	local := PathForProtoc(b.storageDir, version)
	exists, err := util.IsFileExists(local)
	if err != nil {
		return "", err
	}
	if exists {
		return local, nil
	}

	release, err := GetProtocRelease(ctx, version)
	if err != nil {
		return "", err
	}
	defer release.Clear()
	// merge include files
	includeDir := PathForInclude(b.storageDir)
	if err := util.CopyDirectory(release.GetIncludePath(), includeDir); err != nil {
		return "", err
	}
	// download protoc file
	if err := util.CopyFile(release.GetProtocPath(), local); err != nil {
		return "", err
	}
	// * it is required on unix system
	if err := os.Chmod(local, fs.ModePerm); err != nil {
		return "", err
	}
	return local, nil
}

// IncludePath returns the default include path
func (b *BasicPluginManager) IncludePath(ctx context.Context) string {
	return PathForInclude(b.storageDir)
}
