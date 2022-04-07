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
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"

	"github.com/ppaanngggg/powerproto/pkg/util"
)

// ProtocRelease defines the release of protoc
type ProtocRelease struct {
	workspace string
}

// GetIncludePath is used to get the include path
func (p *ProtocRelease) GetIncludePath() string {
	return filepath.Join(p.workspace, "include")
}

// GetProtocPath is used to get the protoc path
func (p *ProtocRelease) GetProtocPath() string {
	return filepath.Join(p.workspace, "bin", util.GetBinaryFileName("protoc"))
}

// Clear is used to clear the workspace
func (p *ProtocRelease) Clear() error {
	return os.RemoveAll(p.workspace)
}

// GetProtocRelease is used to download protoc release
func GetProtocRelease(ctx context.Context, version string) (*ProtocRelease, error) {
	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}
	workspace, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	suffix, err := inferProtocReleaseSuffix()
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("protoc-%s-%s.zip", version, suffix)
	url := fmt.Sprintf("https://github.com/protocolbuffers/protobuf/"+
		"releases/download/v%s/%s", version, filename)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &ErrHTTPDownload{
			Url: url,
			Err: err,
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &ErrHTTPDownload{
			Url: url,
			Err: err,
		}
	}
	zipFilePath := filepath.Join(workspace, filename)
	if err := downloadFile(resp, zipFilePath); err != nil {
		return nil, &ErrHTTPDownload{
			Url:  url,
			Err:  err,
			Code: resp.StatusCode,
		}
	}
	zip := archiver.NewZip()
	if err := zip.Unarchive(zipFilePath, workspace); err != nil {
		return nil, err
	}
	return &ProtocRelease{
		workspace: workspace,
	}, nil
}

// IsProtocInstalled is used to check whether the protoc version is installed
func IsProtocInstalled(ctx context.Context, storageDir string, version string) (bool, string, error) {
	local := PathForProtoc(storageDir, version)
	exists, err := util.IsFileExists(local)
	if err != nil {
		return false, "", err
	}
	return exists, local, nil
}

func inferProtocReleaseSuffix() (string, error) {
	goos := strings.ToLower(runtime.GOOS)
	arch := strings.ToLower(runtime.GOARCH)
	switch goos {
	case "linux":
		switch arch {
		case "arm64":
			return "linux-aarch_64", nil
		case "ppc64le":
			return "linux-ppcle_64", nil
		case "s390x":
			return "linux-s390_64", nil
		case "386":
			return "linux-x86_32", nil
		case "amd64":
			return "linux-x86_64", nil
		}
	case "darwin":
		return "osx-x86_64", nil
	case "windows":
		switch arch {
		case "386":
			return "win32", nil
		case "amd64":
			return "win64", nil
		}
	}
	return "", errors.New("protoc did not release on this platform")
}

func downloadFile(resp *http.Response, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), fs.ModePerm); err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.Errorf("unexpected code %d for url: %s", resp.StatusCode, resp.Request.URL.String())
	}
	file, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fs.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
