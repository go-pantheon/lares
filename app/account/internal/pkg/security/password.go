package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

// Argon2Params includes the parameters for the Argon2 hash function
type Argon2Params struct {
	Memory     uint32 // memory consumption size (KB)
	Time       uint32 // iteration count
	Threads    uint8  // parallelism
	KeyLength  uint32 // output key length
	SaltLength int    // salt length
}

// DefaultArgon2Params returns the recommended Argon2id parameter settings
func DefaultArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:     64 * 1024, // 64MB
		Time:       3,         // 3 iterations
		Threads:    4,         // 4 threads
		KeyLength:  32,        // 32 bytes key
		SaltLength: 16,        // 16 bytes salt
	}
}

// HashPassword uses the Argon2id algorithm to hash the password
// Returns a formatted hash string that contains all the information needed for verification
func HashPassword(password string) (string, error) {
	// Use default parameters
	params := DefaultArgon2Params()

	// Generate random salt
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", errors.Wrapf(err, "failed to generate salt")
	}

	// Calculate hash using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLength,
	)

	// Encode parameters and hash result into a storable format
	// Format: $argon2id$v=19$m=65536,t=3,p=4$[salt]$[hash]
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

// VerifyPassword verifies plaintext password against stored hash
// Supports checking old format hash but uses Argon2id for verification
func VerifyPassword(password, storedHash string) (bool, error) {
	// If not Argon2id hash format, return error
	if len(storedHash) <= 9 || storedHash[0:9] != "$argon2id" {
		return false, errors.New("unsupported hash format, system uses Argon2id algorithm")
	}

	// Parse hash parameters and values
	params, salt, hash, err := decodeArgon2idHash(storedHash)
	if err != nil {
		return false, errors.Wrapf(err, "failed to parse hash")
	}

	// Calculate hash using same parameters
	calcHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLength,
	)

	// Use constant time comparison to resist timing attacks
	return subtle.ConstantTimeCompare(hash, calcHash) == 1, nil
}

// decodeArgon2idHash parses the argon2id hash string into its components
func decodeArgon2idHash(encodedHash string) (params *Argon2Params, salt, hash []byte, err error) {
	// Split the hash string by the $ symbol
	// Format: $argon2id$v=19$m=65536,t=3,p=4$[salt]$[hash]
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format, expected format: $argon2id$v=[version]$m=[memory],t=[time],p=[threads]$[salt]$[hash]")
	}

	// Check algorithm type
	if parts[1] != "argon2id" {
		return nil, nil, nil, fmt.Errorf("unsupported algorithm: %s, only supports argon2id", parts[1])
	}

	// Parse version number
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse version number: %w", err)
	}

	// Parse parameters
	params = &Argon2Params{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Decode salt
	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	params.SaltLength = len(salt)

	// Decode hash
	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode hash: %w", err)
	}
	lenHash := len(hash)
	if lenHash > math.MaxUint32 {
		return nil, nil, nil, fmt.Errorf("hash length is too long: %d", lenHash)
	}
	params.KeyLength = uint32(lenHash)

	return params, salt, hash, nil
}

// UpdateHashIfNeeded checks if the password hash needs to be upgraded
// Returns: new/old hash, whether it has been updated, error
func UpdateHashIfNeeded(password, storedHash string) (string, bool, error) {
	// Check if the hash needs to be updated
	if needsRehash(storedHash) {
		// Create a new hash using the latest Argon2id parameters
		newHash, err := HashPassword(password)
		return newHash, true, err
	}
	return storedHash, false, nil
}

// needsRehash checks if the password hash needs to be rehashed using an updated algorithm or parameters
func needsRehash(hash string) bool {
	// If the hash is empty, it needs to be rehashed
	if hash == "" {
		return true
	}

	// If the hash is not in Argon2id format, it needs to be rehashed
	if len(hash) <= 9 || hash[0:9] != "$argon2id" {
		return true
	}

	// Check if the Argon2id parameters meet the current standard
	params, _, _, err := decodeArgon2idHash(hash)
	if err != nil {
		return true // Parsing failed, needs to be rehashed
	}

	// Get the current recommended parameters
	recommended := DefaultArgon2Params()

	// If the current hash uses parameters below the recommended values, it needs to be updated
	return params.Memory < recommended.Memory ||
		params.Time < recommended.Time ||
		params.Threads < recommended.Threads ||
		params.KeyLength < recommended.KeyLength
}
