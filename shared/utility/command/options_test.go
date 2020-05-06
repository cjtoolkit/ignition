// +build debug

package command

import "testing"

func TestBuildOptions(t *testing.T) {
	args := []string{"abc", "--o", "--option1", "def", "--option2", "ghi", "--option1", "def"}

	op := BuildOptions(args)

	if op.Values[""][0] != "abc" && op.Values[""][1] != "--o" {
		t.Errorf("Not 'abc' and '--o'")
	}

	if op.Values["option1"][0] != "def" && op.Values["option1"][1] != "def" {
		t.Errorf("Not 'def'")
	}

	if op.Values["option2"][0] != "ghi" {
		t.Errorf("Not 'ghi'")
	}

	{
		data := ""
		op.ExecOption("notset", func(_ []string) {
			data = "set"
		})

		if data != "" {
			t.Error("Should not be set")
		}
	}

	{
		data := ""
		op.ExecOption("option1", func(_ []string) {
			data = "set"
		})

		if data != "set" {
			t.Error("Should be set")
		}
	}
}
