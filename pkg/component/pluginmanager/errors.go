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

package pluginmanager

import (
	"fmt"

	"github.com/ppaanngggg/powerproto/pkg/util/command"
)

// ErrGoInstall defines the go install error
type ErrGoInstall struct {
	*command.ErrCommandExec
}

// ErrGoList defines the go list error
type ErrGoList struct {
	*command.ErrCommandExec
}

// ErrHTTPDownload defines the download error
type ErrHTTPDownload struct {
	Url  string
	Err  error
	Code int
}

// Error implements the error interface
func (err *ErrHTTPDownload) Error() string {
	return fmt.Sprintf("failed to download %s, code: %d, err: %s", err.Url, err.Code, err.Err)
}
