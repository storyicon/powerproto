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

package progressbar

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

// ProgressBar implements a customizable progress bar
type ProgressBar interface {
	// Incr is used to increase progress
	Incr()
	// Wait is used to wait for the rendering of the progress bar to complete
	Wait()
	// SetPrefix is used to set the prefix of progress bar
	SetPrefix(format string, args ...interface{})
	// SetSuffix is used to set the suffix of progress bar
	SetSuffix(format string, args ...interface{})
}

type progressBar struct {
	container *mpb.Progress
	bar       *mpb.Bar
	prefix    string
	suffix    string
}

// SetPrefix is used to set the prefix of progress bar
func (s *progressBar) SetPrefix(format string, args ...interface{}) {
	s.prefix = fmt.Sprintf(format, args...)
}

// SetSuffix is used to set the suffix of progress bar
func (s *progressBar) SetSuffix(format string, args ...interface{}) {
	s.suffix = fmt.Sprintf(format, args...)
}

func newEmbedProgressBar(container *mpb.Progress, bar *mpb.Bar) *progressBar {
	return &progressBar{
		container: container,
		bar:       bar,
	}
}

func (s *progressBar) Incr() {
	s.bar.Increment()
}

func (s *progressBar) Wait() {
	s.container.Wait()
}

func getSpinner() []string {
	activeState := "[ " + color.GreenString("●") + " ] "
	defaultState := "[   ] "
	return []string{
		activeState,
		activeState,
		activeState,
		defaultState,
		defaultState,
		defaultState,
	}
}

// GetProgressBar is used to get progress bar
func GetProgressBar(count int) ProgressBar {
	var progressBar *progressBar
	container := mpb.New()
	bar := container.Add(int64(count),
		mpb.NewBarFiller(mpb.BarStyle().Lbound("[").
			Filler(color.GreenString("=")).
			Tip(color.GreenString(">")).Padding(" ").Rbound("]")),
		mpb.PrependDecorators(
			func() decor.Decorator {
				frames := getSpinner()
				var count uint
				return decor.Any(func(statistics decor.Statistics) string {
					if statistics.Completed {
						return frames[0]
					}
					frame := frames[count%uint(len(frames))]
					count++
					return frame
				})
			}(),
			decor.Any(func(statistics decor.Statistics) string {
				if progressBar != nil {
					return progressBar.prefix
				}
				return ""
			}),
		),
		mpb.AppendDecorators(
			decor.NewPercentage("%d  "),
			decor.Any(func(statistics decor.Statistics) string {
				if progressBar != nil {
					return fmt.Sprintf("(%d/%d) %s", statistics.Current, count, progressBar.suffix)
				}
				return ""
			}),
		),
		mpb.BarWidth(15),
	)
	progressBar = newEmbedProgressBar(container, bar)
	return progressBar
}
