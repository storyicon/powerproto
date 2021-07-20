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
	Name    string
	Pkg     string
	Options []string
}

// String implements the standard string interface
func (options *Plugin) String() string {
	return fmt.Sprintf("%s: %s", options.Name, options.Pkg)
}

// GetWellKnownPlugins is used to get well known plugins
func GetWellKnownPlugins() []*Plugin {
	return []*Plugin{
		{
			Name: "protoc-gen-go",
			Pkg:  "google.golang.org/protobuf/cmd/protoc-gen-go@latest",
			Options: []string{
				"--go_out=.",
				"--go_opt=paths=source_relative",
			},
		},
		{
			Name: "protoc-gen-go-grpc",
			Pkg:  "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
			Options: []string{
				"--go-grpc_out=.",
				"--go-grpc_opt=paths=source_relative",
			},
		},
		{
			Name: "protoc-gen-grpc-gateway",
			Pkg:  "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest",
			Options: []string{
				"--grpc-gateway_out=.",
			},
		},
		{
			Name: "protoc-gen-deepcopy",
			Pkg:  "istio.io/tools/cmd/protoc-gen-deepcopy@latest",
			Options: []string{
				"--deepcopy_out=source_relative:.",
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
		if plugin.String() == val {
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
		packages = append(packages, plugin.String())
	}
	return packages
}
