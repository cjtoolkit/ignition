package configuration

import (
	"encoding/json"
	"os"
	"path/filepath"

	ctx "github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/environment"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

type Base struct {
	Database     Database
	CsrfKey      string
	CookieKey    string
	PasswordSalt string
	HmacKey      string
}

type Database struct {
	MainSqlDsn string
	MongoDial  string
	MongoDb    string
	Redis      Redis
}

type Redis struct {
	Addr string
}

func GetConfig(context ctx.Context) Base {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		config := &Base{}
		ParseConfig(context, "base.json", config)
		return *config, nil
	}).(Base)
}

func ParseConfig(context ctx.Context, fileName string, v interface{}) {
	location := environment.GetEnvironment(context).ParseConfigDirectory() + filepath.FromSlash("/"+fileName)
	errorService := loggers.GetBlankErrorService(context)

	file, err := os.Open(location)
	errorService.CheckErrorAndPanic(err)
	defer file.Close()

	err = json.NewDecoder(file).Decode(v)
	errorService.CheckErrorAndPanic(err)
}
