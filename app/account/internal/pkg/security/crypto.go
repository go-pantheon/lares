package security

import (
	"crypto/cipher"
	"encoding/base64"

	"github.com/go-pantheon/fabrica-util/security/aes"
	"github.com/pkg/errors"
)

var (
	tokenAESKey   []byte
	tokenAESBlock cipher.Block

	sessionAESKey   []byte
	sessionAESBlock cipher.Block

	platformAESKey   []byte
	platformAESBlock cipher.Block
)

func Init(tokenKey, sessionKey, platformKey string) (err error) {
	tokenAESKey = []byte(tokenKey)
	if tokenAESBlock, err = aes.NewBlock(tokenAESKey); err != nil {
		return err
	}
	sessionAESKey = []byte(sessionKey)
	if sessionAESBlock, err = aes.NewBlock(sessionAESKey); err != nil {
		return err
	}
	platformAESKey = []byte(platformKey)
	if platformAESBlock, err = aes.NewBlock(platformAESKey); err != nil {
		return err
	}
	return
}

func EncryptToken(org []byte) (string, error) {
	secret, err := aes.Encrypt(tokenAESKey, tokenAESBlock, org)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secret), nil
}

func EncryptSession(org []byte) (string, error) {
	ser, err := aes.Encrypt(sessionAESKey, sessionAESBlock, org)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(ser), nil
}

func DecryptSession(secret string) (org []byte, err error) {
	ser, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		err = errors.Wrapf(err, "session key base64 decode failed.")
		return
	}
	org, err = aes.Decrypt(sessionAESKey, sessionAESBlock, ser)
	if err != nil {
		return
	}
	return
}

func EncryptPlatform(org []byte) (string, error) {
	ser, err := aes.Encrypt(platformAESKey, platformAESBlock, org)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(ser), nil
}

func DecryptPlatform(secret string) (org []byte, err error) {
	ser, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		err = errors.Wrapf(err, "platform key base64 decode failed.")
		return
	}
	org, err = aes.Decrypt(platformAESKey, platformAESBlock, ser)
	if err != nil {
		return
	}
	return
}
