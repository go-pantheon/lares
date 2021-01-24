package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

// Argon2Params 包含Argon2哈希函数的参数
type Argon2Params struct {
	Memory     uint32 // 内存消耗大小(KB)
	Time       uint32 // 迭代次数
	Threads    uint8  // 并行度
	KeyLength  uint32 // 输出密钥长度
	SaltLength uint32 // 盐值长度
}

// DefaultArgon2Params 返回当前推荐的Argon2id参数设置
func DefaultArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:     64 * 1024, // 64MB
		Time:       3,         // 3次迭代
		Threads:    4,         // 4个线程
		KeyLength:  32,        // 32字节的密钥
		SaltLength: 16,        // 16字节的盐值
	}
}

// HashPassword 使用Argon2id算法对密码进行哈希处理
// 返回格式化的哈希字符串，该字符串包含所有验证所需的信息
func HashPassword(password string) (string, error) {
	// 使用默认参数
	params := DefaultArgon2Params()

	// 生成随机盐值
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", errors.Wrapf(err, "生成盐值失败")
	}

	// 使用Argon2id计算哈希
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLength,
	)

	// 将参数和哈希结果编码成可存储的格式
	// 格式: $argon2id$v=19$m=65536,t=3,p=4$[salt]$[hash]
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Time,
		params.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodedHash, nil
}

// VerifyPassword 验证明文密码是否与存储的哈希匹配
// 支持检查旧格式哈希但统一使用Argon2id进行验证
func VerifyPassword(password, storedHash string) (bool, error) {
	// 如果不是Argon2id哈希格式，返回错误
	if len(storedHash) <= 9 || storedHash[0:9] != "$argon2id" {
		return false, errors.New("不支持的哈希格式，系统已统一使用Argon2id算法")
	}

	// 解析哈希参数和值
	params, salt, hash, err := decodeArgon2idHash(storedHash)
	if err != nil {
		return false, errors.Wrapf(err, "解析哈希失败")
	}

	// 使用相同参数计算哈希
	calcHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLength,
	)

	// 使用恒定时间比较以抵抗计时攻击
	return subtle.ConstantTimeCompare(hash, calcHash) == 1, nil
}

// decodeArgon2idHash 解析argon2id哈希字符串到其组成部分
func decodeArgon2idHash(encodedHash string) (params *Argon2Params, salt, hash []byte, err error) {
	// 按照$符号分割哈希字符串
	// 格式: $argon2id$v=19$m=65536,t=3,p=4$[salt]$[hash]
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, fmt.Errorf("格式无效的散列，预期格式: $argon2id$v=[version]$m=[memory],t=[time],p=[threads]$[salt]$[hash]")
	}

	// 检查算法类型
	if parts[1] != "argon2id" {
		return nil, nil, nil, fmt.Errorf("不支持的算法: %s, 只支持argon2id", parts[1])
	}

	// 解析版本号
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("无法解析版本号: %w", err)
	}

	// 解析参数
	params = &Argon2Params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("无法解析参数: %w", err)
	}

	// 解码盐值
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("无法解码盐值: %w", err)
	}
	params.SaltLength = uint32(len(salt))

	// 解码哈希
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("无法解码哈希: %w", err)
	}
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}

// UpdateHashIfNeeded 检查密码哈希是否需要升级，如需要则返回新哈希
// 返回值: 新/旧哈希, 是否已更新, 错误
func UpdateHashIfNeeded(password, storedHash string) (string, bool, error) {
	// 检查哈希是否需要更新
	if needsRehash(storedHash) {
		// 使用最新的Argon2id参数创建新哈希
		newHash, err := HashPassword(password)
		return newHash, true, err
	}
	return storedHash, false, nil
}

// needsRehash 检查密码哈希是否需要使用更新的算法或参数重新哈希
func needsRehash(hash string) bool {
	// 如果是空哈希值，需要重新哈希
	if hash == "" {
		return true
	}

	// 如果不是Argon2id哈希格式，需要统一升级
	if len(hash) <= 9 || hash[0:9] != "$argon2id" {
		return true
	}

	// 检查Argon2id参数是否符合当前标准
	params, _, _, err := decodeArgon2idHash(hash)
	if err != nil {
		return true // 解析失败，需要重新哈希
	}

	// 获取当前推荐参数
	recommended := DefaultArgon2Params()

	// 如果当前哈希使用的参数低于推荐值，需要更新
	return params.Memory < recommended.Memory ||
		params.Time < recommended.Time ||
		params.Threads < recommended.Threads ||
		params.KeyLength < recommended.KeyLength
}
