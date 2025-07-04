package security

import (
	"encoding/base64"

	"github.com/go-pantheon/fabrica-util/errors"
)

func EncryptToken(org []byte) (string, error) {
	if !manager.initialized.Load() {
		return "", errors.New("crypto manager not initialized, call Init() first")
	}

	manager.mu.RLock()
	tokenAESCipher := manager.tokenAES
	manager.mu.RUnlock()

	secret, err := tokenAESCipher.Encrypt(org)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}

func EncryptSession(org []byte) (string, error) {
	if !manager.initialized.Load() {
		return "", errors.New("crypto manager not initialized, call Init() first")
	}

	manager.mu.RLock()
	sessionAESCipher := manager.sessionAES
	manager.mu.RUnlock()

	secret, err := sessionAESCipher.Encrypt(org)
	if err != nil {
		return "", errors.Wrap(err, "session aes encrypt failed")
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}

func DecryptSession(secret string) (org []byte, err error) {
	if !manager.initialized.Load() {
		return nil, errors.New("crypto manager not initialized, call Init() first")
	}

	ser, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, errors.Wrap(err, "base64 DecodeString failed")
	}

	manager.mu.RLock()
	sessionAESCipher := manager.sessionAES
	manager.mu.RUnlock()

	org, err = sessionAESCipher.Decrypt(ser)
	if err != nil {
		return nil, errors.Wrap(err, "session aes decrypt failed")
	}

	return org, nil
}

func EncryptPlatform(org []byte) (string, error) {
	if !manager.initialized.Load() {
		return "", errors.New("crypto manager not initialized, call Init() first")
	}

	manager.mu.RLock()
	platformAESCipher := manager.platformAES
	manager.mu.RUnlock()

	secret, err := platformAESCipher.Encrypt(org)
	if err != nil {
		return "", errors.Wrap(err, "platform aes encrypt failed")
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}

func DecryptPlatform(secret string) (org []byte, err error) {
	if !manager.initialized.Load() {
		return nil, errors.New("crypto manager not initialized, call Init() first")
	}

	ser, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, errors.Wrap(err, "base64 DecodeString failed")
	}

	manager.mu.RLock()
	platformAESCipher := manager.platformAES
	manager.mu.RUnlock()

	org, err = platformAESCipher.Decrypt(ser)
	if err != nil {
		return nil, errors.Wrap(err, "platform aes decrypt failed")
	}

	return org, nil
}
