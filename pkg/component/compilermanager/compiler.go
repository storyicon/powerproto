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
	// GetConfig is used to return config that the compiler used
	GetConfig(ctx context.Context) configs.ConfigItem
}

var _ Compiler = &BasicCompiler{}

// BasicCompiler is the basic implement of Compiler
type BasicCompiler struct {
	logger.Logger
	config        configs.ConfigItem
	pluginManager pluginmanager.PluginManager

	protocPath     string
	arguments      []string
	googleapisPath string
	dir            string
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
	basic := &BasicCompiler{
		Logger:        log.NewLogger("compiler"),
		config:        config,
		pluginManager: pluginManager,
	}
	if err := basic.calcGoogleAPIs(ctx); err != nil {
		return nil, err
	}
	if err := basic.calcProto(ctx); err != nil {
		return nil, err
	}
	if err := basic.calcArguments(ctx); err != nil {
		return nil, err
	}
	if err := basic.calcDir(ctx); err != nil {
		return nil, err
	}
	return basic, nil
}

// Compile is used to compile proto file
func (b *BasicCompiler) Compile(ctx context.Context, protoFilePath string) error {
	arguments := make([]string, len(b.arguments))
	copy(arguments, b.arguments)

	// contains source relative
	if util.Contains(b.config.Config().ImportPaths,
		consts.KeySourceRelative) {
		arguments = append(arguments, "--proto_path="+filepath.Dir(protoFilePath))
	}

	arguments = append(arguments, protoFilePath)
	_, err := command.Execute(ctx,
		b.Logger, b.dir, b.protocPath, arguments, nil)
	if err != nil {
		return &ErrCompile{
			ErrCommandExec: err,
		}
	}
	return nil
}

// GetConfig is used to return config that the compiler used
func (b *BasicCompiler) GetConfig(ctx context.Context) configs.ConfigItem {
	return b.config
}

func (b *BasicCompiler) calcDir(ctx context.Context) error {
	if dir := b.config.Config().ProtocWorkDir; dir != "" {
		dir = util.RenderPathWithEnv(dir)
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(b.config.Path(), dir)
		}
		b.dir = dir
	} else {
		b.dir = filepath.Dir(b.config.Path())
	}
	return nil
}

func (b *BasicCompiler) calcGoogleAPIs(ctx context.Context) error {
	cfg := b.config.Config()
	commitId := cfg.GoogleAPIs
	if commitId == "" {
		return nil
	}
	if !util.Contains(cfg.ImportPaths, consts.KeyPowerProtoGoogleAPIs) {
		return nil
	}
	if commitId == "latest" {
		latestVersion, err := b.pluginManager.GetGoogleAPIsLatestVersion(ctx)
		if err != nil {
			return err
		}
		commitId = latestVersion
	}
	local, err := b.pluginManager.InstallGoogleAPIs(ctx, commitId)
	if err != nil {
		return err
	}
	b.googleapisPath = local
	return nil
}

func (b *BasicCompiler) calcProto(ctx context.Context) error {
	cfg := b.config
	protocVersion := cfg.Config().Protoc
	if protocVersion == "latest" {
		latestVersion, err := b.pluginManager.GetProtocLatestVersion(ctx)
		if err != nil {
			return err
		}
		protocVersion = latestVersion
	}
	localPath, err := b.pluginManager.InstallProtoc(ctx, protocVersion)
	if err != nil {
		return err
	}
	b.protocPath = localPath
	return nil
}

func (b *BasicCompiler) calcArguments(ctx context.Context) error {
	cfg := b.config
	arguments := make([]string, len(cfg.Config().Options))
	copy(arguments, cfg.Config().Options)

	dir := filepath.Dir(cfg.Path())
	// build import paths
Loop:
	for _, path := range cfg.Config().ImportPaths {
		switch path {
		case consts.KeyPowerProtoInclude:
			path = b.pluginManager.IncludePath(ctx)
		case consts.KeyPowerProtoGoogleAPIs:
			if b.googleapisPath != "" {
				path = b.googleapisPath
			}
		case consts.KeySourceRelative:
			continue Loop
		}
		path = util.RenderPathWithEnv(path)
		if !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
		}
		arguments = append(arguments, "--proto_path="+path)
	}

	// build plugin options
	for name, pkg := range cfg.Config().Plugins {
		path, version, ok := util.SplitGoPackageVersion(pkg)
		if !ok {
			return errors.Errorf("failed to parse: %s", pkg)
		}
		if version == "latest" {
			latestVersion, err := b.pluginManager.GetPluginLatestVersion(ctx, path)
			if err != nil {
				return err
			}
			version = latestVersion
		}
		local, err := b.pluginManager.InstallPlugin(ctx, path, version)
		if err != nil {
			return errors.Wrap(err, "failed to get plugin path")
		}
		arg := fmt.Sprintf("--plugin=%s=%s", name, local)
		arguments = append(arguments, arg)
	}
	b.arguments = arguments
	return nil
}
