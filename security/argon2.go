package security

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
	"encoding/base64"
	"fmt"
	"github.com/alexedwards/argon2id"
)

const (
	Time	uint32 = 1
	Memory	uint32 = 64 * 1024
	Threads uint8 = 4
	KeyLen 	uint32 = 32
	SaltLen uint32 = 16
	Version int	= argon2.Version
)

func GenerateHashedPassword(rawPassword string) (string, error) {
	salt := make([]byte, SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(rawPassword), salt, Time, Memory, Threads, KeyLen)
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", Version, Memory, Time, Threads, base64Salt, base64Hash)
	
	return encodedHash, nil
}

func VerifyPassword(rawPassword string, hashedPassword string) bool {
	match, err := argon2id.ComparePasswordAndHash(rawPassword, hashedPassword)
	if err != nil {
		return false
	}
	return match
}