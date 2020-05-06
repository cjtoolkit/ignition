// +build debug

package client

import "testing"

func TestCleanHostAddress(t *testing.T) {
	t.Run("No colon", func(t *testing.T) {
		if cleanHostAddress("hello") != "hello" {
			t.Fail()
		}
	})

	t.Run("With colon on port 80", func(t *testing.T) {
		if cleanHostAddress("hello:80") != "hello" {
			t.Fail()
		}
	})

	t.Run("With colon on another port", func(t *testing.T) {
		if cleanHostAddress("hello:8080") != "hello:8080" {
			t.Fail()
		}
	})
}
