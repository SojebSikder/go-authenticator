package totp

import (
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func normalizeSecret(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", ""))
}

func CurrentCode(secret string) (string, error) {
	return totp.GenerateCode(normalizeSecret(secret), time.Now())
}

func ValidateCode(secret, code string) bool {
	ok, _ := totp.ValidateCustom(code, normalizeSecret(secret), time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return ok
}
