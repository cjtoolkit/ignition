//go:generate gobox tools/easymock

package signature

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/httpError"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

type HmacUtil interface {
	SignWithKey(key []byte, message []byte) []byte
	Sign(message []byte) []byte
	CheckWithKey(key []byte, message []byte) []byte
	Check(message []byte) []byte
}

func GetHmacUtil(context ctx.Context) HmacUtil {
	type HmacUtilContext struct{}
	return context.Persist(HmacUtilContext{}, func() (interface{}, error) {
		return HmacUtil(hmacUtil{
			key:          convertToByte(configuration.GetConfig(context).HmacKey),
			errorService: loggers.GetErrorService(context),
		}), nil
	}).(HmacUtil)
}

type hmacUtil struct {
	key          []byte
	errorService loggers.ErrorService
}

func (u hmacUtil) SignWithKey(key []byte, message []byte) []byte {
	sum := hmacSum(message, key)

	return append(sum, message...)
}

func (u hmacUtil) Sign(message []byte) []byte {
	return u.SignWithKey(u.key, message)
}

func (u hmacUtil) CheckWithKey(key []byte, message []byte) []byte {
	checkErrorAndForbid(checkSize(message))

	currentSum := message[:sha512.Size]
	message = message[sha512.Size:]

	checkBoolAndForbid(hmac.Equal(currentSum, hmacSum(message, key)))

	return message
}

func (u hmacUtil) Check(message []byte) []byte {
	return u.CheckWithKey(u.key, message)
}

func hmacSum(message []byte, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
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

func convertToByte(hmacKeyStr string) []byte {
	hmacKey, err := hex.DecodeString(hmacKeyStr)
	if err != nil {
		hmacKey = []byte(hmacKeyStr)
	}
	return hmacKey
}

func checkSize(message []byte) error {
	if len(message) < sha512.Size {
		return fmt.Errorf("it is less '%d' bytes in size", sha512.Size)
	}

	return nil
}
