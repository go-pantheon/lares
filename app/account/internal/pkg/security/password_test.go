package security

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
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
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == tt.password {
					t.Errorf("HashPassword() failed, hash equals original password")
				}
				if !strings.HasPrefix(hash, "$argon2id$") {
					t.Errorf("HashPassword() generated invalid hash format: %s", hash)
				}
				// Check if the hash format is correct
				_, _, _, err := decodeArgon2idHash(hash)
				if err != nil {
					t.Errorf("HashPassword() generated invalid hash: %v", err)
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// Generate a test password hash
	password := "test_password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name       string
		password   string
		storedHash string
		want       bool
		wantErr    bool
	}{
		{
			name:       "correct password",
			password:   password,
			storedHash: hash,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "wrong password",
			password:   "wrong_password",
			storedHash: hash,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "empty password",
			password:   "",
			storedHash: hash,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "invalid hash format",
			password:   password,
			storedHash: "invalid_hash",
			want:       false,
			wantErr:    true,
		},
		{
			name:       "non-Argon2id hash",
			password:   password,
			storedHash: "$bcrypt$notvalid",
			want:       false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyPassword(tt.password, tt.storedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeArgon2idHash(t *testing.T) {
	// Generate a valid hash for testing
	password := "test_password"
	validHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

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
			params, salt, hash, err := decodeArgon2idHash(tt.encodedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeArgon2idHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if params == nil {
					t.Errorf("DecodeArgon2idHash() params is nil")
				}
				if len(salt) == 0 {
					t.Errorf("DecodeArgon2idHash() salt is empty")
				}
				if len(hash) == 0 {
					t.Errorf("DecodeArgon2idHash() hash is empty")
				}
			}
		})
	}
}

func TestUpdateHashIfNeeded(t *testing.T) {
	// Generate a standard hash for testing
	password := "test_password"
	standardHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

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
			newHash, updated, err := UpdateHashIfNeeded(tt.password, tt.storedHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateHashIfNeeded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if updated != tt.wantUpdate {
				t.Errorf("UpdateHashIfNeeded() updated = %v, want %v", updated, tt.wantUpdate)
			}
			if updated {
				if !strings.HasPrefix(newHash, "$argon2id$") {
					t.Errorf("UpdateHashIfNeeded() invalid new hash format: %s", newHash)
				}
			} else {
				if newHash != tt.storedHash {
					t.Errorf("UpdateHashIfNeeded() hash changed when not expected")
				}
			}
		})
	}
}

func TestNeedsRehash(t *testing.T) {
	// Generate a standard hash for testing
	password := "test_password"
	standardHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

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
			if got := needsRehash(tt.hashValue); got != tt.want {
				t.Errorf("needsRehash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultArgon2Params(t *testing.T) {
	params := DefaultArgon2Params()

	if params == nil {
		t.Fatal("DefaultArgon2Params() returned nil")
	}

	// Verify if the parameter values are as expected
	if params.Memory != 64*1024 {
		t.Errorf("DefaultArgon2Params() Memory = %v, want %v", params.Memory, 64*1024)
	}
	if params.Time != 3 {
		t.Errorf("DefaultArgon2Params() Time = %v, want %v", params.Time, 3)
	}
	if params.Threads != 4 {
		t.Errorf("DefaultArgon2Params() Threads = %v, want %v", params.Threads, 4)
	}
	if params.KeyLength != 32 {
		t.Errorf("DefaultArgon2Params() KeyLength = %v, want %v", params.KeyLength, 32)
	}
	if params.SaltLength != 16 {
		t.Errorf("DefaultArgon2Params() SaltLength = %v, want %v", params.SaltLength, 16)
	}
}

// Benchmark tests

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmark_password"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword() failed: %v", err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmark_password"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := VerifyPassword(password, hash)
		if err != nil {
			b.Fatalf("VerifyPassword() failed: %v", err)
		}
	}
}

func BenchmarkDecodeArgon2idHash(b *testing.B) {
	password := "benchmark_password"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := decodeArgon2idHash(hash)
		if err != nil {
			b.Fatalf("DecodeArgon2idHash() failed: %v", err)
		}
	}
}

func BenchmarkUpdateHashIfNeeded(b *testing.B) {
	password := "benchmark_password"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatalf("setup failed: %v", err)
	}

	b.Run("NoUpdate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := UpdateHashIfNeeded(password, hash)
			if err != nil {
				b.Fatalf("UpdateHashIfNeeded() failed: %v", err)
			}
		}
	})

	b.Run("NeedsUpdate", func(b *testing.B) {
		oldHash := "$bcrypt$somehash"
		for i := 0; i < b.N; i++ {
			_, _, err := UpdateHashIfNeeded(password, oldHash)
			if err != nil {
				b.Fatalf("UpdateHashIfNeeded() failed: %v", err)
			}
		}
	})
}
