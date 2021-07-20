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

	"github.com/storyicon/powerproto/pkg/util"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

type dryRun struct{}
type ignoreDryRun struct{}
type disableAction struct{}

// WithIgnoreDryRun is used to inject ignore dryRun flag into context
func WithIgnoreDryRun(ctx context.Context) context.Context {
	return context.WithValue(ctx, ignoreDryRun{}, "true")
}

// WithDisableAction is used to disable post action/shell
func WithDisableAction(ctx context.Context) context.Context {
	return context.WithValue(ctx, disableAction{}, "true")
}

// WithDryRun is used to inject dryRun flag into context
func WithDryRun(ctx context.Context) context.Context {
	return context.WithValue(ctx, dryRun{}, "true")
}

// IsDisableAction is used to decide whether to disable PostAction and PostShell
func IsDisableAction(ctx context.Context) bool {
	return ctx.Value(disableAction{}) != nil
}

// IsIgnoreDryRun is used to determine whether it is currently in ignore DryRun mode
func IsIgnoreDryRun(ctx context.Context) bool {
	return ctx.Value(ignoreDryRun{}) != nil
}

// IsDryRun is used to determine whether it is currently in DryRun mode
func IsDryRun(ctx context.Context) bool {
	return ctx.Value(dryRun{}) != nil
}

// Execute is used to execute commands, return stdout and execute errors
func Execute(ctx context.Context,
	log logger.Logger,
	dir string, name string, arguments []string, env []string) ([]byte, *ErrCommandExec) {
	cmd := exec.CommandContext(ctx, name, arguments...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = dir

	if IsDryRun(ctx) && !IsIgnoreDryRun(ctx) {
		log.LogInfo(map[string]interface{}{
			"command": cmd.String(),
			"dir":     cmd.Dir,
		}, "DryRun")
		return nil, nil
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
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
