// Copyright © 2019 D.Snezhkov <dsnezhkov@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"deepsea/cmd"
	"deepsea/global"
	jlog "github.com/spf13/jwalterweatherman"
	"os"
)

func main() {

	var err error
	var logSink *os.File
	var logFile = "deepsea.log"

	jlog.SetLogThreshold(jlog.LevelDebug)
	jlog.SetStdoutThreshold(jlog.LevelDebug)

	if !global.FileExists(logFile) {
		jlog.DEBUG.Printf("Creating log %s", logFile)
		logSink, err = os.Create(logFile)
		if err != nil {
			jlog.ERROR.Printf("Log not available %v", err)
		}
		jlog.DEBUG.Printf("Created log %s", logFile)
	} else {
		logSink, err = os.Open(logFile)
		if err != nil {
			jlog.ERROR.Printf("Log not available %v", err)
		}
	}
	defer func() {
		_ = logSink.Close()
	}()

	jlog.SetLogOutput(logSink)
	cmd.Execute()
}
