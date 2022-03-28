package util

import "testing"

func TestGetCommonRootDirOfPaths(t *testing.T) {
	type args struct {
		filepaths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{filepaths: []string{"/root/a/a.proto", "/root/a/c.proto"}},
			want: "/root/a",
		},
		{
			args: args{filepaths: []string{"/root/a/a.proto", "/root/a/c.proto", "/root/b/b.proto"}},
			want: "/root",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCommonRootDirOfPaths(tt.args.filepaths); got != tt.want {
				t.Errorf("GetCommonRootDirOfPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
