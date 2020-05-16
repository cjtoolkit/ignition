// +build debug

package embedder

import (
	"fmt"
	"testing"
)

func TestCheckErr(t *testing.T) {
	t.Run("No Error", func(t *testing.T) {
		defer func() {
			if recover() != nil {
				t.Fail()
			}
		}()

		errCheck(nil)
	})

	t.Run("With Error", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fail()
			}
		}()

		errCheck(fmt.Errorf("I am error"))
	})
}
