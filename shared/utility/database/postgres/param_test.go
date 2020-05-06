// +build debug

package postgres

import (
	"testing"

	"github.com/cjtoolkit/ctx"
)

func TestBuildParamTemplateFromSql(t *testing.T) {
	query := ":apple ::boolean :pair ::boolean :apple"

	query, param := GetParamBuilder(ctx.NewBackgroundContext()).BuildParamTemplate(query)

	if query != "$1 ::boolean $2 ::boolean $1" {
		t.Error("Not '$1 ::boolean $2 ::boolean $1'")
	}

	if 2 != len(param.(paramTemplate)) {
		t.Error("Not 2")
	}
	if param.(paramTemplate)[0] != "apple" {
		t.Error("Not 'apple'")
	}
	if param.(paramTemplate)[1] != "pair" {
		t.Error("Not 'pair'")
	}
}
