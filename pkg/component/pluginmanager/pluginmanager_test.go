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

package pluginmanager_test

import (
	"context"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/storyicon/powerproto/pkg/component/pluginmanager"
	"github.com/storyicon/powerproto/pkg/util/logger"
)

var _ = Describe("Pluginmanager", func() {
	cfg := pluginmanager.NewConfig()
	cfg.StorageDir, _ = filepath.Abs("./tests")

	const pluginPkg = "google.golang.org/protobuf/cmd/protoc-gen-go"
	var manager pluginmanager.PluginManager
	It("should able to init", func() {
		pluginManager, err := pluginmanager.NewPluginManager(cfg, logger.NewDefault("pluginmanager"))
		Expect(err).To(BeNil())
		Expect(pluginManager).To(Not(BeNil()))
		manager = pluginManager
	})
	It("should able to install protoc", func() {
		versions, err := manager.ListProtocVersions(context.TODO())
		Expect(err).To(BeNil())
		Expect(len(versions) > 0).To(BeTrue())
		latestVersion, err := manager.GetProtocLatestVersion(context.TODO())
		Expect(err).To(BeNil())
		Expect(latestVersion).To(Equal(versions[len(versions)-1]))

		local, err := manager.InstallProtoc(context.TODO(), latestVersion)
		Expect(err).To(BeNil())
		exists, local, err := manager.IsProtocInstalled(context.TODO(), latestVersion)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
		Expect(len(local) != 0).To(BeTrue())
	})
	It("should able to install plugin", func() {
		versions, err := manager.ListPluginVersions(context.TODO(), pluginPkg)
		Expect(err).To(BeNil())
		Expect(len(versions) > 0).To(BeTrue())

		latestVersion, err := manager.GetPluginLatestVersion(context.TODO(), pluginPkg)
		Expect(err).To(BeNil())
		Expect(latestVersion).To(Equal(versions[len(versions)-1]))

		local, err := manager.InstallPlugin(context.TODO(), pluginPkg, latestVersion)
		Expect(err).To(BeNil())
		Expect(len(local) > 0).To(BeTrue())

		exists, local, err := manager.IsPluginInstalled(context.TODO(), pluginPkg, latestVersion)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
		Expect(len(local) != 0).To(BeTrue())
	})
	It("should able to install googleapis", func() {
		versions, err := manager.ListGoogleAPIsVersions(context.TODO())
		Expect(err).To(BeNil())
		Expect(len(versions) > 0).To(BeTrue())

		latestVersion, err := manager.GetGoogleAPIsLatestVersion(context.TODO())
		Expect(err).To(BeNil())
		Expect(latestVersion).To(Equal(versions[len(versions)-1]))

		local, err := manager.InstallGoogleAPIs(context.TODO(), latestVersion)
		Expect(err).To(BeNil())
		Expect(len(local) > 0).To(BeTrue())

		exists, local, err := manager.IsGoogleAPIsInstalled(context.TODO(), latestVersion)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
		Expect(len(local) != 0).To(BeTrue())
	})
})
