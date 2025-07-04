package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "normal password",
			password: "normal_password",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 1000),
			wantErr:  false,
		},
		{
			name:     "special character password",
			password: "p@$$w0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash, err := HashPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NotEqual(t, hash, tt.password)
			assert.True(t, strings.HasPrefix(hash, "$argon2id$"))

			_, _, _, err = decodeArgon2idHash(hash)
			assert.NoError(t, err)
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	t.Parallel()

	password := "test_password"

	hash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name       string
		password   string
		storedHash string
		verified   bool
		wantErr    bool
	}{
		{
			name:       "correct password",
			password:   password,
			storedHash: hash,
			verified:   true,
			wantErr:    false,
		},
		{
			name:       "wrong password",
			password:   "wrong_password",
			storedHash: hash,
			verified:   false,
			wantErr:    false,
		},
		{
			name:       "empty password",
			password:   "",
			storedHash: hash,
			verified:   false,
			wantErr:    false,
		},
		{
			name:       "invalid hash format",
			password:   password,
			storedHash: "invalid_hash",
			verified:   false,
			wantErr:    true,
		},
		{
			name:       "non-Argon2id hash",
			password:   password,
			storedHash: "$bcrypt$notvalid",
			verified:   false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			verifed, err := VerifyPassword(tt.password, tt.storedHash)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.verified, verifed)
		})
	}
}

func TestDecodeArgon2idHash(t *testing.T) {
	t.Parallel()

	// Generate a valid hash for testing
	password := "test_password"

	validHash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name        string
		encodedHash string
		wantErr     bool
	}{
		{
			name:        "valid hash",
			encodedHash: validHash,
			wantErr:     false,
		},
		{
			name:        "empty hash",
			encodedHash: "",
			wantErr:     true,
		},
		{
			name:        "invalid format",
			encodedHash: "invalid_format",
			wantErr:     true,
		},
		{
			name:        "partial format",
			encodedHash: "$argon2id$v=19",
			wantErr:     true,
		},
		{
			name:        "non-Argon2id algorithm",
			encodedHash: "$bcrypt$v=19$m=65536,t=3,p=4$salt$hash",
			wantErr:     true,
		},
		{
			name:        "version format error",
			encodedHash: "$argon2id$version=19$m=65536,t=3,p=4$salt$hash",
			wantErr:     true,
		},
		{
			name:        "parameter format error",
			encodedHash: "$argon2id$v=19$memory=65536,time=3,parallel=4$salt$hash",
			wantErr:     true,
		},
		{
			name:        "invalid Base64 salt",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$invalid_base64$hash",
			wantErr:     true,
		},
		{
			name:        "invalid Base64 hash",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$c2FsdA$invalid_base64",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			params, salt, hash, err := decodeArgon2idHash(tt.encodedHash)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, uint32(64*1024), params.Memory)
			assert.Equal(t, uint32(3), params.Time)
			assert.Equal(t, uint8(4), params.Threads)
			assert.Equal(t, uint32(32), params.KeyLength)
			assert.Equal(t, 16, params.SaltLength)
			assert.NotEmpty(t, salt)
			assert.NotEmpty(t, hash)
		})
	}
}

func TestUpdateHashIfNeeded(t *testing.T) {
	t.Parallel()

	// Generate a standard hash for testing
	password := "test_password"

	standardHash, err := HashPassword(password)
	require.NoError(t, err)

	// Create an old format hash (non-Argon2id)
	oldFormatHash := "$bcrypt$somehash"

	// Create an invalid format hash
	invalidHash := "invalid_hash"

	tests := []struct {
		name       string
		password   string
		storedHash string
		wantUpdate bool
		wantErr    bool
	}{
		{
			name:       "standard hash no need to update",
			password:   password,
			storedHash: standardHash,
			wantUpdate: false,
			wantErr:    false,
		},
		{
			name:       "old format hash need to update",
			password:   password,
			storedHash: oldFormatHash,
			wantUpdate: true,
			wantErr:    false,
		},
		{
			name:       "invalid hash need to update",
			password:   password,
			storedHash: invalidHash,
			wantUpdate: true,
			wantErr:    false,
		},
		{
			name:       "empty hash need to update",
			password:   password,
			storedHash: "",
			wantUpdate: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newHash, updated, err := UpdateHashIfNeeded(tt.password, tt.storedHash)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.wantUpdate, updated)
			assert.True(t, strings.HasPrefix(newHash, "$argon2id$"))
		})
	}
}

func TestNeedsRehash(t *testing.T) {
	t.Parallel()

	// Generate a standard hash for testing
	password := "test_password"

	standardHash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name      string
		hashValue string
		want      bool
	}{
		{
			name:      "standard hash no need to rehash",
			hashValue: standardHash,
			want:      false,
		},
		{
			name:      "empty hash need to rehash",
			hashValue: "",
			want:      true,
		},
		{
			name:      "non-Argon2id hash need to rehash",
			hashValue: "$bcrypt$somehash",
			want:      true,
		},
		{
			name:      "invalid format need to rehash",
			hashValue: "invalid_hash",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, needsRehash(tt.hashValue))
		})
	}
}

func TestDefaultArgon2Params(t *testing.T) {
	t.Parallel()

	params := DefaultArgon2Params()
	require.NotNil(t, params)

	// Verify if the parameter values are as expected
	assert.Equal(t, uint32(64*1024), params.Memory)
	assert.Equal(t, uint32(3), params.Time)
	assert.Equal(t, uint8(4), params.Threads)
	assert.Equal(t, uint32(32), params.KeyLength)
	assert.Equal(t, 16, params.SaltLength)
}

// Benchmark tests

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmark_password"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = HashPassword(password)
		}
	})
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmark_password"

	hash, err := HashPassword(password)
	require.NoError(b, err)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = VerifyPassword(password, hash)
		}
	})
}

func BenchmarkDecodeArgon2idHash(b *testing.B) {
	password := "benchmark_password"

	hash, err := HashPassword(password)
	require.NoError(b, err)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _, _ = decodeArgon2idHash(hash)
		}
	})
}

func BenchmarkUpdateHashIfNeeded(b *testing.B) {
	password := "benchmark_password"

	hash, err := HashPassword(password)
	require.NoError(b, err)

	b.Run("NoUpdate", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, _ = UpdateHashIfNeeded(password, hash)
			}
		})
	})

	b.Run("NeedsUpdate", func(b *testing.B) {
		oldHash := "$bcrypt$somehash"

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, _ = UpdateHashIfNeeded(password, oldHash)
			}
		})
	})
}
