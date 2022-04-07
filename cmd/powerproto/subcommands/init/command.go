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

package build

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/ppaanngggg/powerproto/pkg/configs"
	"github.com/ppaanngggg/powerproto/pkg/consts"
	"github.com/ppaanngggg/powerproto/pkg/util"
	"github.com/ppaanngggg/powerproto/pkg/util/logger"
)

// UserPreference defines the model of user preference
type UserPreference struct {
	Plugins      []string `survey:"plugins"`
	Repositories []string `survey:"repositories"`
}

// GetUserPreference is used to get user preference
func GetUserPreference() (*UserPreference, error) {
	var preference UserPreference
	err := survey.Ask([]*survey.Question{
		{
			Name: "plugins",
			Prompt: &survey.MultiSelect{
				Message: "select plugins to use. Later, you can also manually add in the configuration file",
				Options: GetWellKnownPluginsOptionValues(),
			},
		},
		{
			Name: "repositories",
			Prompt: &survey.MultiSelect{
				Message: "select repositories to use. Later, you can also manually add in the configuration file",
				Options: GetWellKnownRepositoriesOptionValues(),
			},
		},
	}, &preference)
	if len(preference.Plugins) == 0 {
		preference.Plugins = []string{
			GetPluginProtocGenGo().GetOptionsValue(),
			GetPluginProtocGenGoGRPC().GetOptionsValue(),
		}
	}
	if len(preference.Repositories) == 0 {
		preference.Repositories = []string{
			GetRepositoryGoogleAPIs().GetOptionsValue(),
		}
	}
	return &preference, err
}

// GetDefaultConfig is used to get default config
func GetDefaultConfig() *configs.Config {
	return &configs.Config{
		Scopes: []string{
			"./",
		},
		Protoc:       "latest",
		Plugins:      map[string]string{},
		Repositories: map[string]string{},
		Options:      []string{},
		ImportPaths: []string{
			".",
			"$GOPATH",
			consts.KeyPowerProtoInclude,
			consts.KeySourceRelative,
		},
	}
}

// CommandInit is used to initialize the configuration
// file in the current directory
func CommandInit(log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "init a config file in current directory",
		Run: func(cmd *cobra.Command, args []string) {
			exists, err := util.IsFileExists(consts.ConfigFileName)
			if err != nil {
				log.LogError(nil, err.Error())
				return
			}
			if exists {
				log.LogInfo(nil, "config file already exists in this directory")
				return
			}

			preference, err := GetUserPreference()
			if err != nil {
				log.LogError(nil, err.Error())
				return
			}
			config := GetDefaultConfig()
			for _, val := range preference.Plugins {
				if plugin, ok := GetPluginFromOptionsValue(val); ok {
					config.Plugins[plugin.Name] = plugin.Pkg
					config.Options = append(config.Options, plugin.Options...)
				}
			}
			for _, val := range preference.Repositories {
				if repo, ok := GetRepositoryFromOptionsValue(val); ok {
					config.Repositories[repo.Name] = repo.Pkg
					config.ImportPaths = append(config.ImportPaths, repo.ImportPaths...)
				}
			}
			if err := configs.SaveConfigs(consts.ConfigFileName, config); err != nil {
				log.LogFatal(nil, "failed to save config: %s", err)
			}
			log.LogInfo(nil, "succeed! You can use `powerproto tidy` to tidy this config file")
		},
	}
}
