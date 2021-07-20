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

	"github.com/storyicon/powerproto/pkg/util/logger"
)

// CommonOptions is common options of action
type CommonOptions struct {
	ConfigFilePath string
}

// ActionFunc defines the common prototype of action func
type ActionFunc func(ctx context.Context,
	log logger.Logger,
	args []string, options *CommonOptions) error
