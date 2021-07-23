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
	"path/filepath"
	"reflect"
	"testing"
)

func TestRenderPathWithEnv(t *testing.T) {
	type args struct {
		path string
		ext  map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				path: "$POWERPROTO_INCLUDE/protobuf",
				ext: map[string]string{
					"POWERPROTO_INCLUDE": "/mnt/powerproto/include",
				},
			},
			want: filepath.Clean("/mnt/powerproto/include/protobuf"),
		},
		{
			args: args{
				path: "$POWERPROTO_INCLUDE/protobuf",
			},
			want: filepath.Clean("$POWERPROTO_INCLUDE/protobuf"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderPathWithEnv(tt.args.path, tt.args.ext); got != tt.want {
				t.Errorf("RenderPathWithEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeduplicateSliceStably(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				items: []string{"a", "B", "c", "B"},
			},
			want: []string{"a", "B", "c"},
		},
		{
			args: args{
				items: []string{"B", "c", "B", "a"},
			},
			want: []string{"B", "c", "a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeduplicateSliceStably(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeduplicateSliceStably() = %v, want %v", got, tt.want)
			}
		})
	}
}