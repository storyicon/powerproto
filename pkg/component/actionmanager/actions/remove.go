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

	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// ActionRemove is used to delete directories or files
// Its args prototype is:
// 		args: (path ...string)
func ActionRemove(ctx context.Context, log logger.Logger, args []string, options *CommonOptions) error {
	for _, arg := range args {
		if filepath.IsAbs(arg) {
			return errors.Errorf("absolute path %s is not allowed in action remove", arg)
		}
		path := filepath.Join(filepath.Dir(options.ConfigFilePath), arg)

		if consts.IsDryRun(ctx) {
			log.LogInfo(map[string]interface{}{
				"action": "remove",
				"target": path,
			}, consts.TextDryRun)
			return nil
		}

		log.LogDebug(map[string]interface{}{
			"action": "remove",
			"target": path,
		}, consts.TextExecuteAction)

		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}
