package postgres

import (
	"database/sql"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/embedder"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	_ "github.com/lib/pq"
)

func GetMainSqlDatabase(context ctx.Context) *sql.DB {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		sqlDsn := configuration.GetConfig(context).Database.MainSqlDsn

		db, err := sql.Open("postgres", sqlDsn)
		return db, err
	}).(*sql.DB)
}

func GetMainSqlPrepareKit(context ctx.Context) *PrepareKit {
	return &PrepareKit{
		DB:           GetMainSqlDatabase(context),
		Builder:      GetParamBuilder(context),
		ErrorService: loggers.GetErrorService(context),
	}
}

type PrepareKit struct {
	DB           *sql.DB
	Builder      ParamBuilder
	ErrorService loggers.ErrorService
}

func (k *PrepareKit) Prepare(query string) *sql.Stmt {
	stmt, err := k.DB.Prepare(query)
	k.ErrorService.CheckErrorAndPanic(err)
	return stmt
}

func (k *PrepareKit) DecodePrepare(query string) *sql.Stmt {
	return k.Prepare(embedder.DecodeValueStr(query))
}

func (k *PrepareKit) PrepareParam(query string) StmtParam {
	stmt, param, err := k.Builder.BuildParamTemplateAndPrepare(k.DB, query)
	k.ErrorService.CheckErrorAndPanic(err)
	return StmtParam{
		Stmt:  stmt,
		Param: param,
	}
}

func (k *PrepareKit) DecodePrepareParam(query string) StmtParam {
	return k.PrepareParam(embedder.DecodeValueStr(query))
}
