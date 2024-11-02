package hash

import "errors"

type PasswordHash struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type PasswordParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewPasswordHash() *PasswordHash {
	return &PasswordHash{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

var (
	ErrInvalidHash         = errors.New("❌ Encoded hash is in the incorrect format")
	ErrIncompatibleVersion = errors.New("❌ Incompatible version of Argon2")
)
