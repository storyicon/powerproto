# PowerProto

[![Go Report Card](https://goreportcard.com/badge/github.com/storyicon/powerproto?t=2)](https://goreportcard.com/report/github.com/storyicon/powerproto) ![TotalLine](https://img.shields.io/tokei/lines/github/storyicon/powerproto?color=77%2C199%2C31) ![last-commit](https://img.shields.io/github/last-commit/storyicon/powerproto) [![GoDoc](https://godoc.org/github.com/storyicon/powerproto?status.svg)](https://godoc.org/github.com/storyicon/powerproto) 

**English** | [中文](README_CN.md)

![exmpales](./docs/images/exmaple.gif)

PowerProto is used to solve the following three main problems:

1. lower the usage threshold and usage cost of gRPC.
2. solve the version control problem of protoc and its related plugins (such as protoc-gen-go, protoc-gen-grpc-gateway).
3. efficiently manage the compilation of proto to achieve multi-platform compatibility, one-click installation and compilation.


- [PowerProto](#powerproto)
  - [🎉 Features](#-features)
  - [Installation and Dependencies](#installation-and-dependencies)
    - [I. Installation via Go](#i-installation-via-go)
    - [II. out-of-the-box version](#ii-out-of-the-box-version)
  - [Command Introduction](#command-introduction)
    - [I. Initial Config](#i-initial-config)
    - [II. Tidy Config](#ii-tidy-config)
    - [III. Compiling Proto files](#iii-compiling-proto-files)
    - [IV. View environment variables](#iv-view-environment-variables)
  - [Examples](#examples)
  - [Config File](#config-file)
    - [Definition](#definition)
      - [Matching patterns and working directory](#matching-patterns-and-working-directory)
      - [Multi-config](#multi-config)
    - [PostAction](#postaction)
      - [1. copy](#1-copy)
      - [2. move](#2-move)
      - [3. remove](#3-remove)
      - [4. replace](#4-replace)


## 🎉 Features

1. one-click installation and multi-version management of protoc.
2. one-click installation and multi-version management of protoc related plugins (such as protoc-gen-go).
3. manage the compilation of proto through config file instead of shell script to improve readability and compatibility.
4. bootstrap generation of config files, cross-platform compatibility, a config can be compiled in multiple platforms with one click.
5. support batch and recursive compilation of proto files to improve efficiency.
6. cross-platform support PostAction, you can perform some routine operations (such as replacing "omitempty" in all generated files) after the compilation.
7. support PostShell, execute specific shell scripts after the compilation.
8. one-click installation and version control of google apis。

## Installation and Dependencies

1. The current version of `PowerProto` relies on `go` and `git` (in the future it may use CDN to pull built binaries directly), please make sure the runtime environment contains these two commands.
2. `protoc` download source is Github, `PowerProto` respects `HTTP_PROXY`, `HTTPS_PROXY` environment variables when downloading `protoc`, if you encounter network problems, please configure your own proxy.
3. When querying the version list of `protoc`, `git ls-remote` is used for `github.com`, if you encounter network problems, please configure the proxy for `git` by yourself.
4. In the current version, downloading and querying plugin versions rely on the `go` command, so if you encounter network problems, please configure the `GOPROXY` environment variable yourself.
5. By default, `user directory/.powerproto` is used as the installation directory, which is used to place the downloaded plug-ins and global config.
6. If you think the name `powerproto` is too long, you can `alias` it into a simpler name to improve the input efficiency, for example, no one will mind if you call it `pp`.


### I. Installation via Go

Installation can be performed by executing the following command directly:

```
go install github.com/storyicon/powerproto/cmd/powerproto@latest
```

### II. out-of-the-box version

You can download the out-of-the-box version via the [`Github Release Page`](https://github.com/storyicon/powerproto/releases)


## Command Introduction

You can view help with `powerproto -h`, e.g.

```
powerproto -h
powerproto init -h
powerproto tidy -h
powerproto build -h
powerproto env -h
```

It has the advantage that the documentation on the command line is always consistent with your binary version.

### I. Initial Config

The config can be initialized with the following command.

```
powerproto init
```

### II. Tidy Config

The config can be tidied with the following command.

```
powerproto tidy
```

It will search for a config file named `powerproto.yaml` from the current directory to the parent directory, and will read and tidy the config.

You can also specify which config file to tidy.

```
powerproto tidy [the path of proto file]
```

Tidy the config consists of two main operations:

1. replacing the latest in the version with the real latest version number by querying.
2. install all dependencies defined in the config file.


Supports entering `debug mode` by appending the `-d` argument to see more detailed logs.


### III. Compiling Proto files

The Proto file can be compiled with the following command.

```
// Compile the specified proto file
powerproto build xxxx.proto

// Compile all the proto files in the current directory
powerproto build .

// Compile all proto files in the current directory recursively, including subfolders.
powerproto build -r .
```

The execution logic is that for each proto file, the `powerproto.yaml` config file will be searched from the directory where the proto file is located to the ancestor directory:

1. For the found config file, match it with the `scope` in it and use it if it matches.
2. Check and install the dependencies declared in the config file.
3. Compile the proto file according to the `plugins`, `protoc`, `options`, `importPaths` and other configs in the config file。 After all the proto files are compiled, if you specify the `-p` argument, `PostAction` and `PostShell` will also be executed.

Note: The default `working directory` of `PowerProto` is the directory where the `proto file` matches to the config file, it is equivalent to the directory where you execute the `protoc` command. You can change it via `protocWorkDir` in the config file.


Supports entering `debug mode` by appending the `-d` argument to see more detailed logs.

Supports entering `dryRun mode` by appending the `-y` argument, in this mode the commands are not actually executed, but just printed out, which is very useful for debugging.

### IV. View environment variables

If your command keeps getting stuck in a certain state, there is a high probability that there is a network problem.        

You can check if the environment variables are configured successfully with the following command:

```
powerproto env
```


## Examples

For example, you have the following file structure in the `/mnt/data/hello` directory:

```
$ pwd
/mnt/data/hello

$ tree
./apis
└── hello.proto
powerproto.yaml
```

The contents of the `powerproto.yaml` file (you can easily generate the config file with the `powerproto init` command) are:

```
scopes:
    - ./
protoc: latest
protocWorkDir: ""
plugins:
    protoc-gen-go: google.golang.org/protobuf/cmd/protoc-gen-go@latest
    protoc-gen-go-grpc: google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
options:
    - --go_out=.
    - --go_opt=paths=source_relative
    - --go-grpc_out=.
    - --go-grpc_opt=paths=source_relative
importPaths:
    - .
    - $GOPATH
    - $POWERPROTO_INCLUDE
postActions: []
postShell: ""
```

Execute in any directory:

```
powerproto build -r /mnt/data/hello/apis
```

You can get the compiled file:

```
$ pwd
/mnt/data/hello

$ tree
./apis
├── hello.pb.go
├── hello.proto
└── hello_grpc.pb.go
powerproto.yaml
```

It is equivalent to if you were in the directory where `powerproto.yaml` is located and executed:

```shell script
$POWERPROTO_HOME/protoc/3.17.3/protoc --go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
--proto_path=/mnt/data/hello \
--proto_path=$GOPATH \
--proto_path=$POWERPROTO_HOME/include \
--plugin=protoc-gen-go=$POWERPROTO_HOME/plugins/google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1/protoc-gen-go \
--plugin=protoc-gen-go-grpc=$POWERPROTO_HOME/plugins/google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0/protoc-gen-go-grpc \
/mnt/data/hello/apis/hello.proto
```

## Config File

The config file is used to describe the versions of various dependencies and parameters when compiling the proto file.

It can be easily initialized with `powerproto init`.

### Definition

Take the following config file as an example:

```yaml
# required. scopes is used to define scopes. 
# i.e. which directories in the project the current config item is valid for
scopes:
    - ./
# required. the version of protoc.
# you can fill in the 'latest', will be automatically converted to the latest version
protoc: 3.17.3
# optional. The working directory for executing the protoc command, 
# the default is the directory where the config file is located.
# support mixed environment variables in path, such as $GOPATH
protocWorkDir: ""
# optional. If you need to use googleapis, you should fill in the commit id of googleapis here.
# You can fill in the latest, it will be automatically converted to the latest version.
googleapis: 75e9812478607db997376ccea247dd6928f70f45
# required. it is used to describe which plug-ins are required for compilation
plugins:
    # the name, path, and version number of the plugin.
    # the address of the plugin must be in path@version format, 
    # and version can be filled with 'latest', which will be automatically converted to the latest version.
    protoc-gen-deepcopy: istio.io/tools/cmd/protoc-gen-deepcopy@latest
    protoc-gen-go: google.golang.org/protobuf/cmd/protoc-gen-go@latest
    protoc-gen-go-json: github.com/mitchellh/protoc-gen-go-json@v1.0.0
    protoc-gen-grpc-gateway: github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
# required. defines the parameters of protoc when compiling proto files
options:
    - --go_out=paths=source_relative:.
    - --go-json_out=.
    - --deepcopy_out=source_relative:.
    - --grpc-gateway_out=.
    - --go-grpc_out=paths=source_relative:.
# required. defines the path of the proto dependency, which will be converted to the --proto_path (-I) parameter.
importPaths:
    # Special variables. Will be replaced with the folder where the current configuration file is located.
    - .
    # Environment variables. Environment variables can be used in importPaths.
    # Support mixed writing like $GOPATH/include
    - $GOPATH
    # Special variables. Will be replaced with the local path to the public proto file that comes with protoc by default
    - $POWERPROTO_INCLUDE
    # Special variables. Reference to the directory where the proto file to be compiled is located
    # For example, if /a/b/data.proto is to be compiled, then the /a/b directory will be automatically referenced
    - $SOURCE_RELATIVE
    # Special variables. Will be replaced with the local path to the version of google apis specified by the googleapis field
    - $POWERPROTO_GOOGLEAPIS
# optional. The operation is executed after compilation.
# its working directory is the directory where the config file is located.
# postActions is cross-platform compatible.
# Note that the "-p" parameter must be appended to the "powerproto build" to allow execution of the postActions in the config file
postActions: []
# optional. The shell script that is executed after compilation.
# its working directory is the directory where the config file is located.
# postShell is not cross-platform compatible.
# Note that the "-p" parameter must be appended to the "powerproto build" to allow execution of the postShell in the config file
postShell: |
    // do something
```

#### Matching patterns and working directory

When building the proto file, the `powerproto.yaml` config file will be searched from the directory where the proto file is located to the ancestor directory, match with the `scope` in.
The first matched config item will be used for the compilation of this proto file.
When PowerProto executes protoc (and also when it executes postActions and postShell), the default is to use the directory where the config file is located as the working directory. (working directory is equivalent to the directory where you execute the protoc command.)

#### Multi-config

A config file can be filled with multiple configs, which are separated by "---".

In the example below, the apis1 directory uses protoc-gen-go with v1.25.0, while the apis2 directory uses protoc-gen-go with v1.27.0.

```
scopes:
    - ./apis1
protoc: v3.17.3
protocWorkDir: ""
googleapis: 75e9812478607db997376ccea247dd6928f70f45
plugins:
    protoc-gen-go: google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
    protoc-gen-go-grpc: google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
options:
    - --go_out=.
    - --go_opt=paths=source_relative
    - --go-grpc_out=.
    - --go-grpc_opt=paths=source_relative
importPaths:
    - .
    - $GOPATH
    - $POWERPROTO_INCLUDE
postActions: []
postShell: ""

---

scopes:
    - ./apis2
protoc: v3.17.3
protocWorkDir: ""
googleapis: 75e9812478607db997376ccea247dd6928f70f45
plugins:
    protoc-gen-go: google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.0
    protoc-gen-go-grpc: google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
options:
    - --go_out=.
    - --go_opt=paths=source_relative
    - --go-grpc_out=.
    - --go-grpc_opt=paths=source_relative
importPaths:
    - .
    - $GOPATH
    - $POWERPROTO_INCLUDE
postActions: []
postShell: ""
```



### PostAction

PostAction allows to perform specific actions after all proto files have been compiled. In contrast to `PostShell`, it is cross-platform supported.
For security reasons, the `PostAction` and `PostShell` defined in the config file will only be executed if the `-p` argument is appended to the execution of `powerproto build`.


Currently, PostAction supports the following commands:

| Command    | Description                  | Function Prototype                                              |
| ------- | ---------------------- | ----------------------------------------------------- |
| copy    | Copy file or folder       | copy(src string, dest string) error                   |
| move    | Move file or folder       | move(src string, dest string) error                   |
| remove  | Delete file or folder       | remove(path ...string) error                          |
| replace | Batch replacement of strings in files | replace(pattern string, from string, to string) error |

#### 1. copy

For copying files or folders, the function prototype is:

```
copy(src string, dest string) error
```

For security and config compatibility, only relative paths are allowed in the parameters.

If the target folder already exists, it will be merged.

The following example will copy 'a' from the directory where the config file is located to 'b'.

```yaml
postActions:
    - name: copy
      args:
        - ./a
        - ./b
```

#### 2. move

For moving files or folders, the function prototype is:

```
move(src string, dest string) error
```

For security and config compatibility, only relative paths are allowed in the parameters.

If the target folder already exists, it will be merged.

The following example will move 'a' in the directory where the config file is located to 'b'.

```yaml
postActions:
    - name: move
      args:
        - ./a
        - ./b
```

#### 3. remove

For deleting files or folders, the function prototype is:

```
remove(path ...string) error
```

For security and config compatibility, only relative paths are allowed in the parameters.

The following example will remove 'a', 'b' and 'c' from the directory where the config file is located:

```yaml
postActions:
    - name: remove
      args:
        - ./a
        - ./b
        - ./c
```

#### 4. replace

Used for batch replacement of strings in files. Its function prototype is:

```
replace(pattern string, from string, to string) error
```

* `pattern` is a relative path that supports wildcard characters.
* `from` is the string to be replaced.
* `to` is the string to replace with.

The following example will replace `,omitempty` with the empty string in all go files in the apis directory and its subdirectories:

```
postActions:
    - name: replace
      args:
        - ./apis/**/*.go
        - ',omitempty'
        - ""
```





