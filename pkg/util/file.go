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

package util

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	filecopy "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

// MatchPath is used to match path with specified pattern
func MatchPath(pattern string, path string) (bool, error) {
	return doublestar.PathMatch(pattern, path)
}

// CopyDirectory is used to copy directory
// If dst already exists, it will be merged
func CopyDirectory(src, dst string) error {
	return filecopy.Copy(src, dst)
}

// CopyFile is used to copy file from src to dst
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return errors.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	if err := os.MkdirAll(filepath.Dir(dst), fs.ModePerm); err != nil {
		return err
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// IsFileExists is used to check whether the file exists
func IsFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if info.IsDir() {
		return false, errors.Errorf("%s is not a file", path)
	}
	return true, nil
}

// IsDirExists is used to check whether the dir exists
func IsDirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !info.IsDir() {
		return false, errors.Errorf("%s is not a directory", path)
	}
	return true, nil
}

// GetFilesWithExtRecursively is used to recursively list files with a specific suffix
// expectExt should contain the prefix '.'
func GetFilesWithExtRecursively(target string, targetExt string) ([]string, error) {
	var data []string
	err := filepath.Walk(target, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == targetExt {
			data = append(data, path)
		}
		return nil
	})
	return data, err
}

// GetFilesWithExt is used to list files with a specific suffix
// expectExt should contain the prefix '.'
func GetFilesWithExt(dir string, targetExt string) ([]string, error) {
	children, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var data []string
	for _, child := range children {
		if child.IsDir() {
			continue
		}
		if ext := filepath.Ext(child.Name()); ext != targetExt {
			continue
		}
		data = append(data, filepath.Join(dir, child.Name()))
	}
	return data, nil
}
