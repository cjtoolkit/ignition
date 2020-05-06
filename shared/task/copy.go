package task

import (
	"github.com/cjtoolkit/ignition/shared/task/internal"
	"github.com/cjtoolkit/taskforce"
)

func CopyFolder(tf *taskforce.TaskForce) func(dest, src string) {
	return func(dest, src string) {
		tf.CheckError(internal.CopyFolder(dest, src))
	}
}

func CopyFile(tf *taskforce.TaskForce) func(dest string, srcs ...string) {
	return func(dest string, srcs ...string) {
		for _, src := range srcs {
			tf.CheckError(internal.CopyFile(dest, src))
		}
	}
}
