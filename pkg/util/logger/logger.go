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

package logger

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus"
)

// Logger defines the basic log library implementation
type Logger interface {
	// LogDebug print a message with debug level.
	LogDebug(fields map[string]interface{}, format string, args ...interface{})
	// LogInfo print a message with info level.
	LogInfo(fields map[string]interface{}, format string, args ...interface{})
	// LogWarn print a message with warn level.
	LogWarn(fields map[string]interface{}, format string, args ...interface{})
	// LogError print a message with error level.
	LogError(fields map[string]interface{}, format string, args ...interface{})
	// LogFatal print a message with fatal level.
	LogFatal(fields map[string]interface{}, format string, args ...interface{})
	// NewLogger is used to derive a new child Logger
	NewLogger(component string) Logger
	// SetLogLevel is used to set log level
	SetLogLevel(level Level) Logger
}

// BasicLogger simply implements Logger
type BasicLogger struct {
	cfg *Config

	component  string
	registerer prometheus.Registerer
}

// Config defines the config structure
type Config struct {
	Pretty bool
	Level  Level
}

// Level is the structure of log level
type Level struct {
	Name   string
	Color  color.Attribute
	Index  int
	Writer io.Writer
}

// defines a set of log level
var (
	LevelDebug = Level{Name: "debug", Color: color.FgWhite, Index: 0, Writer: os.Stdout}
	LevelInfo  = Level{Name: "info", Color: color.FgWhite, Index: 1, Writer: os.Stdout}
	LevelWarn  = Level{Name: "warn", Color: color.FgYellow, Index: 2, Writer: os.Stderr}
	LevelError = Level{Name: "error", Color: color.FgHiRed, Index: 3, Writer: os.Stderr}
	LevelFatal = Level{Name: "fatal", Color: color.FgRed, Index: 4, Writer: os.Stderr}
)

// NewConfig is used to init config with default values
func NewConfig() *Config {
	return &Config{
		Pretty: false,
		Level:  LevelDebug,
	}
}

// NewDefault is used to initialize a simple Logger
func NewDefault(component string) Logger {
	logger, err := New(NewConfig(), component, prometheus.DefaultRegisterer)
	if err != nil {
		panic(err)
	}
	return logger
}

// New is used to init service
func New(cfg *Config, component string, registerer prometheus.Registerer) (Logger, error) {
	if cfg == nil {
		cfg = NewConfig()
	}
	service := &BasicLogger{
		cfg:        cfg,
		component:  component,
		registerer: registerer,
	}
	return service, nil
}

// LogDebug print a message with debug level.
func (b *BasicLogger) LogDebug(fields map[string]interface{}, format string, args ...interface{}) {
	b.log(LevelDebug, fields, format, args...)
}

// LogInfo print a message with info level.
func (b *BasicLogger) LogInfo(fields map[string]interface{}, format string, args ...interface{}) {
	b.log(LevelInfo, fields, format, args...)
}

// LogWarn print a message with warn level.
func (b *BasicLogger) LogWarn(fields map[string]interface{}, format string, args ...interface{}) {
	b.log(LevelWarn, fields, format, args...)
}

// LogError print a message with error level.
func (b *BasicLogger) LogError(fields map[string]interface{}, format string, args ...interface{}) {
	b.log(LevelError, fields, format, args...)
}

// LogFatal print a message with fatal level.
func (b *BasicLogger) LogFatal(fields map[string]interface{}, format string, args ...interface{}) {
	b.log(LevelFatal, fields, format, args...)
}

// NewLogger is used to derive a new child Logger
func (b *BasicLogger) NewLogger(component string) Logger {
	name := strings.Join([]string{b.component, component}, ".")
	logger, err := New(b.cfg, name, b.registerer)
	if err != nil {
		b.LogWarn(map[string]interface{}{
			"name": name,
		}, "failed to extend logger: %s", err)
		return b
	}
	return logger
}

// SetLogLevel is used to set log level
func (b *BasicLogger) SetLogLevel(level Level) Logger {
	b.cfg.Level = level
	return b
}

func (b *BasicLogger) log(level Level, fields map[string]interface{}, format string, args ...interface{}) {
	if b.cfg.Level.Index > level.Index {
		return
	}
	if fields == nil {
		fields = map[string]interface{}{}
	}
	// if b.cfg.Level == LevelDebug {
	// 	if _, file, line, ok := runtime.Caller(4); ok {
	// 		fields["file"] = fmt.Sprintf("%s:%d", path.Base(file), line)
	// 	}
	// 	if b.component != "" {
	// 		fields["component"] = b.component
	// 	}
	// }
	if b.cfg.Pretty {
		dict := map[string]interface{}{}
		for key, val := range fields {
			dict[key] = val
		}
		dict["level"] = level.Name
		dict["message"] = fmt.Sprintf(format, args...)
		_ = jsoniter.NewEncoder(level.Writer).Encode(dict)
		return
	}
	var buf []byte
	buf = appendString(buf, fmt.Sprintf(format, args...))
	if len(fields) > 0 {
		buf = appendTab(buf)
		buf = appendFields(buf, fields)
	}
	buf = appendLF(buf)
	fmt.Fprintf(level.Writer, color.New(level.Color).Sprint(string(buf)))
	if level == LevelFatal {
		os.Exit(1)
	}
}

func appendFields(dst []byte, fields map[string]interface{}) []byte {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		dst = appendField(dst, key, fields[key])
		dst = appendTab(dst)
	}
	return dst
}

func appendTab(dst []byte) []byte {
	return append(dst, '	')
}

func appendField(dst []byte, key string, val interface{}) []byte {
	dst = appendString(dst, fmt.Sprintf("%s=%v", key, val))
	return dst
}

func appendString(dst []byte, s string) []byte {
	dst = append(dst, []byte(s)...)
	return dst
}

func appendLF(dst []byte) []byte {
	return append(dst, '\n')
}
