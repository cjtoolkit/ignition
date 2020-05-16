//go:generate gobox tools/easymock

package loggers

import (
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
)

type Logger interface {
	OutputRegistry() LogOutputRegistry
	Clone(callDepth int) Logger
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type customLog struct {
	*log.Logger
	logOutputRegistry LogOutputRegistry
	callDepth         int
}

func (l customLog) OutputRegistry() LogOutputRegistry { return l.logOutputRegistry }

func (l customLog) Output(calldepth int, kind, s string) error {
	stack := debug.Stack()
	_, file, line, ok := runtime.Caller(calldepth)
	file, line = checkCallStatus(file, line, ok)
	for _, outputter := range l.logOutputRegistry.Outputters() {
		outputter.Output(line, file, kind, s, stack)
	}
	return l.Logger.Output(calldepth+1, s)
}

func (l customLog) Clone(callDepth int) Logger {
	return customLog{
		Logger:            l.Logger,
		logOutputRegistry: l.logOutputRegistry,
		callDepth:         callDepth,
	}
}

func (l customLog) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(l.callDepth, "Panic", s)
	panic(s)
}

func (l customLog) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(l.callDepth, "Panic", s)
	panic(s)
}

func (l customLog) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(l.callDepth, "Panic", s)
	panic(s)
}

func (l customLog) Print(v ...interface{}) { l.Output(l.callDepth, "Print", fmt.Sprint(v...)) }

func (l customLog) Printf(format string, v ...interface{}) {
	l.Output(l.callDepth, "Print", fmt.Sprintf(format, v...))
}

func (l customLog) Println(v ...interface{}) { l.Output(l.callDepth, "Print", fmt.Sprintln(v...)) }

type LoggerOutputer interface {
	Output(line int, file, kind, s string, stack []byte) error
}

type LogOutputRegistry interface {
	Outputters() []LoggerOutputer
	Register(outputers ...LoggerOutputer) LogOutputRegistry
	Lock()
}

type logOutputRegistryBlank struct{}

func (_ logOutputRegistryBlank) Outputters() []LoggerOutputer                           { return nil }
func (l logOutputRegistryBlank) Register(outputers ...LoggerOutputer) LogOutputRegistry { return l }
func (_ logOutputRegistryBlank) Lock()                                                  {}

type logOutputRegistry struct {
	outputters []LoggerOutputer
	register   func(outputer ...LoggerOutputer)
}

func newLogOutputRegistry() *logOutputRegistry {
	l := &logOutputRegistry{}
	l.register = func(outputer ...LoggerOutputer) {
		l.outputters = append(l.outputters, outputer...)
	}

	return l
}

func (l *logOutputRegistry) Outputters() []LoggerOutputer { return l.outputters }

func (l *logOutputRegistry) Register(outputers ...LoggerOutputer) LogOutputRegistry {
	l.register(outputers...)
	return l
}

func (l *logOutputRegistry) Lock() {
	l.register = func(outputer ...LoggerOutputer) { panic("Registry is Locked") }
}

func checkCallStatus(file string, line int, ok bool) (string, int) {
	if !ok {
		file = "???"
		line = 0
	}
	return file, line
}
