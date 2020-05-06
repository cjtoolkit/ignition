//go:generate gobox tools/gmock

package postgres

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/cjtoolkit/ctx"
)

type Param map[string]interface{}

type ParamTemplate interface {
	ToSlice(param Param) []interface{}
}

type paramTemplate []string

func (p paramTemplate) ToSlice(param Param) []interface{} {
	s := []interface{}{}
	for _, name := range p {
		s = append(s, param[name])
	}
	return s
}

type ParamBuilder interface {
	BuildParamTemplate(query string) (string, ParamTemplate)
	BuildParamTemplateAndPrepare(dbConn *sql.DB, query string) (*sql.Stmt, ParamTemplate, error)
}

func GetParamBuilder(context ctx.BackgroundContext) ParamBuilder {
	type paramBuilderContext struct{}
	return context.Persist(paramBuilderContext{}, func() (interface{}, error) {
		return ParamBuilder(paramBuilder{
			paramRegExp: regexp.MustCompile(`:([a-zA-Z_]+)`),
			searchWith:  string([]byte{':', ':'}),
			replaceWith: string([]byte{0, ':', 0, ':', 0}),
		}), nil
	}).(ParamBuilder)
}

type paramBuilder struct {
	paramRegExp *regexp.Regexp
	searchWith  string
	replaceWith string
}

func (p paramBuilder) BuildParamTemplate(query string) (string, ParamTemplate) {
	paramTemplate := paramTemplate{}
	got := map[string]bool{}
	count := 0

	query = strings.Replace(query, p.searchWith, p.replaceWith, -1)
	matches := p.paramRegExp.FindAllString(query, -1)
	for _, match := range matches {
		if !got[match] {
			got[match] = true
			count++
			paramTemplate = append(paramTemplate, match[1:])
			query = strings.Replace(query, match, fmt.Sprintf("$%d", count), -1)
		}
	}
	query = strings.Replace(query, p.replaceWith, p.searchWith, -1)

	return query, paramTemplate
}

func (p paramBuilder) BuildParamTemplateAndPrepare(dbConn *sql.DB, query string) (*sql.Stmt, ParamTemplate, error) {
	query, template := p.BuildParamTemplate(query)
	stmt, err := dbConn.Prepare(query)
	return stmt, template, err
}

type StmtParam struct {
	Stmt  *sql.Stmt
	Param ParamTemplate
}
