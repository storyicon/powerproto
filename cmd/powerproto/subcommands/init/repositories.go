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

package build

import (
	"fmt"
	"strings"
)

// Repository defines the plugin options
type Repository struct {
	OptionsValue string

	Name        string
	Pkg         string
	ImportPaths []string
}

// OptionsValue is used to return the options value
func (repo *Repository) GetOptionsValue() string {
	if repo.OptionsValue != "" {
		return repo.OptionsValue
	}
	return fmt.Sprintf("%s: %s", strings.ToLower(repo.Name), repo.Pkg)
}

// GetWellKnownRepositories is used to get well known plugins
func GetWellKnownRepositories() []*Repository {
	return []*Repository{
		GetRepositoryGoogleAPIs(),
		GetRepositoryGoGoProtobuf(),
	}
}

// GetRepositoryGoGoProtobuf is used to get gogo protobuf repository
func GetRepositoryGoGoProtobuf() *Repository {
	return &Repository{
		Name: "GOGO_PROTOBUF",
		Pkg:  "https://github.com/gogo/protobuf@226206f39bd7276e88ec684ea0028c18ec2c91ae",
		ImportPaths: []string{
			"$GOGO_PROTOBUF",
		},
	}
}

// GetRepositoryGoogleAPIs is used to get google apis repository
func GetRepositoryGoogleAPIs() *Repository {
	return &Repository{
		Name: "GOOGLE_APIS",
		Pkg:  "https://github.com/googleapis/googleapis@75e9812478607db997376ccea247dd6928f70f45",
		ImportPaths: []string{
			"$GOOGLE_APIS/github.com/googleapis/googleapis",
		},
	}
}

// GetRepositoryFromOptionsValue is used to get plugin by option value
func GetRepositoryFromOptionsValue(val string) (*Repository, bool) {
	repositories := GetWellKnownRepositories()
	for _, repo := range repositories {
		if repo.GetOptionsValue() == val {
			return repo, true
		}
	}
	return nil, false
}

// GetWellKnownRepositoriesOptionValues is used to get option values of well known plugins
func GetWellKnownRepositoriesOptionValues() []string {
	repos := GetWellKnownRepositories()
	packages := make([]string, 0, len(repos))
	for _, repo := range repos {
		packages = append(packages, repo.GetOptionsValue())
	}
	return packages
}
