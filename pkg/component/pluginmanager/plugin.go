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
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"

	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/command"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// IsPluginInstalled is used to check whether a plugin is installed
func IsPluginInstalled(ctx context.Context,
	storageDir string,
	path string, version string) (bool, string, error) {
	local, err := PathForPlugin(storageDir, path, version)
	if err != nil {
		return false, "", err
	}
	exists, err := util.IsFileExists(local)
	if err != nil {
		return false, "", err
	}
	if exists {
		return true, local, nil
	}
	return false, "", nil
}

// InstallPluginUsingGo is used to install plugin using golang
func InstallPluginUsingGo(ctx context.Context,
	log logger.Logger,
	storageDir string,
	path string, version string) (string, error) {
	exists, local, err := IsPluginInstalled(ctx, storageDir, path, version)
	if err != nil {
		return "", err
	}
	if exists {
		return local, nil
	}

	local, err = PathForPlugin(storageDir, path, version)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(local)

	uri := util.JoinGoPackageVersion(path, version)
	_, err2 := command.Execute(ctx, log, "", "go", []string{
		"install", uri,
	}, []string{"GOBIN=" + dir, "GO111MODULE=on"})
	if err2 != nil {
		return "", &ErrGoInstall{
			ErrCommandExec: err2,
		}
	}
	return local, nil
}

// ///////////////// Version Control /////////////////

// Module defines the model of go list data
type Module struct {
	Path      string       // module path
	Version   string       // module version
	Versions  []string     // available module versions (with -versions)
	Replace   *Module      // replaced by this module
	Time      *time.Time   // time version was created
	Update    *Module      // available update, if any (with -u)
	Main      bool         // is this the main module?
	Indirect  bool         // is this module only an indirect dependency of main module?
	Dir       string       // directory holding files for this module, if any
	GoMod     string       // path to go.mod file used when loading this module, if any
	GoVersion string       // go version used in module
	Retracted string       // retraction information, if any (with -retracted or -u)
	Error     *ModuleError // error loading module
}

// ModuleError defines the module error
type ModuleError struct {
	Err string // the error itself
}

// ListGoPackageVersions is list go package versions
func ListGoPackageVersions(ctx context.Context, log logger.Logger, path string) ([]string, error) {
	// query from latest version
	// If latest is not specified here, the queried version
	// may be restricted to the current project go.mod/go.sum
	pkg := util.JoinGoPackageVersion(path, "latest")
	data, err := command.Execute(ctx, log, "", "go", []string{
		"list", "-m", "-json", "-versions", pkg,
	}, []string{
		"GO111MODULE=on",
	})
	if err != nil {
		return nil, &ErrGoList{
			ErrCommandExec: err,
		}
	}
	var module Module
	if err := jsoniter.Unmarshal(data, &module); err != nil {
		return nil, err
	}
	if len(module.Versions) != 0 {
		return module.Versions, nil
	}
	return []string{module.Version}, nil
}

// ListsGoPackageVersionsAmbiguously is used to list go package versions ambiguously
func ListsGoPackageVersionsAmbiguously(ctx context.Context, log logger.Logger, pkg string) ([]string, error) {
	type Result struct {
		err      error
		pkg      string
		versions []string
	}
	items := strings.Split(pkg, "/")
	dataMap := make([]*Result, len(items))
	notify := make(chan struct{}, 1)
	maxIndex := len(items) - 1
	for i := maxIndex; i >= 1; i-- {
		go func(i int) {
			pkg := strings.Join(items[0:i+1], "/")
			versions, err := ListGoPackageVersions(context.TODO(), log, pkg)
			dataMap[maxIndex-i] = &Result{
				pkg:      pkg,
				versions: versions,
				err:      err,
			}
			notify <- struct{}{}
		}(i)
	}
OutLoop:
	for {
		select {
		case <-notify:
			var errs error
			for _, data := range dataMap {
				if data == nil {
					continue OutLoop
				}
				if data.err != nil {
					errs = multierror.Append(errs, data.err)
				}
				if data.versions != nil {
					return data.versions, nil
				}
			}
			return nil, errs
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
