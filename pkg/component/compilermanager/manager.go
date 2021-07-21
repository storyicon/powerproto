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

package compilermanager

import (
	"context"
	"sync"

	"github.com/storyicon/powerproto/pkg/component/configmanager"
	"github.com/storyicon/powerproto/pkg/component/pluginmanager"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// CompilerManager is to manage compiler
type CompilerManager interface {
	// GetCompiler is used to get compiler of specified proto file path
	GetCompiler(ctx context.Context, protoFilePath string) (Compiler, error)
}

// BasicCompilerManager is the basic implement of CompilerManager
type BasicCompilerManager struct {
	logger.Logger

	configManager configmanager.ConfigManager
	pluginManager pluginmanager.PluginManager

	tree     map[string]Compiler
	treeLock sync.RWMutex
}

var _ CompilerManager = &BasicCompilerManager{}

// NewCompilerManager is used to create CompilerManager
func NewCompilerManager(ctx context.Context,
	log logger.Logger,
	configManager configmanager.ConfigManager,
	pluginManager pluginmanager.PluginManager,
) (CompilerManager, error) {
	return NewBasicCompilerManager(ctx, log, configManager, pluginManager)
}

// NewBasicCompilerManager is used to create basic CompilerManager
func NewBasicCompilerManager(ctx context.Context,
	log logger.Logger,
	configManager configmanager.ConfigManager,
	pluginManager pluginmanager.PluginManager,
) (*BasicCompilerManager, error) {
	return &BasicCompilerManager{
		Logger:        log.NewLogger("compilermanager"),
		configManager: configManager,
		pluginManager: pluginManager,
		tree:          map[string]Compiler{},
	}, nil
}

// GetCompiler is used to get compiler of specified proto file path
func (b *BasicCompilerManager) GetCompiler(ctx context.Context, protoFilePath string) (Compiler, error) {
	config, err := b.configManager.GetConfig(ctx, protoFilePath)
	if err != nil {
		return nil, err
	}

	b.treeLock.Lock()
	defer b.treeLock.Unlock()
	c, ok := b.tree[config.Path()]
	if ok {
		return c, nil
	}

	compiler, err := NewCompiler(ctx, b.Logger, b.pluginManager, config)
	if err != nil {
		return nil, err
	}
	b.tree[config.ID()] = compiler
	return compiler, nil
}
