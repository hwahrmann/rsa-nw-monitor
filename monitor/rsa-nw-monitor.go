//: ----------------------------------------------------------------------------
//: Copyright (C) 2019 Helmut Wahrmann.
//:
//: file:    rsa-nw-monitor.go
//: details: Main Program
//: author:  Helmut Wahrmann
//: date:    03/09/2019
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------

// Package main is the rsa-nw-monitor binary
package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	opts *Options
)

func main() {
	var (
		signalCh = make(chan os.Signal, 1)
	)

	// Retrieve the Optons
	opts = GetOptions()

	// Notify on SIGINT and SIGTERM
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	apiserver := NewAPIServer()

	go func() {
		err := apiserver.run()
		if err != nil {
			close(signalCh)
		}
	}()

	<-signalCh

	apiserver.shutdown()
	opts.Logger.Info("Stopping RSA NetWitness Monitor")
}
