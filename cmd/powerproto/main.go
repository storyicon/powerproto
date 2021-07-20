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

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdbuild "github.com/storyicon/powerproto/cmd/powerproto/subcommands/build"
	cmdenv "github.com/storyicon/powerproto/cmd/powerproto/subcommands/env"
	cmdinit "github.com/storyicon/powerproto/cmd/powerproto/subcommands/init"
	cmdtidy "github.com/storyicon/powerproto/cmd/powerproto/subcommands/tidy"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

// Version is set via build flag -ldflags -X main.Version
var (
	Version   string
	Branch    string
	Revision  string
	BuildDate string
)

var log = logger.NewDefault("command")

// GetRootCommand is used to get root command
func GetRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "[powerproto]",
		Version: fmt.Sprintf("%s, branch: %s, revision: %s, buildDate: %s", Version, Branch, Revision, BuildDate),
		Short:   "powerproto is used to build proto files and version control of protoc and related plug-ins",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
}

func main() {
	cmdRoot := GetRootCommand()
	cmdRoot.AddCommand(
		cmdbuild.CommandBuild(log),
		cmdinit.CommandInit(log),
		cmdtidy.CommandTidy(log),
		cmdenv.CommandEnv(log),
	)
	cmdRoot.Execute()
}
