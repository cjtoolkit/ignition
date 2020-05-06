package command

import (
	"unicode/utf8"
)

type Options struct {
	Has    map[string]bool
	Values map[string][]string
}

func BuildOptions(args []string) Options {
	op := &Options{
		Has: map[string]bool{
			"": true,
		},
		Values: map[string][]string{
			"": {},
		},
	}

	option := ""
	for _, value := range args {
		if utf8.RuneCountInString(value) > 4 && value[:2] == "--" {
			option = value[2:]
			op.Has[option] = true
			op.Values[option] = append(op.Values[option], []string{}...)
		} else {
			op.Values[option] = append(op.Values[option], value)
		}
	}

	return *op
}

func (op Options) ExecOption(name string, fn func([]string)) {
	if nil == fn || !op.Has[name] {
		return
	}

	fn(op.Values[name])
}
