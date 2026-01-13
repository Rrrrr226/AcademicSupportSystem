// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sse

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/flamego/flamego"
)

// Options contains options for the sse.Bind middleware.
type Options struct {
	// PingInterval is the time internal to wait between sending pings to the
	// client. Default is 10 seconds.
	PingInterval time.Duration
}

type connection struct {
	Options

	// sender is the channel used for sending out data to the client. This channel
	// gets mapped for the next handler to use with the right type and is
	// asynchronous unless the SendChannelBuffer is set to 0.
	sender reflect.Value
}

// Bind returns a middleware handler that uses the given bound object as the
// date type for sending events.
func Bind(obj interface{}, opts ...Options) flamego.Handler {
	return func(c flamego.Context, log *log.Logger) {
		c.ResponseWriter().Header().Set("Content-Type", "text/event-stream")
		c.ResponseWriter().Header().Set("Cache-Control", "no-cache")
		c.ResponseWriter().Header().Set("Connection", "keep-alive")
		c.ResponseWriter().Header().Set("X-Accel-Buffering", "no")

		sse := &connection{
			Options: newOptions(opts),
			// Create a chan of the given type as a reflect.Value.
			sender: reflect.MakeChan(reflect.ChanOf(reflect.BothDir, reflect.PtrTo(reflect.TypeOf(obj))), 0),
		}
		c.Set(reflect.ChanOf(reflect.SendDir, sse.sender.Type().Elem()), sse.sender)

		go sse.handle(log, c)
	}
}

// newOptions creates new default options and assigns any given options.
func newOptions(opts []Options) Options {
	if len(opts) == 0 {
		return Options{
			PingInterval: 10 * time.Second,
		}
	}
	return opts[0]
}

func (c *connection) handle(log *log.Logger, ctx flamego.Context) {
	// 捕获可能的 panic，防止连接断开后写入导致崩溃
	defer func() {
		if r := recover(); r != nil {
			log.Printf("sse: recovered from panic: %v", r)
		}
	}()

	w := ctx.ResponseWriter()
	ticker := time.NewTicker(c.PingInterval)
	defer ticker.Stop()

	// closed 用于标记连接是否已关闭
	closed := false

	write := func(msg string) bool {
		if closed {
			return false
		}
		// 检查 context 是否已取消
		select {
		case <-ctx.Request().Context().Done():
			closed = true
			return false
		default:
		}
		_, err := w.Write([]byte(msg))
		if err != nil {
			log.Printf("sse: failed to write message: %v", err)
			closed = true
			return false
		}
		return true
	}

	flush := func() bool {
		if closed {
			return false
		}
		// 检查 context 是否已取消
		select {
		case <-ctx.Request().Context().Done():
			closed = true
			return false
		default:
		}
		w.Flush()
		return true
	}

	if !write(": ping\n\n") {
		return
	}
	if !write("events: stream opened\n\n") {
		return
	}
	if !flush() {
		return
	}

	const (
		senderSend = iota
		tickerTick
		timeout
		ctxClosed
	)
	cases := make([]reflect.SelectCase, 4)
	cases[senderSend] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: c.sender, Send: reflect.ValueOf(nil)}
	cases[tickerTick] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ticker.C), Send: reflect.ValueOf(nil)}
	cases[timeout] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(time.Hour)), Send: reflect.ValueOf(nil)}
	cases[ctxClosed] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Request().Context().Done()), Send: reflect.ValueOf(nil)}

loop:
	for {
		if closed {
			break loop
		}

		chosen, message, ok := reflect.Select(cases)
		switch chosen {
		case senderSend:
			if !ok {
				// Sender channel has been closed.
				return
			}

			if !write("data: ") {
				return
			}
			evt, err := json.Marshal(message.Interface())
			if err != nil {
				log.Printf("sse: failed to marshal message: %v", err)
				continue
			}
			if !write(string(evt)) {
				return
			}
			if !write("\n\n") {
				return
			}
			if !flush() {
				return
			}

		case tickerTick:
			if !write(": ping\n\n") {
				return
			}
			if !flush() {
				return
			}

		case timeout:
			write("events: stream timeout\n\n")
			flush()
			break loop

		case ctxClosed:
			// 客户端已断开连接
			return
		}
	}

	write("events: error\ndata: eof\n\n")
	flush()
	write("events: stream closed")
	flush()
}
