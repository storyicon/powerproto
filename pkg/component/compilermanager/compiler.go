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
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/storyicon/powerproto/pkg/component/pluginmanager"
	"github.com/storyicon/powerproto/pkg/configs"
	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/command"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// Compiler is used to compile proto file
type Compiler interface {
	// Compile is used to compile proto file
	Compile(ctx context.Context, protoFilePath string) error
	// BatchCompile is used to compile a batch of proto files
	BatchCompile(ctx context.Context, protoFilePaths []string) error
	// GetConfig is used to return config that the compiler used
	GetConfig(ctx context.Context) configs.ConfigItem
}

var _ Compiler = &BasicCompiler{}

// BasicCompiler is the basic implement of Compiler
type BasicCompiler struct {
	logger.Logger
	config        configs.ConfigItem
	pluginManager pluginmanager.PluginManager
}

// NewCompiler is used to create a compiler
func NewCompiler(
	ctx context.Context,
	log logger.Logger,
	pluginManager pluginmanager.PluginManager,
	config configs.ConfigItem,
) (Compiler, error) {
	return NewBasicCompiler(ctx, log, pluginManager, config)
}

// NewBasicCompiler is used to create a basic compiler
func NewBasicCompiler(
	ctx context.Context,
	log logger.Logger,
	pluginManager pluginmanager.PluginManager,
	config configs.ConfigItem,
) (*BasicCompiler, error) {
	return &BasicCompiler{
		Logger:        log.NewLogger("compiler"),
		config:        config,
		pluginManager: pluginManager,
	}, nil
}

// Compile is used to compile proto file
func (b *BasicCompiler) Compile(ctx context.Context, protoFilePath string) error {
	dir := b.calcDir()
	protocPath, err := b.calcProtocPath(ctx)
	if err != nil {
		return err
	}

	arguments, err := b.calcArguments(ctx, protoFilePath)
	if err != nil {
		return err
	}
	arguments = append(arguments, protoFilePath)
	_, err = command.Execute(ctx,
		b.Logger, dir, protocPath, arguments, nil)
	if err != nil {
		return &ErrCompile{
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	return nil
}

// GetConfig is used to return config that the compiler used
func (b *BasicCompiler) GetConfig(ctx context.Context) configs.ConfigItem {
	return b.config
}

func (b *BasicCompiler) calcProtocPath(ctx context.Context) (string, error) {
	return b.pluginManager.GetPathForProtoc(ctx, b.config.Config().Protoc)
}

func (b *BasicCompiler) calcDir() string {
	if dir := b.config.Config().ProtocWorkDir; dir != "" {
		dir = util.RenderPathWithEnv(dir, nil)
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(b.config.Path(), dir)
		}
		return dir
	} else {
		return filepath.Dir(b.config.Path())
	}
}

func (b *BasicCompiler) calcArguments(ctx context.Context, protoFilePath string) ([]string, error) {
	variables, err := b.calcVariables(ctx, protoFilePath)
	if err != nil {
		return nil, err
	}

	cfg := b.config
	var arguments []string

	// build plugin options
	for name, pkg := range cfg.Config().Plugins {
		path, version, ok := util.SplitGoPackageVersion(pkg)
		if !ok {
			return nil, errors.Errorf("failed to parse: %s", pkg)
		}
		local, err := b.pluginManager.GetPathForPlugin(ctx, path, version)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get plugin path")
		}
		arg := fmt.Sprintf("--plugin=%s=%s", name, local)
		arguments = append(arguments, arg)
	}

	// build compile options
	for _, option := range cfg.Config().Options {
		option = util.RenderWithEnv(option, variables)
		arguments = append(arguments, option)
	}

	// build import paths
	dir := filepath.Dir(cfg.Path())
	for _, path := range cfg.Config().ImportPaths {
		path = util.RenderPathWithEnv(path, variables)
		if !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
		}
		arguments = append(arguments, "--proto_path="+path)
	}

	arguments = util.DeduplicateSliceStably(arguments)

	return arguments, nil
}

func (b *BasicCompiler) calcVariables(ctx context.Context, protoFilePath string) (map[string]string, error) {
	cfg := b.config
	variables := map[string]string{}
	for name, pkg := range cfg.Config().Repositories {
		_, version, ok := util.SplitGoPackageVersion(pkg)
		if !ok {
			return nil, errors.Errorf("failed to parse: %s", pkg)
		}
		repoPath, err := b.pluginManager.GitRepoPath(ctx, version)
		if err != nil {
			return nil, err
		}
		variables[name] = repoPath
	}
	includePath, err := b.pluginManager.IncludePath(ctx)
	if err != nil {
		return nil, err
	}
	variables[consts.KeyNamePowerProtocInclude] = includePath
	variables[consts.KeyNameSourceRelative] = filepath.Dir(protoFilePath)
	return variables, nil
}

// BatchCompile is used to compile a batch of proto files
func (b *BasicCompiler) BatchCompile(ctx context.Context, protoFilePaths []string) error {
	dir := b.calcDir()
	protocPath, err := b.calcProtocPath(ctx)
	if err != nil {
		return err
	}

	arguments, err := b.calcArguments(ctx, util.GetCommonRootDirOfPaths(protoFilePaths))
	if err != nil {
		return err
	}
	arguments = append(arguments, protoFilePaths...)
	_, err = command.Execute(ctx,
		b.Logger, dir, protocPath, arguments, nil)
	if err != nil {
		return &ErrCompile{
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	return nil
}
