//: ----------------------------------------------------------------------------
//: Copyright (C) 2019 Helmut Wahrmann.
//:
//: file:    options.go
//: details: Handle the Configuration
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/google/logger"
	"gopkg.in/yaml.v2"
)

var (
	version = "1.0.0"
)

const (
	productname = "RSA NetWitness Monitor"
	configFile  = "/etc/rsa-nw-monitor/rsa-nw-monitor.conf"
)

// Options represents options
type Options struct {
	// global options
	Verbose        bool   `yaml:"verbose"`
	PIDFile        string `yaml:"pid-file"`
	Logger         *logger.Logger
	version        bool
	MonitorPort    int    `yaml:"monitorport"`
	EndPointServer string `yaml:"endpointserver"`
	APIUser        string `yaml:"user"`
	APIUserPwd     string `yaml:"password"`
}

func init() {
	if version == "" {
		version = "1.0.0"
	}
}

// NewOptions constructs new options
func NewOptions() *Options {
	options := Options{}
	options.Verbose = false
	options.PIDFile = "/var/run/rsa-nw-monitor.pid"
	options.MonitorPort = 8080
	options.Logger = logger.Init("", options.Verbose, true, ioutil.Discard)
	logger.SetFlags(0)
	return &options
}

// GetOptions gets options through cmd and file
func GetOptions() *Options {
	options := NewOptions()

	options.productFlagSet()
	options.productVersion()

	if ok := options.productIsRunning(); ok {
		options.Logger.Fatalf("%s is already running!", productname)
	}

	options.pidWrite()

	options.Logger.Infof("Welcome to %s v.%s GPL v3", productname, version)
	options.Logger.Info("Copyright (C) 2019 Helmut Wahrmann.")

	return options
}

func (opts Options) pidWrite() {
	f, err := os.OpenFile(opts.PIDFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		opts.Logger.Info(err)
		return
	}

	_, err = fmt.Fprintf(f, "%d", os.Getpid())
	if err != nil {
		opts.Logger.Info(err)
	}
}

func (opts Options) productIsRunning() bool {
	b, err := ioutil.ReadFile(opts.PIDFile)
	if err != nil {
		return false
	}

	cmd := exec.Command("kill", "-0", string(b))
	_, err = cmd.Output()
	if err != nil {
		return false
	}

	return true
}

func (opts Options) productVersion() {
	if opts.version {
		fmt.Printf("%s version: %s\n", productname, version)
		os.Exit(0)
	}
}

func (opts *Options) productFlagSet() {

	var config string
	flag.StringVar(&config, "config", configFile, "path to config file")

	loadConfig(opts)

	// global options
	flag.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "enable/disable verbose logging")
	flag.BoolVar(&opts.version, "version", opts.version, "show version")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
    Example:
	rsa-nw-monitor -config %s"
	`, configFile)
	}

	flag.Parse()
}

// Load the configuration from the config file
func loadConfig(opts *Options) {
	var file = configFile

	// Check, if we received a config file via arguments
	for i, flag := range os.Args {
		if flag == "-config" {
			file = os.Args[i+1]
			break
		}
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		opts.Logger.Info(err)
		return
	}
	err = yaml.Unmarshal(b, opts)
	if err != nil {
		opts.Logger.Info(err)
	}
}
