package security

import (
	"sync"
	"sync/atomic"

	"github.com/go-pantheon/fabrica-util/security/aes"
	"github.com/go-pantheon/lares/app/account/internal/conf"
)

var (
	once    sync.Once
	manager = &CryptoManager{}
)

type CryptoManager struct {
	mu          sync.RWMutex
	initialized atomic.Bool

	tokenAES    *aes.Cipher
	sessionAES  *aes.Cipher
	platformAES *aes.Cipher
}

func Init(config *conf.Secret) (err error) {
	once.Do(func() {
		err = initCryptoManager(config)
	})

	return err
}

func initCryptoManager(c *conf.Secret) (err error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	tokenCipher, err := aes.NewAESCipher([]byte(c.TokenKey))
	if err != nil {
		return err
	}

	sessionCipher, err := aes.NewAESCipher([]byte(c.SessionKey))
	if err != nil {
		return err
	}

	platformCipher, err := aes.NewAESCipher([]byte(c.PlatformKey))
	if err != nil {
		return err
	}

	manager.tokenAES = tokenCipher
	manager.sessionAES = sessionCipher
	manager.platformAES = platformCipher
	manager.initialized.Store(true)

	return nil
}
