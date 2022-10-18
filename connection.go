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
	"errors"
	ws "github.com/gorilla/websocket"
	"github.com/taouniverse/tao"
	"net/http"
)

// Connection of websocket
type Connection struct {
	*ws.Conn

	requestBuff  chan []byte
	responseBuff chan []byte
	close        chan bool
}

// New constructor of websocket connection
func New(options []func(*Connection) error) (conn *Connection, err error) {
	for _, option := range options {
		err = option(conn)
		if err != nil {
			break
		}
	}
	if conn.Conn == nil {

	}
	return
}

// Response to websocket client
func (conn *Connection) Response() error {
	for {
		select {
		case <-conn.close:
			return errors.New("websocket: connection has been closed")
		case response := <-conn.responseBuff:
			err := conn.WriteMessage(ws.TextMessage, response)
			if err != nil {
				return tao.NewErrorWrapped("websocket: fail to write response payload", err)
			}
		}
	}
}

// Request from websocket client
func (conn *Connection) Request() error {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return tao.NewErrorWrapped("websocket: fail to read request payload", err)
		}
		select {
		case <-conn.close:
			return errors.New("websocket: connection has been closed")
		case conn.requestBuff <- p:
		}
	}
}

// Write response payload
func (conn *Connection) Write(p []byte) error {
	select {
	case <-conn.close:
		return errors.New("websocket: connection has been closed")
	case conn.responseBuff <- p:
		return nil
	}
}

// Read request payload
func (conn *Connection) Read() ([]byte, error) {
	select {
	case <-conn.close:
		return nil, errors.New("websocket: connection has been closed")
	case p := <-conn.requestBuff:
		return p, nil
	}
}

// Close websocket connection
func (conn *Connection) Close() error {
	select {
	case conn.close <- true:
		close(conn.close)
		return conn.Conn.Close()
	default:
		return nil
	}
}

/**
Optional Function for Websocket Connection
*/

// Bind responseWriter & request of handler
func Bind(w http.ResponseWriter, r *http.Request) func(conn *Connection) error {
	return func(conn *Connection) (err error) {
		conn.Conn, err = HttpUpgrade.Upgrade(w, r, nil)
		return
	}
}

// Standard websocket connection
func Standard() func(conn *Connection) error {
	return func(conn *Connection) error {
		conn.requestBuff = make(chan []byte, 1000)
		conn.responseBuff = make(chan []byte, 1000)
		conn.close = make(chan bool)
		go func() {
			err := conn.Request()
			if err != nil {
				tao.Error(err)
			}
		}()
		go func() {
			err := conn.Response()
			if err != nil {
				tao.Error(err)
			}
		}()
		return nil
	}
}
