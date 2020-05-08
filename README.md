# Ignition

Web Application Initializer

## To Install

```sh
$ go get github.com/cjtoolkit/ignition/ignite
```

## To Initialise an application

First run, second param is the name of directory and the third is the proposed repo location.

```sh
$ ignite base baseDir github.com/org/base
```

It's should create a new directory, update the constant variables in `constant` directory
inside the new directory; then commit and push into the new repository.

Then run, the second param is the name of app directory,
the third is the proposed repo location for the app and the fourth is the base repo location.

```sh
$ ignite app appDir github.com/org/app github.com/org/base
```

Then run the following

```
$ cd appDir
$ go get github.com/org/base
$ go mod tidy
```

## Other required tools
* [gnode](https://github.com/cjtoolkit/gnode) To manage nodejs and to generate asset
(e.g css and javascript)
* [gobox](https://github.com/cjtoolkit/gobox) To manage tools for `go generate`.