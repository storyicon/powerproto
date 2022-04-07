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

package actions

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ppaanngggg/powerproto/pkg/consts"
	"github.com/ppaanngggg/powerproto/pkg/util"
	"github.com/ppaanngggg/powerproto/pkg/util/logger"
)

// ActionMove is used to move directory or file from src to dest
// Its args prototype is:
// 		args: (src string, dest string)
func ActionMove(ctx context.Context, log logger.Logger, args []string, options *CommonOptions) error {
	if len(args) != 2 || util.ContainsEmpty(args...) {
		return errors.Errorf("expected length of args is 3, but received %d", len(args))
	}
	var (
		source         = args[0]
		destination    = args[1]
		path           = filepath.Dir(options.ConfigFilePath)
		absSource      = filepath.Join(path, source)
		absDestination = filepath.Join(path, destination)
	)

	if filepath.IsAbs(source) {
		return errors.Errorf("absolute source %s is not allowed in action move", source)
	}

	if filepath.IsAbs(destination) {
		return errors.Errorf("absolute destination %s is not allowed in action move", destination)
	}

	if consts.IsDryRun(ctx) {
		log.LogInfo(map[string]interface{}{
			"action": "move",
			"from":   absSource,
			"to":     absDestination,
		}, consts.TextDryRun)
		return nil
	}

	log.LogDebug(map[string]interface{}{
		"action": "move",
		"from":   absSource,
		"to":     absDestination,
	}, consts.TextExecuteAction)

	if err := util.CopyDirectory(absSource, absDestination); err != nil {
		return err
	}
	if err := os.RemoveAll(absSource); err != nil {
		return err
	}
	return nil
}
