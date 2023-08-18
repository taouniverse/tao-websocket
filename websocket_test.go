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
	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
	"net/http"
	"testing"
)

func TestTao(t *testing.T) {
	err := tao.DevelopMode()
	assert.Nil(t, err)

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":9000", nil)

	err = tao.Run(nil, nil)
	assert.Nil(t, err)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := New(w, r)
	if err != nil {
		tao.Error(err)
		return
	}
	defer conn.Close()

	for {
		data, err := conn.Read()
		if err != nil {
			tao.Error(err)
			return
		}
		err = conn.Write([]byte(reverseString(string(data))))
		if err != nil {
			tao.Error(err)
			return
		}
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
