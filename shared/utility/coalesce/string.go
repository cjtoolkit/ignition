package coalesce

import "strings"

func Strings(fns ...func() string) string {
	for _, fn := range fns {
		if str := strings.TrimSpace(fn()); "" != str {
			return str
		}
	}
	return ""
}
