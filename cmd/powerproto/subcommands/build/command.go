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

package build

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/storyicon/powerproto/pkg/bootstraps"
	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

const description = `
Examples:
compile specific proto file
	powerproto build [proto file] 

compile the proto file in the folder, excluding sub folders:
	powerproto build [dir] 

compile all proto files in the folder recursively, including sub folders:
	powerproto build -r [dir] 

compile proto files and execute the post actions/shells:
	powerproto build -r -a [dir]
`

// CommandBuild is used to compile proto files
// powerproto build -r .
// powerproto build .
// powerproto build xxxxx.proto
func CommandBuild(log logger.Logger) *cobra.Command {
	var recursive bool
	var dryRun bool
	var debugMode bool
	var postScriptEnabled bool
	cmd := &cobra.Command{
		Use:   "build [dir|proto file]",
		Short: "compile proto files",
		Long:  strings.TrimSpace(description),
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.SetLogLevel(logger.LevelInfo)
			ctx := cmd.Context()
			if debugMode {
				ctx = consts.WithDebugMode(ctx)
				log.LogWarn(nil, "running in debug mode")
				log.SetLogLevel(logger.LevelDebug)
			}
			if dryRun {
				ctx = consts.WithDryRun(ctx)
				log.LogWarn(nil, "running in dryRun mode")
			}
			if !postScriptEnabled {
				ctx = consts.WithDisableAction(ctx)
			}

			target, err := filepath.Abs(args[0])
			if err != nil {
				log.LogFatal(nil, "failed to abs target path: %s", err)
			}
			fileInfo, err := os.Stat(target)
			if err != nil {
				log.LogFatal(map[string]interface{}{
					"target": target,
				}, "failed to stat target: %s", err)
			}

			var targets []string
			if fileInfo.IsDir() {
				log.LogInfo(nil, "search proto files...")
				if recursive {
					targets, err = util.GetFilesWithExtRecursively(target, ".proto")
					if err != nil {
						log.LogFatal(nil, "failed to walk directory: %s", err)
					}
				} else {
					targets, err = util.GetFilesWithExt(target, ".proto")
					if err != nil {
						log.LogFatal(nil, "failed to walk directory: %s", err)
					}
				}
			} else {
				targets = append(targets, target)
			}

			if len(targets) == 0 {
				log.LogWarn(nil, "no file to compile")
				return
			}
			if err := bootstraps.StepTidyConfig(ctx, targets); err != nil {
				log.LogFatal(nil, "failed to tidy config: %+v", err)
				return
			}

			if err := bootstraps.Compile(ctx, targets); err != nil {
				log.LogFatal(nil, "failed to compile: %+v", err)
			}

			log.LogInfo(nil, "succeed! you are ready to go :)")
		},
	}
	flags := cmd.PersistentFlags()
	flags.BoolVarP(&recursive, "recursive", "r", recursive, "whether to recursively traverse all child folders")
	flags.BoolVarP(&postScriptEnabled, "postScriptEnabled", "p", postScriptEnabled, "when this flag is attached, it will allow the execution of postActions and postShell")
	flags.BoolVarP(&debugMode, "debug", "d", debugMode, "debug mode")
	flags.BoolVarP(&dryRun, "dryRun", "y", dryRun, "dryRun mode")
	return cmd
}
