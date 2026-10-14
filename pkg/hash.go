package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	msgerr "github.com/MGT06/e-wallet_backend/internal/message"
	"golang.org/x/crypto/argon2"
)

var ()

type ConfigHash struct {
	memory      uint32
	time        uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewConfigHash() *ConfigHash {
	return &ConfigHash{
		memory:      64 * 1024,
		time:        2,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

func (c *ConfigHash) generateRandomBytes() ([]byte, error) {
	res := make([]byte, c.saltLength)
	_, err := rand.Read(res)
	if err != nil {
		log.Printf("[PkgHash.generateRandomBytes] : %v", err)
		return nil, err
	}

	return res, nil
}

func (c *ConfigHash) HashPassword(password string) (string, error) {
	salt, err := c.generateRandomBytes()
	if err != nil {
		log.Printf("[PkgHash.HashPassword] generateRandomBytes : %v", err)
		return "", fmt.Errorf("failed to generate random salt: %w", err)
	}

	key := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.parallelism, c.keyLength)

	return c.CompleteHash(key, salt), nil
}

func (c *ConfigHash) CompleteHash(key, salt []byte) string {

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, c.memory, c.time, c.parallelism, b64Salt, b64Key)
	return hash
}

func (c *ConfigHash) DecodeHash(hash string) (config *ConfigHash, salt, key []byte, err error) {
	split := strings.Split(hash, "$")

	if len(split) != 6 {
		return nil, nil, nil, msgerr.ErrInvalidHash
	}

	if split[1] != "argon2id" {
		return nil, nil, nil, msgerr.ErrIncompatibleVariant
	}

	var version int
	if _, err := fmt.Sscanf(split[2], "v=%d", &version); err != nil {
		log.Printf("[PkgHash.DecodeHash] Sscanf 1 : %v", err)
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, msgerr.ErrIncompatibleVersion
	}

	config = &ConfigHash{}

	var memory, time uint32
	var parallelism uint8

	if _, err := fmt.Sscanf(split[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism); err != nil {
		log.Printf("[PkgHash.DecodeHash] Sscanf 2 : %v", err)
		return nil, nil, nil, err
	}

	config.memory = memory
	config.time = time
	config.parallelism = parallelism

	salt, err = base64.RawStdEncoding.DecodeString(split[4])
	if err != nil {
		log.Printf("[PkgHash.DecodeHash] DecodeString salt : %v", err)
		return nil, nil, nil, err
	}

	config.saltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.DecodeString(split[5])
	if err != nil {
		log.Printf("[PkgHash.DecodeHash] DecodeString key : %v", err)
		return nil, nil, nil, err
	}

	config.keyLength = uint32(len(key))

	return config, salt, key, nil
}

func (c *ConfigHash) Compare(password string, hash string) (string, error) {
	_, salt, key, err := c.DecodeHash(hash)
	if err != nil {
		log.Printf("[PkgHash.Compare] DecodeHash : %v", err)
		return "", err
	}

	newkey := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.parallelism, c.keyLength)

	if subtle.ConstantTimeCompare(key, newkey) == 0 {
		return "", msgerr.ChangePassword
	}
	return c.CompleteHash(newkey, salt), nil
}
