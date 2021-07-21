# PowerProto

[![Go Report Card](https://goreportcard.com/badge/github.com/storyicon/powerproto?t=1)](https://goreportcard.com/report/github.com/storyicon/powerproto) ![TotalLine](https://img.shields.io/tokei/lines/github/storyicon/powerproto?color=77%2C199%2C31) ![last-commit](https://img.shields.io/github/last-commit/storyicon/powerproto) [![GoDoc](https://godoc.org/github.com/storyicon/powerproto?status.svg)](https://godoc.org/github.com/storyicon/powerproto) 

**中文** | [English](README_CN.md)

![exmpales](./docs/images/exmaple.gif)    

PowerProto主要用于解决下面三个问题：

1. 降低gRPC的使用门槛与使用成本。
2. 解决protoc以及其相关插件（比如protoc-gen-go、protoc-gen-grpc-gateway）的版本控制问题。
3. 高效管理proto的编译，实现多平台兼容、一键安装与编译。


- [PowerProto](#powerproto)
  - [功能](#功能)
  - [安装与依赖](#安装与依赖)
    - [一、通过Go进行安装](#一通过go进行安装)
    - [二、开箱即用版本](#二开箱即用版本)
  - [命令介绍](#命令介绍)
    - [一、初始化配置](#一初始化配置)
    - [二、整理配置](#二整理配置)
    - [三、编译Proto文件](#三编译proto文件)
    - [四、查看环境变量](#四查看环境变量)
  - [示例](#示例)
  - [配置文件](#配置文件)
    - [解释](#解释)
      - [匹配模式与工作目录](#匹配模式与工作目录)
      - [多配置组合](#多配置组合)
    - [PostAction](#postaction)
      - [1. copy](#1-copy)
      - [2. move](#2-move)
      - [3. remove](#3-remove)
      - [4. replace](#4-replace)


## 功能

1. 实现protoc的一键安装与多版本管理。
2. 实现protoc相关插件（比如protoc-gen-go）的一键安装与多版本管理。
3. 通过配置文件管理proto的编译，而非shell脚本，提高可读性与兼容性。
4. 引导式生成配置文件，跨平台兼容，一份配置在多个平台均可以实现一键编译。
5. 支持批量、递归编译proto文件，提高效率。
6. 跨平台支持PostAction，可以在编译完成之后执行一些常规操作（比如替换掉所有生成文件中的"omitempty"）。
7. 支持PostShell，在编译完成之后执行特定的shell脚本。
8. 支持 `google api` 的一键安装与版本控制。

## 安装与依赖

1. 目前版本的 `PowerProto` 依赖 `go` 以及 `git`（未来可能会直接使用CDN拉取构建好的二进制），请确保运行环境中包含这两个命令。
2. `protoc`的下载源是Github，`PowerProto`在下载`protoc`时尊重 `HTTP_PROXY`、`HTTPS_PROXY`环境变量，如果遇到网络问题，请自行配置代理。
3. 在查询`protoc`的版本列表时，会对`github.com`使用`git ls-remote`，如果遇到网络问题，请自行为`git`配置代理。
4. 在当前版本，下载和查询插件版本均依赖`go`命令，所以如果遇到网络问题，请自行配置 `GOPROXY`环境变量。
5. 默认会使用 `用户目录/.powerproto`作为安装目录，用于放置下载好的各种插件以及全局配置，可以通过 `POWERPROTO_HOME`环境变量来修改安装目录。
6. 如果认为`powerproto`名字太长，可以通过`alias`成一个更简单的名字来提高输入效率，比如没有人会介意你叫它`pp`。


### 一、通过Go进行安装

直接执行下面的命令即可进行安装：

```
go install github.com/storyicon/powerproto/cmd/powerproto@latest
```

### 二、开箱即用版本

可以通过 `Github Release Page` 下载开箱即用版本。

## 命令介绍

你可以通过 powerproto -h 来查看帮助，比如：

```
powerproto -h
powerproto init -h
powerproto tidy -h
powerproto build -h
```

它的好处是命令行中的文档永远和你的二进制版本保持一致。而Github上的文档可能会一直是对应最新的二进制。

### 一、初始化配置

可以通过下面的命令进行配置的初始化：

```
powerproto init
```

### 二、整理配置

可以通过下面的命令整理配置：

```
powerproto tidy
```

它将会从当前目录开始向父级目录搜索名为 `powerproto.yaml` 的配置文件，并对配置进行读取和整理。

你也可以指定整理哪个配置文件：

```
powerproto tidy [the path of proto file]
```

整理配置主要包含两个操作：

1. 通过查询，将版本中的latest替换为真实的最新版本号。
2. 安装配置文件中定义的所有依赖。

支持通过 `-d` 参数来进入到`debug模式`，查看更详细的日志。

### 三、编译Proto文件

可以通过下面的命令进行Proto文件的编译：

```
// 编译指定的proto文件
powerproto build xxxx.proto

// 编译当前目录下的所有proto文件
powerproto build .

// 递归编译当前目录下的所有proto文件，包括子文件夹。
powerproto build -r .
```

其执行逻辑是，对于每一个proto文件，从其文件所在目录开始向父级目录寻找 `powerproto.yaml` 配置文件:

1. 对于找到的配置文件，与其中的`scope`进行匹配，如果匹配则采用。
2. 检查并安装配置文件中声明的依赖。
3. 根据配置文件中的`plugins`、`protoc`、`options`、`importPaths`等配置对proto文件进行编译。 当所有的proto文件都编译完成之后，如果你指定了 `-p` 参数，还会进行`PostAction`与`PostShell`的执行。

注意：`protoc`执行的工作目录默认是`proto文件`匹配到的配置文件所在的目录，它相当于你在配置文件所在目录执行protoc命令。你可以通过配置文件中的 `protocWorkDir` 来进行修改。

支持通过 `-d` 参数来进入到`debug模式`，查看更详细的日志。
支持通过 `-y` 参数来进入到`dryRun模式`，只打印命令而不真正执行，这对于调试非常有用。

### 四、查看环境变量

如果你的命令一直卡在某个状态，大概率是出现网络问题了。
你可以通过下面的命令来查看环境变量是否配置成功：
```
powerproto env
```

## 示例

比如你在 `/mnt/data/hello` 目录下拥有下面这样的文件结构：

```
$ pwd
/mnt/data/hello

$ tree
./apis
└── hello.proto
powerproto.yaml
```

`powerproto.yaml` 的文件内容是（你可以通过 `powerproto init` 命令很方便的生成配置文件）：

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

在任意目录执行：

```
powerproto build -r /mnt/data/hello/apis
```

你都可以得到编译后的文件

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

它相当于你在 `powerproto.yaml` 所在目录，执行：

```
$POWERPROTO_HOME/protoc/3.17.3/protoc --go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
--proto_path=/mnt/data/hello \
--proto_path=$GOPATH \
--proto_path=$POWERPROTO_HOME/include \
--plugin=protoc-gen-go=$POWERPROTO_HOME/plugins/google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1/protoc-gen-go \
--plugin=protoc-gen-go-grpc=$POWERPROTO_HOME/plugins/google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0/protoc-gen-go-grpc
/mnt/data/hello/apis/hello.proto
```

## 配置文件

配置文件用于描述编译proto文件时，各种依赖的版本以及参数等。

可以方便的通过 `powerproto init`进行配置文件的初始化。

### 解释

以下面这份配置文件为例：


```yaml
# 必填，scopes 用于定义作用域，即当前配置项对项目中的哪些目录生效
scopes:
    - ./
# 必填，protoc的版本，可以填 latest，会自动转换成最新的版本
protoc: 3.17.3
# 选填，执行protoc命令的工作目录，默认是配置文件所在目录
# 支持路径中混用环境变量，比如$GOPATH
protocWorkDir: ""
# 选填，如果需要使用 googleapis，你应该在这里填写googleapis的commit id
# 可以填 latest，会自动转换成最新的版本
googleapis: 75e9812478607db997376ccea247dd6928f70f45
# 必填，代表scope匹配的目录中的proto文件，在编译时需要用到哪些插件
plugins:
    # 插件的名字、路径以及版本号。
    # 插件的地址必须是 path@version 的格式，version可以填latest，会自动转换成最新的版本。
    protoc-gen-deepcopy: istio.io/tools/cmd/protoc-gen-deepcopy@latest
    protoc-gen-go: google.golang.org/protobuf/cmd/protoc-gen-go@latest
    protoc-gen-go-json: github.com/mitchellh/protoc-gen-go-json@v1.0.0
    protoc-gen-grpc-gateway: github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
# 必填，定义了编译proto文件时 protoc 的参数
options:
    - --go_out=paths=source_relative:.
    - --go-json_out=.
    - --deepcopy_out=source_relative:.
    - --grpc-gateway_out=.
    - --go-grpc_out=paths=source_relative:.
# 必填，定义了构建时 protoc 的引用路径，会被转换为 --proto_path (-I) 参数。
importPaths:
    # 特殊变量。代表当前配置文件所在文件夹
    - .
    # 环境变量。可以使用环境变量
    # 也支持 $GOPATH/include 这样的混合写法
    - $GOPATH
    - $POWERPROTO_INCLUDE
    # 特殊变量。引用待编译的proto文件所在的目录
    # 比如将要编译 /a/b/data.proto，那么 /a/b 目录将会被自动引用
    - $SOURCE_RELATIVE
    # 特殊变量。引用googleapis字段所指定的版本的google apis
    - $POWERPROTO_GOOGLEAPIS
# 选填，构建完成之后执行的操作，工作目录是配置文件所在目录
# postActions是跨平台兼容的
# 注意，必须在 powerproto build 时附加 -p 参数，才会执行配置文件中的postActions
postActions: []
# 选填，构建完成之后执行的shell脚本，工作目录是配置文件所在目录
# postShell不是跨平台兼容的。
# 注意，必须在 powerproto build 时附加 -p 参数，才会执行配置文件中的postShell
postShell: |
    // do something
```


#### 匹配模式与工作目录

在构建proto文件时，将会从proto文件所在目录开始，向父级目录搜索 `powerproto.yaml` 配置文件，并与其中的 scope进行匹配，第一个匹配到的配置，将会被用于此proto文件的编译。
在执行protoc时（执行postActions、postShell时也是如此），是以配置文件所在目录作为工作目录的，即相当于你在这个目录执行protoc命令。


#### 多配置组合

一个配置文件中支持填写多份配置，多份配置之间以 "---" 进行分割。

在下面的示例中，apis1目录使用的是v1.25.0的protoc-gen-go，而apis2目录使用的则是v1.27.0的protoc-gen-go。

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

PostAction允许在所有的proto文件都编译完成之后，执行特定的操作。与`PostShell`相比，它是跨平台支持的。

为了安全起见，只有在执行 `powerproto build`时附加上 `-p` 参数，才会执行配置文件中定义的`PostAction`与`PostShell`。

目前，PostAction支持下面这些命令：

| 命令    | 描述                   | 函数原型                                              |
| ------- | ---------------------- | ----------------------------------------------------- |
| copy    | 复制文件或文件夹       | copy(src string, dest string) error                   |
| move    | 移动文件或文件夹       | move(src string, dest string) error                   |
| remove  | 删除文件或文件夹       | remove(path ...string) error                          |
| replace | 批量替换文件中的字符串 | replace(pattern string, from string, to string) error |

#### 1. copy

用于复制文件或文件夹，其函数原型为：

```
copy(src string, dest string) error
```

为了安全以及配置的兼容性，参数中只允许填写相对路径。

如果目标文件夹已经存在，将会合并。

下面的例子将会把配置文件所在目录下的a复制到b：

```yaml
postActions:
    - name: copy
      args:
        - ./a
        - ./b
```

#### 2. move

用于移动文件或文件夹，其函数原型为：

```
move(src string, dest string) error
```

为了安全以及配置的兼容性，参数中只允许填写相对路径。

如果目标文件夹已经存在，将会合并。

下面的例子将会把配置文件所在目录下的a移动到b：

```yaml
postActions:
    - name: move
      args:
        - ./a
        - ./b
```

#### 3. remove

用于删除文件或文件夹，其函数原型为：

```
remove(path ...string) error
```

为了安全以及配置的兼容性，参数中只允许填写相对路径。

下面的例子将会删除配置文件所在目录下的a、b、c：

```yaml
postActions:
    - name: remove
      args:
        - ./a
        - ./b
        - ./c
```

#### 4. replace

用于批量替换文件中的字符串，其函数原型为：

```
replace(pattern string, from string, to string) error
```

其中：

* pattern是支持通配符的相对路径。
* from是要被替换的字符串。
* to是替换为的字符串。

下面的例子将会把apis目录以及其子目录下的所有go文件中的 `,omitempty` 替换为空字符串：

```
postActions:
    - name: replace
      args:
        - ./apis/**/*.go
        - ',omitempty'
        - ""
```





