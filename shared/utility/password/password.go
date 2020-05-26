package password

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
	"golang.org/x/crypto/bcrypt"
)

type Password interface {
	SaltPassword(password string) string
	CheckPassword(password, hash string) bool
}

type password struct {
	salt         []byte
	errorService loggers.ErrorService
}

func GetPassword(context ctx.Context) Password {
	type c struct{}
	return context.Persist(c{}, func() (interface{}, error) {
		return Password(password{
			salt:         convertToByte(configuration.GetConfig(context).PasswordSalt),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(Password)
}

func (p password) SaltPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword(append(p.salt, []byte(password)...), 14)
	p.errorService.CheckErrorAndPanic(err)

	return base64.URLEncoding.EncodeToString(hash)
}

func (p password) CheckPassword(password, hash string) (ok bool) {
	hashBytes, err := base64.URLEncoding.DecodeString(hash)
	if nil != err {
		return
	}

	ok = nil == bcrypt.CompareHashAndPassword(hashBytes, append(p.salt, []byte(password)...))
	return
}

func convertToByte(saltStr string) []byte {
	salt, err := hex.DecodeString(saltStr)
	if err != nil {
		salt = []byte(saltStr)
	}
	return salt
}
