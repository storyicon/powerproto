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

package actionmanager

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ppaanngggg/powerproto/pkg/component/actionmanager/actions"
	"github.com/ppaanngggg/powerproto/pkg/configs"
	"github.com/ppaanngggg/powerproto/pkg/util/command"
	"github.com/ppaanngggg/powerproto/pkg/util/logger"
)

// ActionManager is used to manage actions
type ActionManager interface {
	// ExecutePostShell is used to execute post shell in config item
	ExecutePostShell(ctx context.Context, config configs.ConfigItem) error
	// ExecutePostAction is used to execute post action in config item
	ExecutePostAction(ctx context.Context, config configs.ConfigItem) error
}

// BasicActionManager is a basic implement of ActionManager
type BasicActionManager struct {
	logger.Logger

	// map[string]ActionFunc
	actions map[string]actions.ActionFunc
}

// NewActionManager is used to create action manager
func NewActionManager(log logger.Logger) (ActionManager, error) {
	return NewBasicActionManager(log)
}

// NewBasicActionManager is used to create a BasicActionManager
func NewBasicActionManager(log logger.Logger) (*BasicActionManager, error) {
	return &BasicActionManager{
		Logger: log.NewLogger("actionmanager"),
		actions: map[string]actions.ActionFunc{
			"move":    actions.ActionMove,
			"replace": actions.ActionReplace,
			"remove":  actions.ActionRemove,
			"copy":    actions.ActionCopy,
		},
	}, nil
}

// ExecutePostShell is used to execute post shell in config item
func (m *BasicActionManager) ExecutePostShell(ctx context.Context, config configs.ConfigItem) error {
	script := config.Config().PostShell
	if script == "" {
		return nil
	}
	dir := filepath.Dir(config.Path())
	_, err := command.Execute(ctx, m.Logger, dir, "/bin/sh", []string{
		"-c", script,
	}, nil)
	if err != nil {
		return &ErrPostShell{
			Path:           config.Path(),
			ErrCommandExec: err.(*command.ErrCommandExec),
		}
	}
	return nil
}

// ExecutePostAction is used to execute post action in config item
func (m *BasicActionManager) ExecutePostAction(ctx context.Context, config configs.ConfigItem) error {
	for _, action := range config.Config().PostActions {
		actionFunc, ok := m.actions[action.Name]
		if !ok {
			return fmt.Errorf("unknown action: %s", action.Name)
		}
		if err := actionFunc(ctx, m.Logger, action.Args, &actions.CommonOptions{
			ConfigFilePath: config.Path(),
		}); err != nil {
			return &ErrPostAction{
				Path:      config.Path(),
				Name:      action.Name,
				Arguments: action.Args,
				Err:       err,
			}
		}
	}
	return nil
}
