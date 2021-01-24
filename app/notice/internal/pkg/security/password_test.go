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
			name:     "正常密码",
			password: "normal_password",
			wantErr:  false,
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  false,
		},
		{
			name:     "长密码",
			password: strings.Repeat("a", 1000),
			wantErr:  false,
		},
		{
			name:     "特殊字符密码",
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
				// 检验哈希格式是否正确
				_, _, _, err := decodeArgon2idHash(hash)
				if err != nil {
					t.Errorf("HashPassword() generated invalid hash: %v", err)
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// 生成测试用的密码哈希
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
			name:       "正确密码",
			password:   password,
			storedHash: hash,
			want:       true,
			wantErr:    false,
		},
		{
			name:       "错误密码",
			password:   "wrong_password",
			storedHash: hash,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "空密码",
			password:   "",
			storedHash: hash,
			want:       false,
			wantErr:    false,
		},
		{
			name:       "错误格式哈希",
			password:   password,
			storedHash: "invalid_hash",
			want:       false,
			wantErr:    true,
		},
		{
			name:       "非Argon2id哈希",
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
	// 生成有效的哈希用于测试
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
			name:        "有效哈希",
			encodedHash: validHash,
			wantErr:     false,
		},
		{
			name:        "空哈希",
			encodedHash: "",
			wantErr:     true,
		},
		{
			name:        "格式错误",
			encodedHash: "invalid_format",
			wantErr:     true,
		},
		{
			name:        "部分格式",
			encodedHash: "$argon2id$v=19",
			wantErr:     true,
		},
		{
			name:        "非Argon2id算法",
			encodedHash: "$bcrypt$v=19$m=65536,t=3,p=4$salt$hash",
			wantErr:     true,
		},
		{
			name:        "版本格式错误",
			encodedHash: "$argon2id$version=19$m=65536,t=3,p=4$salt$hash",
			wantErr:     true,
		},
		{
			name:        "参数格式错误",
			encodedHash: "$argon2id$v=19$memory=65536,time=3,parallel=4$salt$hash",
			wantErr:     true,
		},
		{
			name:        "无效Base64盐值",
			encodedHash: "$argon2id$v=19$m=65536,t=3,p=4$invalid_base64$hash",
			wantErr:     true,
		},
		{
			name:        "无效Base64哈希",
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
	// 生成标准哈希用于测试
	password := "test_password"
	standardHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// 创建一个旧格式哈希（非Argon2id）
	oldFormatHash := "$bcrypt$somehash"

	// 创建一个无效格式哈希
	invalidHash := "invalid_hash"

	tests := []struct {
		name       string
		password   string
		storedHash string
		wantUpdate bool
		wantErr    bool
	}{
		{
			name:       "标准哈希无需更新",
			password:   password,
			storedHash: standardHash,
			wantUpdate: false,
			wantErr:    false,
		},
		{
			name:       "旧格式哈希需要更新",
			password:   password,
			storedHash: oldFormatHash,
			wantUpdate: true,
			wantErr:    false,
		},
		{
			name:       "无效哈希需要更新",
			password:   password,
			storedHash: invalidHash,
			wantUpdate: true,
			wantErr:    false,
		},
		{
			name:       "空哈希需要更新",
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
	// 直接测试非导出函数 needsRehash
	// 在真实场景中，我们通常通过测试导出函数来间接测试非导出函数
	// 但这里为了提高覆盖率，我们直接测试 needsRehash

	// 生成标准哈希用于测试
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
			name:      "标准哈希不需要重新哈希",
			hashValue: standardHash,
			want:      false,
		},
		{
			name:      "空哈希需要重新哈希",
			hashValue: "",
			want:      true,
		},
		{
			name:      "非Argon2id哈希需要重新哈希",
			hashValue: "$bcrypt$somehash",
			want:      true,
		},
		{
			name:      "无效格式需要重新哈希",
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

	// 验证参数值是否符合预期
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

// 基准测试

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
