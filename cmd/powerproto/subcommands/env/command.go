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

package env

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// CommandEnv is used to print environment variables related to program operation
func CommandEnv(log logger.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "list the environments and binary files",
		Run: func(cmd *cobra.Command, args []string) {
			log.LogInfo(nil, "[ENVIRONMENT]")
			for _, key := range []string{
				consts.EnvHomeDir,
				"HTTP_PROXY",
				"HTTPS_PROXY",
				"GOPROXY",
			} {
				log.LogInfo(nil, "%s=%s", key, os.Getenv(key))
			}

			log.LogInfo(nil, "[BIN]")
			for _, key := range []string{
				"go",
				"git",
			} {
				path, err := exec.LookPath(key)
				if err != nil {
					log.LogError(nil, "failed to find: %s", key)
					continue
				}
				log.LogInfo(nil, "%s=%s", key, path)
			}
		},
	}
}
