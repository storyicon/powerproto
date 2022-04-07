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
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/ppaanngggg/powerproto/pkg/util/logger"
)

// defines a set of const value
const (
	// ConfigFileName defines the config file name
	ConfigFileName = "powerproto.yaml"
	// KeyNamePowerProtocInclude is the key name of powerproto default include
	KeyNamePowerProtocInclude = "POWERPROTO_INCLUDE"
	// The default include can be referenced by this key in import paths
	KeyPowerProtoInclude  = "$" + KeyNamePowerProtocInclude
	KeyNameSourceRelative = "SOURCE_RELATIVE"
	// KeySourceRelative can be specified in import paths to refer to
	// the folder where the current proto file is located
	KeySourceRelative = "$" + KeyNameSourceRelative
	// Defines the program directory of PowerProto, including various binary and include files
	EnvHomeDir = "POWERPROTO_HOME"
	// ProtobufRepository defines the protobuf repository
	ProtobufRepository = "https://github.com/protocolbuffers/protobuf"
	// GoogleAPIsRepository defines the google apis repository
	GoogleAPIsRepository = "https://github.com/googleapis/googleapis"
)

// defines a set of text style
var (
	TextExecuteAction  = color.HiGreenString("EXECUTE ACTION")
	TextExecuteCommand = color.HiGreenString("EXECUTE COMMAND")
	TextDryRun         = color.HiGreenString("DRY RUN")
)

var homeDir string
var log = logger.NewDefault("consts")

// GetHomeDir is used to get cached homeDir
func GetHomeDir() string {
	return homeDir
}

// PathForGlobalConfig is used to get path of global config
func PathForGlobalConfig() string {
	return filepath.Join(GetHomeDir(), ConfigFileName)
}

func getHomeDir() (string, error) {
	val := os.Getenv(EnvHomeDir)
	if val != "" {
		return filepath.Abs(val)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".powerproto"), nil
}

func init() {
	var err error
	homeDir, err = getHomeDir()
	if err != nil {
		log.LogFatal(nil, "Please set the working directory of PowerProto by "+
			"configuring the environment variable %s", EnvHomeDir)
	}
}
