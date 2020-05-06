// +debug debug

package coalesce

import "testing"

func TestCoalesceString(t *testing.T) {
	t.Run("Runs both function, because first is empty", func(t *testing.T) {
		var first, second bool

		str := Strings(
			func() string {
				first = true
				return ""
			},
			func() string {
				second = true
				return "Smith"
			},
		)

		if str != "Smith" {
			t.Error("Not 'Smith'")
		}
		if !first {
			t.Error("First was not called")
		}
		if !second {
			t.Error("Second was not called")
		}
	})

	t.Run("Runs one function, because first is not empty", func(t *testing.T) {
		var first, second bool

		str := Strings(
			func() string {
				first = true
				return "John"
			},
			func() string {
				second = true
				return "Smith"
			},
		)

		if str != "John" {
			t.Error("Not 'John'")
		}
		if !first {
			t.Error("First was not called")
		}
		if second {
			t.Error("Second was called")
		}
	})
}
