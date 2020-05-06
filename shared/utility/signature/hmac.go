//go:generate gobox tools/gmock

package signature

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/httpError"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

type hmacData struct {
	Message string
	Hash    []byte
}

type HmacUtil interface {
	Sign(context ctx.Context, message string) string
	Check(context ctx.Context, message string) string
}

func GetHmacUtil(context ctx.BackgroundContext) HmacUtil {
	type HmacUtilContext struct{}
	return context.Persist(HmacUtilContext{}, func() (interface{}, error) {
		return HmacUtil(hmacUtil{
			key:          []byte(configuration.GetConfig(context).HmacKey),
			errorService: loggers.GetErrorService(context),
			hasher:       GetHasher(context),
		}), nil
	}).(HmacUtil)
}

type hmacUtil struct {
	key          []byte
	errorService loggers.ErrorService
	hasher       Hasher
}

func (u hmacUtil) Sign(context ctx.Context, message string) string {
	data := hmacData{
		Message: message,
		Hash:    hmacSum(message, u.hasher.Sum(context), u.key),
	}

	b, err := json.Marshal(data)
	u.errorService.CheckErrorAndPanic(err)

	return base64.URLEncoding.EncodeToString(b)
}

func (u hmacUtil) Check(context ctx.Context, message string) string {
	b, err := base64.URLEncoding.DecodeString(message)
	checkErrorAndForbid(err)

	data := hmacData{}
	err = json.Unmarshal(b, &data)
	checkErrorAndForbid(err)

	checkBoolAndForbid(hmac.Equal(data.Hash, hmacSum(data.Message, u.hasher.Sum(context), u.key)))

	return data.Message
}

func hmacSum(message string, extraHash, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(extraHash)
	mac.Write([]byte(message))
	return mac.Sum(nil)
}

func checkErrorAndForbid(err error) {
	if err != nil {
		callHalt()
	}
}

func checkBoolAndForbid(ok bool) {
	if !ok {
		callHalt()
	}
}

func callHalt() { httpError.HaltForbidden("Hmac Signature Check Failed.") }
