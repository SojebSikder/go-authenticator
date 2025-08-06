package qrcodeutil

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/tuotoo/qrcode"
)

func DecodeFromFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	qrMatrix, err := qrcode.Decode(f)
	if err != nil {
		return "", err
	}
	return qrMatrix.Content, nil
}

func ParseOtpauth(uri string) (string, string, error) {
	if !strings.HasPrefix(uri, "otpauth://") {
		return "", "", errors.New("not an otpauth URI")
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}

	if parsed.Scheme != "otpauth" || parsed.Host != "totp" {
		return "", "", errors.New("unsupported otpauth type (only totp supported)")
	}

	account := strings.TrimPrefix(parsed.Path, "/")
	secret := parsed.Query().Get("secret")
	if secret == "" {
		return "", "", errors.New("no secret in otpauth URI")
	}
	return account, secret, nil
}
