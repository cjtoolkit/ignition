package task

import (
	"github.com/cjtoolkit/ignition/shared/task/internal"
	"github.com/cjtoolkit/taskforce"
)

func Zip(tf *taskforce.TaskForce) func(fileName, src string) {
	return func(fileName, src string) {
		tf.CheckError(internal.CreateZip(fileName, src))
	}
}
