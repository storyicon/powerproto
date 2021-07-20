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
	"fmt"
)

// ConfigItem is the wrapper of Config
type ConfigItem interface {
	// ID is used to return to config unique id
	ID() string
	// Path is used to return the config path
	Path() string
	// Config is used to return the Config
	Config() *Config
}

// GetConfigs is used to generate Config from given config entity
func GetConfigItems(data []*Config, path string) []ConfigItem {
	ret := make([]ConfigItem, 0, len(data))
	for i, item := range data {
		ret = append(ret, newConfigItem(item, path, i))
	}
	return ret
}

type configItem struct {
	id   string
	c    *Config
	path string
}

// ID is used to return to config unique id
func (c *configItem) ID() string {
	return c.id
}

// Path is used to return the config path
func (c *configItem) Path() string {
	return c.path
}

// Config is used to return the Config
func (c *configItem) Config() *Config {
	return c.c
}

func newConfigItem(c *Config, path string, idx int) ConfigItem {
	return &configItem{
		id:   fmt.Sprintf("%s:%d", path, idx),
		c:    c,
		path: path,
	}
}
