// +build debug

package signature

import (
	"errors"
	"testing"
)

func TestCheckErrorAndForbid(t *testing.T) {
	t.Run("Has error", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("Recover should not be nil.")
			}
		}()

		checkErrorAndForbid(errors.New("E"))
	})

	t.Run("Has no error", func(t *testing.T) {
		checkErrorAndForbid(nil)
	})
}

func TestCheckBoolAndForbid(t *testing.T) {
	t.Run("Not okay", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("Recover should not be nil.")
			}
		}()

		checkBoolAndForbid(false)
	})

	t.Run("Is okay", func(t *testing.T) {
		checkBoolAndForbid(true)
	})
}
