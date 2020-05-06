package password

import (
	"encoding/base64"

	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"

	"github.com/cjtoolkit/ctx"
	"golang.org/x/crypto/bcrypt"
)

type Password interface {
	SaltPassword(password string) string
	CheckPassword(password, hash string) bool
}

type password struct {
	salt         string
	errorService loggers.ErrorService
}

func GetPassword(context ctx.BackgroundContext) Password {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return Password(password{
			salt:         configuration.GetConfig(context).PasswordSalt,
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(Password)
}

func (p password) SaltPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.salt+password), 14)
	p.errorService.CheckErrorAndPanic(err)

	return base64.URLEncoding.EncodeToString(hash)
}

func (p password) CheckPassword(password, hash string) (ok bool) {
	hashBytes, err := base64.URLEncoding.DecodeString(hash)
	if nil != err {
		return
	}

	ok = nil == bcrypt.CompareHashAndPassword(hashBytes, []byte(p.salt+password))
	return
}
