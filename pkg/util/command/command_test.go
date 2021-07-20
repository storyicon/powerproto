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

package command_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/storyicon/powerproto/pkg/util/command"
)

var _ = Describe("Command", func() {
	ctx := context.Background()
	It("should able to dryRun", func() {
		Expect(command.IsDryRun(ctx)).To(BeFalse())
		ctx = command.WithDryRun(ctx)
		Expect(command.IsDryRun(ctx)).To(BeTrue())
	})
	It("should able to ignore dryRun", func() {
		Expect(command.IsIgnoreDryRun(ctx)).To(BeFalse())
		ctx = command.WithIgnoreDryRun(ctx)
		Expect(command.IsIgnoreDryRun(ctx)).To(BeTrue())
	})
})
