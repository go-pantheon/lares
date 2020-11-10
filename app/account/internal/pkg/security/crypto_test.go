package security

import (
	"encoding/base64"
	"testing"

	jsoniter "github.com/json-iterator/go"
	v1 "github.com/luffy050596/rec-account/gen/api/server/account/interface/account/v1"
	"github.com/luffy050596/rec-util/security/aes"
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
