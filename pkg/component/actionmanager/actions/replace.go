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
	"bytes"
	"context"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ppaanngggg/powerproto/pkg/consts"
	"github.com/ppaanngggg/powerproto/pkg/util"
	"github.com/ppaanngggg/powerproto/pkg/util/logger"
)

// ActionReplace is used to replace text in bulk
// Its args prototype is:
// 		args: (pattern string, from string, to string)
// pattern is used to match files
func ActionReplace(ctx context.Context, log logger.Logger, args []string, options *CommonOptions) error {
	if len(args) != 3 {
		return errors.Errorf("expected length of args is 3, but received %d", len(args))
	}
	var (
		pattern = args[0]
		path    = filepath.Dir(options.ConfigFilePath)
		from    = args[1]
		to      = args[2]
	)

	if pattern == "" || from == "" {
		return errors.Errorf("pattern and from arguments in action replace can not be empty")
	}
	if filepath.IsAbs(pattern) {
		return errors.Errorf("absolute path %s is not allowed in action replace", pattern)
	}
	pattern = filepath.Join(path, pattern)
	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		matched, err := util.MatchPath(pattern, path)
		if err != nil {
			panic(err)
		}
		if matched {
			if consts.IsDryRun(ctx) {
				log.LogInfo(map[string]interface{}{
					"action": "replace",
					"file":   path,
					"from":   from,
					"to":     to,
				}, consts.TextDryRun)
				return nil
			}
			log.LogDebug(map[string]interface{}{
				"action": "replace",
				"file":   path,
				"from":   from,
				"to":     to,
			}, consts.TextExecuteAction)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			data = bytes.ReplaceAll(data, []byte(from), []byte(to))
			return ioutil.WriteFile(path, data, fs.ModePerm)
		}
		return nil
	})
}
