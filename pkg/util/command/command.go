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

package command

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/storyicon/powerproto/pkg/consts"
	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// Execute is used to execute commands, return stdout and execute errors
func Execute(ctx context.Context,
	log logger.Logger,
	dir string, name string, arguments []string, env []string) ([]byte, *ErrCommandExec) {
	cmd := exec.CommandContext(ctx, name, arguments...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = dir

	if consts.IsDryRun(ctx) && !consts.IsIgnoreDryRun(ctx) {
		log.LogInfo(map[string]interface{}{
			"command": cmd.String(),
			"dir":     cmd.Dir,
		}, consts.TextDryRun)
		return nil, nil
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.LogDebug(map[string]interface{}{
		"command": cmd.String(),
		"dir":     cmd.Dir,
	}, consts.TextExecuteCommand)
	if err := cmd.Run(); err != nil {
		return nil, &ErrCommandExec{
			Err:      err,
			Dir:      cmd.Dir,
			Command:  cmd.String(),
			ExitCode: util.GetExitCode(err),
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
		}
	}
	return stdout.Bytes(), nil
}

// ErrCommandExec defines the command error
type ErrCommandExec struct {
	Err      error
	Dir      string
	Command  string
	ExitCode int
	Stdout   string
	Stderr   string
}

// Error implements the error interface
func (err *ErrCommandExec) Error() string {
	return fmt.Sprintf("failed to execute %s in %s, stderr: %s, exit code %d, %s",
		err.Command, err.Dir, err.Stderr, err.ExitCode, err.Err,
	)
}
