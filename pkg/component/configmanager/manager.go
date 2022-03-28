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

package configmanager

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/storyicon/powerproto/pkg/configs"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// ConfigManager is used to manage config
type ConfigManager interface {
	// GetConfig is used to get config of specified proto file path
	GetConfig(ctx context.Context, protoFilePath string) (configs.ConfigItem, error)
}

// NewConfigManager is used to create ConfigManager
func NewConfigManager(log logger.Logger) (ConfigManager, error) {
	return NewBasicConfigManager(log)
}

// BasicConfigManager is the basic implement of ConfigManager
type BasicConfigManager struct {
	logger.Logger

	tree     map[string][]configs.ConfigItem
	treeLock sync.RWMutex
}

// NewBasicConfigManager is used to create a basic ConfigManager
func NewBasicConfigManager(log logger.Logger) (*BasicConfigManager, error) {
	return &BasicConfigManager{
		Logger: log.NewLogger("configmanager"),
		tree:   map[string][]configs.ConfigItem{},
	}, nil
}

// GetConfig is used to get config of specified proto file path
func (b *BasicConfigManager) GetConfig(ctx context.Context, protoFilePath string) (configs.ConfigItem, error) {
	possiblePath := configs.ListConfigPaths(filepath.Dir(protoFilePath))
	for _, configFilePath := range possiblePath {
		items, err := b.loadConfig(configFilePath)
		if err != nil {
			return nil, err
		}
		dir := filepath.Dir(configFilePath)
		for _, config := range items {
			for _, scope := range config.Config().Scopes {
				scopePath := filepath.Join(dir, scope)
				if strings.Contains(protoFilePath, scopePath) {
					return config, nil
				}
			}
		}
	}
	return nil, errors.Errorf("unable to find config: %s", protoFilePath)
}

func (b *BasicConfigManager) loadConfig(configFilePath string) ([]configs.ConfigItem, error) {
	b.treeLock.Lock()
	defer b.treeLock.Unlock()
	configItems, ok := b.tree[configFilePath]
	if !ok {
		fileInfo, err := os.Stat(configFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				b.tree[configFilePath] = nil
				return nil, nil
			}
			return nil, err
		}
		if fileInfo.IsDir() {
			b.tree[configFilePath] = nil
			return nil, nil
		}
		data, err := configs.LoadConfigItems(configFilePath)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to decode: %s", configFilePath)
		}
		b.tree[configFilePath] = data
		return data, nil
	}
	return configItems, nil
}
