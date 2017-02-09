// Copyright 2017 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package gdb provides a gdb remote agent.
package gdb

import (
	"context"
	"io"
	"net"
)

func connMonitor(c io.Closer, done <-chan struct{}) {
	<-done
	_ = c.Close()
}

// StartGDBAgent starts a background GDB agent for remote debugging.
//
// TODO: not implemented yet. the API will change (the system parameter is not final).
//
func StartGDBAgent(ctx context.Context, addr string, system interface{}) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	client := func(ctx context.Context, conn net.Conn) {
		cc, cancel := context.WithCancel(ctx)
		defer cancel()

		// client monitor -- will close conn when cc.Done() is closed
		go connMonitor(conn, cc.Done())

		// TODO: implement the agent
	}

	server := func(ctx context.Context, l net.Listener) {
		lc, cancel := context.WithCancel(ctx)
		defer cancel()

		// server monitor -- -- will close l when lc.Done() is closed
		go connMonitor(l, lc.Done())

		for {
			conn, err := l.Accept()
			if err != nil {
				// TODO: need callback to notify caller that something bad happenned
				// if callback != nil { callback(err) }
				return
			}
			// Handle the connection in a new goroutine.
			go client(lc, conn)
		}
	}
	go server(ctx, l)
	return nil
}
