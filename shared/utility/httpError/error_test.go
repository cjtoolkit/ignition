// +build debug

package httpError

import (
	"errors"
	"testing"
)

func TestCheckParamErr(t *testing.T) {
	t.Run("Has Error", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("Should not be nil")
			}
		}()

		CheckParamErr(errors.New("I am error"))
	})

	t.Run("Has no error", func(t *testing.T) {
		CheckParamErr(nil)
	})
}
