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
	"fmt"
)

// Plugin defines the plugin options
type Plugin struct {
	OptionsValue string

	Name    string
	Pkg     string
	Options []string
}

// GetOptionsValue is used to get options value of plugin
func (plugin *Plugin) GetOptionsValue() string {
	if plugin.OptionsValue != "" {
		return plugin.OptionsValue
	}
	return fmt.Sprintf("%s: %s", plugin.Name, plugin.Pkg)
}

// GetPluginProtocGenGo is used to get protoc-gen-go plugin
func GetPluginProtocGenGo() *Plugin {
	return &Plugin{
		Name: "protoc-gen-go",
		Pkg:  "google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		Options: []string{
			"--go_out=.",
			"--go_opt=paths=source_relative",
		},
	}
}

// GetPluginProtocGenGoGRPC is used to get protoc-gen-go-grpc plugin
func GetPluginProtocGenGoGRPC() *Plugin {
	return &Plugin{
		Name: "protoc-gen-go-grpc",
		Pkg:  "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		Options: []string{
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
		},
	}
}

// GetWellKnownPlugins is used to get well known plugins
func GetWellKnownPlugins() []*Plugin {
	return []*Plugin{
		GetPluginProtocGenGo(),
		GetPluginProtocGenGoGRPC(),
		{
			Name: "protoc-gen-grpc-gateway",
			Pkg:  "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest",
			Options: []string{
				"--grpc-gateway_out=.",
				"--grpc-gateway_opt=paths=source_relative",
			},
		},
		{
			OptionsValue: "protoc-gen-go-grpc (SupportPackageIsVersion6)",
			Name:         "protoc-gen-go-grpc",
			Pkg:          "google.golang.org/grpc/cmd/protoc-gen-go-grpc@ad51f572fd270f2323e3aa2c1d2775cab9087af2",
			Options: []string{
				"--go-grpc_out=.",
				"--go-grpc_opt=paths=source_relative",
			},
		},
		{
			Name: "protoc-gen-openapiv2",
			Pkg:  "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
			Options: []string{
				"--openapiv2_out=.",
			},
		},
		{
			Name: "protoc-gen-gogo",
			Pkg:  "github.com/gogo/protobuf/protoc-gen-gogo@latest",
			Options: []string{
				"--gogo_out=plugins=grpc:.",
				"--gogo_opt=paths=source_relative",
			},
		},
		{
			Name: "protoc-gen-gofast",
			Pkg:  "github.com/gogo/protobuf/protoc-gen-gofast@latest",
			Options: []string{
				"--gofast_out=plugins=grpc:.",
				"--gofast_opt=paths=source_relative",
			},
		},
		{
			Name: "protoc-gen-golang-deepcopy",
			Pkg:  "istio.io/tools/cmd/protoc-gen-golang-deepcopy@latest",
			Options: []string{
				"--golang-deepcopy_out=source_relative:.",
			},
		},
		{
			Name: "protoc-gen-go-json",
			Pkg:  "github.com/mitchellh/protoc-gen-go-json@latest",
			Options: []string{
				"--go-json_out=.",
			},
		},
	}
}

// GetPluginFromOptionsValue is used to get plugin by option value
func GetPluginFromOptionsValue(val string) (*Plugin, bool) {
	plugins := GetWellKnownPlugins()
	for _, plugin := range plugins {
		if plugin.GetOptionsValue() == val {
			return plugin, true
		}
	}
	return nil, false
}

// GetWellKnownPluginsOptionValues is used to get option values of well known plugins
func GetWellKnownPluginsOptionValues() []string {
	plugins := GetWellKnownPlugins()
	packages := make([]string, 0, len(plugins))
	for _, plugin := range plugins {
		packages = append(packages, plugin.GetOptionsValue())
	}
	return packages
}
