package main

import (
	"os"
	"path/filepath"

	"github.com/cjtoolkit/ignition/shared/task"
	"github.com/cjtoolkit/taskforce"
)

func initTask() *taskforce.TaskForce {
	var (
		tf             = taskforce.InitTaskForce()
		chdir          = task.Chdir(tf)
		yarnRun        = task.YarnRun(tf)
		devEnvironment = true
	)

	tf.Register("yarn", func() {
		tf.ExecCmd("gnode", "yarn")
		tf.ExecCmd("./createLink")
	})

	tf.Register("clean", func() {
		os.RemoveAll("live")

		os.Mkdir("live", 0755)
	})

	{
		// Sass
		const (
			dest = "live/stylesheets/styles.css"
			src  = "dev/sass/styles.scss"
		)

		tf.Register("sass", func() {
			args := []string{"--source-map", "true", "--source-map-contents", "true", "--precision", "8", "--output-style", "compressed"}
			args = append(args, src, dest)
			yarnRun("node-sass", args...)
		})
	}

	tf.Register("rollup", func() {
		env := "BUILD:production"
		if devEnvironment {
			env = "BUILD:development"
		}
		yarnRun("rollup", "-c", "-m", "--environment", env)
	})

	tf.Register("copy", func() {
		copyFolder := task.CopyFolder(tf)

		copyFolder("live/fonts", "link/fontawesome/webfonts")
	})

	tf.Register("zip", func() {
		defer chdir("live")()
		zipUtil := task.Zip(tf)
		zipUtil("../../asset.zip", ".")
	})

	tf.Register("dev", func() {
		tf.Run("yarn", "clean", "sass", "rollup", "copy")
	})

	tf.Register("prod", func() {
		devEnvironment = false
		tf.Run("yarn", "clean", "sass", "rollup", "copy", "zip")
	})

	tf.Register("quick:sass", func() {
		os.RemoveAll(filepath.FromSlash("live/stylesheets"))
		tf.Run("sass")
	})

	tf.Register("quick:js", func() {
		os.RemoveAll(filepath.FromSlash("live/javascripts"))
		tf.Run("rollup")
	})

	return tf
}

func main() {
	initTask().Run(os.Args[1:]...)
}
