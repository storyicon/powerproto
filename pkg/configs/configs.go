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

package configs

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ppaanngggg/powerproto/pkg/consts"
	"github.com/ppaanngggg/powerproto/pkg/util"
)

// Config defines the config model
type Config struct {
	Scopes        []string          `json:"scopes" yaml:"scopes"`
	Protoc        string            `json:"protoc" yaml:"protoc"`
	ProtocWorkDir string            `json:"protocWorkDir" yaml:"protocWorkDir"`
	Plugins       map[string]string `json:"plugins" yaml:"plugins"`
	Repositories  map[string]string `json:"repositories" yaml:"repositories"`
	Options       []string          `json:"options" yaml:"options"`
	ImportPaths   []string          `json:"importPaths" yaml:"importPaths"`
	PostActions   []*PostAction     `json:"postActions" yaml:"postActions"`
	PostShell     string            `json:"postShell" yaml:"postShell"`
}

// PostAction defines the Action model
type PostAction struct {
	Name string   `json:"name" yaml:"name"`
	Args []string `json:"args" yaml:"args"`
}

// SaveConfigs is used to save configs into files
func SaveConfigs(path string, configs ...*Config) error {
	parts := make([][]byte, 0, len(configs))
	for _, config := range configs {
		data, err := yaml.Marshal(config)
		if err != nil {
			return err
		}
		parts = append(parts, data)
	}
	data := bytes.Join(parts, []byte("\r\n---\r\n"))
	if err := ioutil.WriteFile(path, data, fs.ModePerm); err != nil {
		return err
	}
	return nil
}

// LoadConfigs is used to load config from specified path
func LoadConfigs(path string) ([]*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	items, err := util.SplitYAML(raw)
	if err != nil {
		return nil, err
	}
	var ret []*Config
	for _, item := range items {
		var config Config
		if err := yaml.Unmarshal(item, &config); err != nil {
			return nil, err
		}
		ret = append(ret, &config)
	}
	return ret, nil
}

// LoadConfigItems is similar to LoadConfigs, but obtains the abstraction of the Config prototype structure
func LoadConfigItems(path string) ([]ConfigItem, error) {
	data, err := LoadConfigs(path)
	if err != nil {
		return nil, err
	}
	return GetConfigItems(data, path), nil
}

// ListConfigPaths is used to list all possible config paths
func ListConfigPaths(sourceDir string) []string {
	var paths []string
	cur := sourceDir
	for {
		filePath := filepath.Join(cur, consts.ConfigFileName)
		paths = append(paths, filePath)
		next := filepath.Dir(cur)
		if next == cur {
			break
		}
		cur = next
	}
	paths = append(paths, consts.PathForGlobalConfig())
	return paths
}
