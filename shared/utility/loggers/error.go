//go:generate gobox tools/easymock

package loggers

import (
	"log"
	"os"
	"strings"

	"github.com/cjtoolkit/ctx/v2"
)

type ErrorService interface {
	CheckErrorAndPanic(err error)
	CheckErrorAndLog(err error)
	GetLogger() Logger
}

func GetBlankErrorService(context ctx.Context) ErrorService {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return ErrorService(initErrorService(logOutputRegistryBlank{})), nil
	}).(ErrorService)
}

func GetErrorService(context ctx.Context) ErrorService {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return ErrorService(initErrorService(newLogOutputRegistry())), nil
	}).(ErrorService)
}

func initErrorService(registry LogOutputRegistry) *errorService {
	return &errorService{
		log: customLog{
			Logger:            log.New(os.Stderr, "INFO: ", log.Lshortfile),
			logOutputRegistry: registry,
			callDepth:         3,
		},
	}
}

type errorService struct {
	log Logger
}

func (e *errorService) CheckErrorAndPanic(err error) {
	if nil != err {
		e.log.Panic(err)
	}
}

func (e *errorService) CheckErrorAndLog(err error) {
	if nil != err {
		e.log.Print(err)
	}
}

func (e *errorService) GetLogger() Logger { return e.log.Clone(2) }

type ErrorCollector []error

func (e ErrorCollector) Error() string {
	str := []string{}
	for _, err := range e {
		str = append(str, err.Error())
	}
	return strings.Join(str, "\n")
}

func (e ErrorCollector) FilterError() error {
	errorCollection := ErrorCollector{}
	for _, err := range e {
		if nil != err {
			errorCollection = append(errorCollection, err)
		}
	}

	if 0 == len(errorCollection) {
		return nil
	}

	return errorCollection
}

type RecoverCollector []interface{}

func (e RecoverCollector) FilterError() error {
	errorCollection := ErrorCollector{}
	for _, err := range e {
		if err, ok := err.(error); ok {
			errorCollection = append(errorCollection, err)
		}
	}

	if 0 == len(errorCollection) {
		return nil
	}

	return errorCollection
}
