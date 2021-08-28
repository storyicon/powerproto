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

package consts

import (
	"context"
	"time"
)

type debugMode struct{}
type dryRun struct{}
type ignoreDryRun struct{}
type disableAction struct{}
type perCommandTimeout struct{}

func GetContextWithPerCommandTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	val := ctx.Value(perCommandTimeout{})
	if val == nil {
		return ctx, func() {}
	}
	duration, ok := val.(time.Duration)
	if !ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, duration)
}

// WithPerCommandTimeout is used to inject per command timeout
func WithPerCommandTimeout(ctx context.Context, timeout time.Duration) context.Context {
	return context.WithValue(ctx, perCommandTimeout{}, timeout)
}

// WithDebugMode is used to set debug mode
func WithDebugMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, debugMode{}, "true")
}

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
	ctx = WithDebugMode(ctx)
	return context.WithValue(ctx, dryRun{}, "true")
}

// IsDebugMode is used to decide whether to disable PostAction and PostShell
func IsDebugMode(ctx context.Context) bool {
	return ctx.Value(debugMode{}) != nil
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
