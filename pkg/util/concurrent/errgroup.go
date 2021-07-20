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

package concurrent

import (
	"context"
	"sync"
)

// ErrGroup is another ErrGroup implement
type ErrGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	limit  chan struct{}

	errOnce sync.Once
	err     error

	wg sync.WaitGroup
}

// NewErrGroup is used to create a new ErrGroup
func NewErrGroup(ctx context.Context, concurrency int) *ErrGroup {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{
		ctx:    ctx,
		cancel: cancel,
		limit:  make(chan struct{}, concurrency),
	}
}

// Wait is used to wait ErrGroup finish
func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

// Go is used to start a new goroutine
func (g *ErrGroup) Go(f func(ctx context.Context) error) {
	if g.err != nil {
		return
	}

	g.limit <- struct{}{}
	g.wg.Add(1)
	go func() {
		defer func() {
			<-g.limit
			g.wg.Done()
		}()
		if err := f(g.ctx); err != nil {
			g.cancel()
			g.errOnce.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}
