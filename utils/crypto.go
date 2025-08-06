package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptN = 1 << 15
	scryptR = 8
	scryptP = 1
	keyLen  = 32
)

func deriveKey(passphrase, salt []byte) ([]byte, error) {
	return scrypt.Key(passphrase, salt, scryptN, scryptR, scryptP, keyLen)
}

func EncryptJSON(plaintext, passphrase []byte) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	buf := bytes.Buffer{}
	buf.Write(salt)
	buf.Write(nonce)
	buf.Write(ciphertext)
	return buf.Bytes(), nil
}

func DecryptJSON(blob, passphrase []byte) ([]byte, error) {
	if len(blob) < 16 {
		return nil, errors.New("ciphertext too short")
	}
	salt := blob[:16]
	rest := blob[16:]

	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(rest) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := rest[:nonceSize]
	ciphertext := rest[nonceSize:]

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
