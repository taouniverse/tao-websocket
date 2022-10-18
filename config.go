// Copyright 2022 huija
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package websocket

import (
	"context"
	"github.com/taouniverse/tao"
)

// ConfigKey for this repo
const ConfigKey = "websocket"

// Config implements tao.Config
// TODO declare the configuration you want & define some default values
type Config struct {
	RunAfters []string `json:"run_after,omitempty"`
}

var defaultWebsocket = &Config{
	RunAfters: []string{},
}

// Name of Config
func (w *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (w *Config) ValidSelf() {
	if w.RunAfters == nil {
		w.RunAfters = defaultWebsocket.RunAfters
	}
}

// ToTask transform itself to Task
func (w *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			// non-block check
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			// TODO JOB code run after RunAfters, you can just do nothing here
			return param, nil
		})
}

// RunAfter defines pre task names
func (w *Config) RunAfter() []string {
	return w.RunAfters
}
