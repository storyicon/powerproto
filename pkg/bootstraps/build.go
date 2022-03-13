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
	"fmt"

	"github.com/pkg/errors"

	"github.com/storyicon/powerproto/pkg/component/actionmanager"
	"github.com/storyicon/powerproto/pkg/component/compilermanager"
	"github.com/storyicon/powerproto/pkg/component/configmanager"
	"github.com/storyicon/powerproto/pkg/component/pluginmanager"
	"github.com/storyicon/powerproto/pkg/configs"
	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/concurrent"
	"github.com/storyicon/powerproto/pkg/util/logger"
	"github.com/storyicon/powerproto/pkg/util/progressbar"
)

// StepLookUpConfigs is used to lookup config files according to target proto files
func StepLookUpConfigs(
	ctx context.Context,
	targets []string,
	configManager configmanager.ConfigManager,
) ([]configs.ConfigItem, error) {
	progress := progressbar.GetProgressBar(ctx, len(targets))
	progress.SetPrefix("Lookup configs of proto files")
	var configItems []configs.ConfigItem
	deduplicate := map[string]struct{}{}
	for _, target := range targets {
		cfg, err := configManager.GetConfig(ctx, target)
		if err != nil {
			return nil, err
		}
		if _, exists := deduplicate[cfg.ID()]; !exists {
			configItems = append(configItems, cfg)
			deduplicate[cfg.ID()] = struct{}{}
		}
		progress.SetSuffix("load %s", cfg.ID())
		progress.Incr()
	}
	progress.SetSuffix("success!")
	progress.Wait()
	fmt.Printf("the following %d configurations will be used: \r\n", len(configItems))
	for _, config := range configItems {
		fmt.Printf("	%s \r\n", config.Path())
	}
	if len(configItems) == 0 {
		return nil, errors.New("no config file matched, please check the scope of config file " +
			"or use 'powerproto init' to create config file")
	}
	return configItems, nil
}

// StepInstallRepositories is used to install repositories
func StepInstallRepositories(ctx context.Context,
	pluginManager pluginmanager.PluginManager,
	configItems []configs.ConfigItem) error {
	deduplicate := map[string]struct{}{}
	for _, config := range configItems {
		for _, pkg := range config.Config().Repositories {
			deduplicate[pkg] = struct{}{}
		}
	}
	if len(deduplicate) == 0 {
		return nil
	}
	progress := progressbar.GetProgressBar(ctx, len(deduplicate))
	progress.SetPrefix("Install repositories")

	repoMap := map[string]struct{}{}
	for pkg := range deduplicate {
		path, version, ok := util.SplitGoPackageVersion(pkg)
		if !ok {
			return errors.Errorf("invalid format: %s, should in path@version format", pkg)
		}
		if version == "latest" {
			progress.SetSuffix("query latest version of %s", path)
			latestVersion, err := pluginManager.GetGitRepoLatestVersion(ctx, path)
			if err != nil {
				return errors.Wrapf(err, "failed to query latest version of %s", path)
			}
			version = latestVersion
			pkg = util.JoinGoPackageVersion(path, version)
		}
		progress.SetSuffix("check cache of %s", pkg)
		exists, _, err := pluginManager.IsGitRepoInstalled(ctx, path, version)
		if err != nil {
			return err
		}
		if exists {
			progress.SetSuffix("the %s version of %s is already cached", version, path)
		} else {
			progress.SetSuffix("install %s version of %s", version, path)
			_, err = pluginManager.InstallGitRepo(ctx, path, version)
			if err != nil {
				return err
			}
			progress.SetSuffix("the %s version of %s is installed", version, path)
		}
		repoMap[pkg] = struct{}{}
		progress.Incr()
	}
	progress.SetSuffix("all repositories have been installed")
	progress.Wait()
	fmt.Println("the following versions of googleapis will be used:")
	for pkg := range repoMap {
		fmt.Printf("	%s\r\n", pkg)
	}
	return nil
}

// StepInstallProtoc is used to install protoc
func StepInstallProtoc(ctx context.Context,
	pluginManager pluginmanager.PluginManager,
	configItems []configs.ConfigItem) error {
	deduplicate := map[string]struct{}{}
	for _, config := range configItems {
		version := config.Config().Protoc
		if version == "" {
			return errors.Errorf("protoc version is required: %s", config.Path())
		}
		deduplicate[version] = struct{}{}
	}
	progress := progressbar.GetProgressBar(ctx, len(deduplicate))
	progress.SetPrefix("Install protoc")

	versionsMap := map[string]struct{}{}
	for version := range deduplicate {
		if version == "latest" {
			progress.SetSuffix("query latest version of protoc")
			latestVersion, err := pluginManager.GetProtocLatestVersion(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to list protoc versions")
			}
			version = latestVersion
		}
		progress.SetSuffix("check cache of protoc %s", version)
		exists, _, err := pluginManager.IsProtocInstalled(ctx, version)
		if err != nil {
			return err
		}
		if exists {
			progress.SetSuffix("the %s version of protoc is already cached", version)
		} else {
			progress.SetSuffix("install %s version of protoc", version)
			_, err = pluginManager.InstallProtoc(ctx, version)
			if err != nil {
				return err
			}
			progress.SetSuffix("the %s version of protoc is installed", version)
		}
		versionsMap[version] = struct{}{}
		progress.Incr()
	}
	progress.Wait()
	fmt.Println("the following versions of protoc will be used:", util.SetToSlice(versionsMap))
	return nil
}

// StepInstallPlugins is used to install plugins
func StepInstallPlugins(ctx context.Context,
	pluginManager pluginmanager.PluginManager,
	configItems []configs.ConfigItem,
) error {
	deduplicate := map[string]struct{}{}
	for _, config := range configItems {
		for _, pkg := range config.Config().Plugins {
			deduplicate[pkg] = struct{}{}
		}
	}
	progress := progressbar.GetProgressBar(ctx, len(deduplicate))
	progress.SetPrefix("Install plugins")
	pluginsMap := map[string]struct{}{}
	for pkg := range deduplicate {
		path, version, ok := util.SplitGoPackageVersion(pkg)
		if !ok {
			return errors.Errorf("invalid format: %s, should in path@version format", pkg)
		}
		if version == "latest" {
			progress.SetSuffix("query latest version of %s", path)
			data, err := pluginManager.GetPluginLatestVersion(ctx, path)
			if err != nil {
				return err
			}
			version = data
			pkg = util.JoinGoPackageVersion(path, version)
		}
		progress.SetSuffix("check cache of %s", pkg)
		exists, _, err := pluginManager.IsPluginInstalled(ctx, path, version)
		if err != nil {
			return err
		}
		if exists {
			progress.SetSuffix("%s is cached", pkg)
		} else {
			progress.SetSuffix("installing %s", pkg)
			_, err := pluginManager.InstallPlugin(ctx, path, version)
			if err != nil {
				panic(err)
			}
			progress.SetSuffix("%s installed", pkg)
		}
		pluginsMap[pkg] = struct{}{}
		progress.Incr()
	}
	progress.SetSuffix("all plugins have been installed")
	progress.Wait()

	fmt.Println("the following plugins will be used:")
	for pkg := range pluginsMap {
		fmt.Printf("	%s\r\n", pkg)
	}
	return nil
}

// StepCompile is used to compile proto files
func StepCompile(ctx context.Context,
	compilerManager compilermanager.CompilerManager,
	targets []string,
) error {
	progress := progressbar.GetProgressBar(ctx, len(targets))
	progress.SetPrefix("Compile Proto Files")
	c := concurrent.NewErrGroup(ctx, 10)
	for _, target := range targets {
		func(target string) {
			c.Go(func(ctx context.Context) error {
				progress.SetSuffix(target)
				comp, err := compilerManager.GetCompiler(ctx, target)
				if err != nil {
					return err
				}
				if err := comp.Compile(ctx, target); err != nil {
					return err
				}
				progress.Incr()
				return nil
			})
		}(target)
	}
	if err := c.Wait(); err != nil {
		return err
	}
	progress.Wait()
	return nil
}

// StepPostAction is used to execute post actions
func StepPostAction(ctx context.Context,
	actionsManager actionmanager.ActionManager,
	configItems []configs.ConfigItem) error {
	progress := progressbar.GetProgressBar(ctx, len(configItems))
	progress.SetPrefix("PostAction")
	for _, cfg := range configItems {
		progress.SetSuffix(cfg.Path())
		if err := actionsManager.ExecutePostAction(ctx, cfg); err != nil {
			return err
		}
		progress.Incr()
	}
	progress.Wait()
	return nil
}

// StepPostShell is used to execute post shell
func StepPostShell(ctx context.Context,
	actionsManager actionmanager.ActionManager,
	configItems []configs.ConfigItem) error {
	progress := progressbar.GetProgressBar(ctx, len(configItems))
	progress.SetPrefix("PostShell")
	for _, cfg := range configItems {
		progress.SetSuffix(cfg.Path())
		if err := actionsManager.ExecutePostShell(ctx, cfg); err != nil {
			return err
		}
		progress.Incr()
	}
	progress.Wait()
	return nil
}

// Compile is used to compile proto files
func Compile(ctx context.Context, targets []string) error {
	log := logger.NewDefault("compile")
	log.SetLogLevel(logger.LevelInfo)
	if consts.IsDebugMode(ctx) {
		log.SetLogLevel(logger.LevelDebug)
	}

	configManager, err := configmanager.NewConfigManager(log)
	if err != nil {
		return err
	}
	pluginManager, err := pluginmanager.NewPluginManager(pluginmanager.NewConfig(), log)
	if err != nil {
		return err
	}
	compilerManager, err := compilermanager.NewCompilerManager(ctx, log, configManager, pluginManager)
	if err != nil {
		return err
	}
	actionManager, err := actionmanager.NewActionManager(log)
	if err != nil {
		return err
	}

	configItems, err := StepLookUpConfigs(ctx, targets, configManager)
	if err != nil {
		return err
	}

	if err := StepInstallProtoc(ctx, pluginManager, configItems); err != nil {
		return err
	}
	if err := StepInstallRepositories(ctx, pluginManager, configItems); err != nil {
		return err
	}
	if err := StepInstallPlugins(ctx, pluginManager, configItems); err != nil {
		return err
	}
	if err := StepCompile(ctx, compilerManager, targets); err != nil {
		return err
	}

	if !consts.IsDisableAction(ctx) {
		if err := StepPostAction(ctx, actionManager, configItems); err != nil {
			return err
		}
		if err := StepPostShell(ctx, actionManager, configItems); err != nil {
			return err
		}
	} else {
		displayWarn := false
		for _, configItem := range configItems {
			if len(configItem.Config().PostActions) > 0 || configItem.Config().PostShell != "" {
				displayWarn = true
				break
			}
		}
		if displayWarn {
			log.LogWarn(nil, "PostAction and PostShell is skipped. If you need to allow execution, please append '-p' to command flags to enable")
		}
	}

	log.LogInfo(nil, "Good job! you are ready to go :)")
	return nil
}
