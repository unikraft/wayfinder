package logs

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2020, Lancaster University.  All rights reserved.
//               2021, Unikraft UG.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

import (
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"

	l2 "github.com/erda-project/erda-infra/base/logs"
)

// LogLevel is an enum-like type that we can use to designate the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
	SUCCESS
)

// Logger is a base struct that could eventually maintain connections to
// something like bugsnag or logging tools.
type Logger struct {
	logLevel LogLevel
	Prefix   string
}

var (
	globalLogger *Logger
	logFile      *os.File
)

func init() {
	globalLogger = &Logger{
		logLevel: INFO,
	}

	// TODO - make this configurable - Who owns the logs?
	err := os.MkdirAll("/var/lib/wayfinder/logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("/var/lib/wayfinder/logs/wayfinder.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logFile = f
}

func New() Logger {
	logger := Logger{
		logLevel: INFO,
	}

	return logger
}

// log is a private function that manages the internal logic about what and how
// to log data depending on the log level.
func (l Logger) log(level LogLevel, format string, messages ...interface{}) {
	var logType string
	var logColor termenv.ANSIColor
	switch level {
	case DEBUG:
		logType = "DEBU"
		logColor = termenv.ANSICyan
	case WARNING:
		logType = "WARN"
		logColor = termenv.ANSIYellow
	case ERROR:
		logType = "ERRO"
		logColor = termenv.ANSIRed
	case FATAL:
		logType = "FATA"
		logColor = termenv.ANSIRed
	case SUCCESS:
		logType = " :) "
		logColor = termenv.ANSIGreen
	default:
		logType = "INFO"
		logColor = termenv.ANSIBlue
	}

	if level < l.logLevel {
		return
	}

	// Add some colours!
	out := termenv.String(logType)
	out = out.Foreground(logColor)

	if len(l.Prefix) > 0 {
		fmt.Printf("[%s][%s] %s\n", out, l.Prefix, fmt.Sprintf(format, messages...))
		fmt.Fprintf(logFile, "[%s][%s] %s\n", out, l.Prefix, fmt.Sprintf(format, messages...))
	} else {
		fmt.Printf("[%s] %s\n", out, fmt.Sprintf(format, messages...))
		fmt.Fprintf(logFile, "[%s] %s\n", out, fmt.Sprintf(format, messages...))
	}
}

func SetLevel(level LogLevel) {
	globalLogger.logLevel = level
}

func (l Logger) SetLevel(lvl string) error {
	return nil
}

func (l Logger) Sub(name string) l2.Logger {
	// globalLogger.logLevel = lvl
	l.Prefix = name
	return l
}

func GetLevel() LogLevel {
	return globalLogger.logLevel
}

func Debug(args ...interface{}) {
	globalLogger.log(DEBUG, "%s", args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.log(DEBUG, format, args...)
}

func (l Logger) Debug(args ...interface{}) {
	l.log(DEBUG, "%s", args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func Info(args ...interface{}) {
	globalLogger.log(INFO, "%s", args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.log(INFO, format, args...)
}

func (l Logger) Info(args ...interface{}) {
	l.log(INFO, "%s", args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func Warn(args ...interface{}) {
	globalLogger.log(WARNING, "%s", args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.log(WARNING, format, args...)
}

func (l Logger) Warn(args ...interface{}) {
	l.log(WARNING, "%s", args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

func Error(args ...interface{}) {
	globalLogger.log(ERROR, "%s", args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.log(ERROR, format, args...)
}

func (l Logger) Error(args ...interface{}) {
	l.log(ERROR, "%s", args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func Fatal(args ...interface{}) {
	globalLogger.log(FATAL, "%s", args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.log(FATAL, format, args...)
}

func (l Logger) Fatal(args ...interface{}) {
	l.log(FATAL, "%s", args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

func Panic(args ...interface{}) {
	globalLogger.log(FATAL, "%s", args...)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.log(FATAL, format, args...)
}

func (l Logger) Panic(args ...interface{}) {
	l.log(FATAL, "%s", args...)
}

func (l Logger) Panicf(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

func Success(args ...interface{}) {
	globalLogger.log(SUCCESS, "%s", args...)
}

func Successf(format string, args ...interface{}) {
	globalLogger.log(SUCCESS, format, args...)
}

func (l Logger) Success(args ...interface{}) {
	l.log(SUCCESS, "%s", args...)
}

func (l Logger) Successf(format string, args ...interface{}) {
	l.log(SUCCESS, format, args...)
}

// Write implements io.Writer
func (l Logger) Write(b []byte) (n int, err error) {
	if len(string(b)) > 0 {
		args := strings.Split(string(b), "\n")
		for _, arg := range args {
			if len(arg) > 0 {
				l.log(INFO, "%s", arg)
			}
		}
	}
	return len(b), nil
}
