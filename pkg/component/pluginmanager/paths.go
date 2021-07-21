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
	"path"
	"path/filepath"

	"golang.org/x/mod/module"

	"github.com/storyicon/powerproto/pkg/util"
)

// PathForInclude is used to get the local directory of include files
func PathForInclude(storageDir string) string {
	return filepath.Join(storageDir, "include")
}

// PathForProtoc is used to get the local binary location where the specified version protoc should be stored
func PathForProtoc(storageDir string, version string) string {
	return filepath.Join(storageDir, "protoc", version, util.GetBinaryFileName("protoc"))
}

// GetPluginPath is used to get the plugin path
func GetPluginPath(path string, version string) (string, error) {
	enc, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	encVer, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}
	return filepath.Join(enc + "@" + encVer), nil
}

// PathForGoogleAPIs is used to get the google apis path
func PathForGoogleAPIs(storageDir string, commitId string) string {
	return filepath.Join(storageDir, "googleapis", commitId)
}

// PathForPluginDir is used to get the local directory where the specified version plug-in should be stored
func PathForPluginDir(storageDir string, path string, version string) (string, error) {
	pluginPath, err := GetPluginPath(path, version)
	if err != nil {
		return "", err
	}
	return filepath.Join(storageDir, "plugins", pluginPath), nil
}

// PathForPlugin is used to get the binary path of plugin
func PathForPlugin(storageDir string, path string, version string) (string, error) {
	name := GetGoPkgExecName(path)
	dir, err := PathForPluginDir(storageDir, path, version)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, util.GetBinaryFileName(name)), nil
}

// isVersionElement reports whether s is a well-formed path version element:
// v2, v3, v10, etc, but not v0, v05, v1.
// `src\cmd\go\internal\load\pkg.go:1209`
func isVersionElement(s string) bool {
	if len(s) < 2 || s[0] != 'v' || s[1] == '0' || s[1] == '1' && len(s) == 2 {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || '9' < s[i] {
			return false
		}
	}
	return true
}

// GetGoPkgExecName is used to parse binary name from pkg uri
// `src\cmd\go\internal\load\pkg.go:1595`
func GetGoPkgExecName(pkgPath string) string {
	_, elem := path.Split(pkgPath)
	if elem != pkgPath && isVersionElement(elem) {
		_, elem = path.Split(path.Dir(pkgPath))
	}
	return elem
}
