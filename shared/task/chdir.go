package task

import (
	"os"

	"github.com/cjtoolkit/taskforce"
)

func Chdir(tf *taskforce.TaskForce) func(dir string) func() {
	return func(dir string) func() {
		currentWd, err := os.Getwd()
		tf.CheckError(err)
		tf.CheckError(os.Chdir(dir))

		return func() {
			tf.CheckError(os.Chdir(currentWd))
		}
	}
}
