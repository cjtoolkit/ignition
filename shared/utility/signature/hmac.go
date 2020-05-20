//go:generate gobox tools/easymock

package signature

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/configuration"
	"github.com/cjtoolkit/ignition/shared/utility/httpError"
	"github.com/cjtoolkit/ignition/shared/utility/loggers"
)

type HmacUtil interface {
	Sign(context ctx.Context, message []byte) string
	Check(context ctx.Context, message string) []byte
}

func GetHmacUtil(context ctx.BackgroundContext) HmacUtil {
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

func (u hmacUtil) Sign(context ctx.Context, message []byte) string {
	sum := hmacSum(message, u.key)

	return base64.URLEncoding.EncodeToString(append(sum, message...))
}

func (u hmacUtil) Check(context ctx.Context, message string) []byte {
	messageB, err := base64.URLEncoding.DecodeString(message)
	checkErrorAndForbid(err)
	checkErrorAndForbid(checkSize(messageB))

	currentSum := messageB[:sha512.Size]
	messageB = messageB[sha512.Size:]

	checkBoolAndForbid(hmac.Equal(currentSum, hmacSum(messageB, u.key)))

	return messageB
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
