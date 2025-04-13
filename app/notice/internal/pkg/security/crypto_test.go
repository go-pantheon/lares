package security

import (
	"encoding/base64"
	"testing"

	"github.com/go-pantheon/fabrica-util/security/aes"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestDecryptPlatform(t *testing.T) {
	key := "test-example-key"
	appleId := "test-example-id-1"
	state := "test-state-abc"

	err := Init(key, key, key)
	assert.Nil(t, err)

	org, err := jsoniter.Marshal(&v1.AppleId{AppleId: appleId, State: state})
	assert.Nil(t, err)

	ser, err := aes.Encrypt(platformAESKey, platformAESBlock, org)
	assert.Nil(t, err)

	secret := base64.RawURLEncoding.EncodeToString(ser)
	data, err := DecryptPlatform(secret)
	assert.Nil(t, err)

	t.Log(secret)

	dest := &v1.AppleId{}
	err = jsoniter.Unmarshal(data, dest)
	assert.Nil(t, err)

	assert.Equal(t, appleId, dest.AppleId)
	assert.Equal(t, state, dest.State)
}
