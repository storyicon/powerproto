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

package bootstraps

import (
	"context"

	"github.com/pkg/errors"

	"github.com/storyicon/powerproto/pkg/component/configmanager"
	"github.com/storyicon/powerproto/pkg/component/pluginmanager"
	"github.com/storyicon/powerproto/pkg/configs"
	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/logger"
	"github.com/storyicon/powerproto/pkg/util/progressbar"
)

// StepTidyConfig is used to tidy configs by proto file targets
func StepTidyConfig(ctx context.Context, targets []string) error {
	log := logger.NewDefault("tidy")
	log.SetLogLevel(logger.LevelError)

	configManager, err := configmanager.NewConfigManager(log)
	if err != nil {
		return err
	}
	pluginManager, err := pluginmanager.NewPluginManager(pluginmanager.NewConfig(), log)
	if err != nil {
		return err
	}

	configPaths := map[string]struct{}{}
	for _, target := range targets {
		cfg, err := configManager.GetConfig(ctx, target)
		if err != nil {
			return err
		}
		configPaths[cfg.Path()] = struct{}{}
	}
	if len(configPaths) == 0 {
		return nil
	}

	progress := progressbar.GetProgressBar(len(configPaths))
	progress.SetPrefix("tidy configs")
	for path := range configPaths {
		progress.SetSuffix(path)
		err := StepTidyConfigFile(ctx, pluginManager, progress, path)
		if err != nil {
			return err
		}
		progress.Incr()
	}
	progress.SetSuffix("success!")
	progress.Wait()
	return nil
}

// StepTidyConfig is used clean config
// It will amend the 'latest' version to the latest version number in 'vx.y.z' format
func StepTidyConfigFile(ctx context.Context,
	pluginManager pluginmanager.PluginManager,
	progress progressbar.ProgressBar,
	configFilePath string,
) error {
	progress.SetSuffix("load config: %s", configFilePath)
	configItems, err := configs.LoadConfigs(configFilePath)
	if err != nil {
		return err
	}
	var cleanable bool
	for _, item := range configItems {
		if item.Protoc == "latest" {
			progress.SetSuffix("query latest version of protoc")
			version, err := pluginManager.GetProtocLatestVersion(ctx)
			if err != nil {
				return err
			}
			item.Protoc = version
			cleanable = true
		}
		for name, pkg := range item.Plugins {
			path, version, ok := util.SplitGoPackageVersion(pkg)
			if !ok {
				return errors.Errorf("invalid package format: %s, should be in path@version format", pkg)
			}
			if version == "latest" {
				progress.SetSuffix("query latest version of %s", path)
				version, err := pluginManager.GetPluginLatestVersion(ctx, path)
				if err != nil {
					return err
				}
				item.Plugins[name] = util.JoinGoPackageVersion(path, version)
				cleanable = true
			}
		}
	}
	if cleanable {
		progress.SetSuffix("save %s", configFilePath)
		if err := configs.SaveConfigs(configFilePath, configItems...); err != nil {
			return err
		}
	}
	progress.SetSuffix("config file tidied: %s", configFilePath)
	return nil
}
