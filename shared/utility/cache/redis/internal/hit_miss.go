//go:generate gobox tools/gmock

package internal

type HitMiss interface {
	Miss() (data interface{}, b []byte, err error)
	Hit(b []byte) (data interface{}, err error)
}
