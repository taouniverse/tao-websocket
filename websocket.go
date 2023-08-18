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
	"github.com/taouniverse/tao"
	"net/http"

	// Load the required dependencies.
	// An error occurs when there was no package in the root directory.
	ws "github.com/gorilla/websocket"
)

/**
import _ "github.com/taouniverse/tao-websocket"
*/

// W config of websocket
var W = new(Config)

func init() {
	err := tao.Register(ConfigKey, W, setup)
	if err != nil {
		panic(err.Error())
	}
}

// httpUpgrade upgrade http to websocket
var httpUpgrade ws.Upgrader

// TODO setup unit with the global config 'W'
// execute when init tao universe
func setup() error {
	httpUpgrade = ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return nil
}
