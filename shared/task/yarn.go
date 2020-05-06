package task

import "github.com/cjtoolkit/taskforce"

func YarnRun(tf *taskforce.TaskForce) func(name string, args ...string) {
	return func(name string, args ...string) {
		args = append([]string{"yarn", "run", name}, args...)
		tf.ExecCmd("gnode", args...)
	}
}
