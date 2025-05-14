package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"

	"golang.org/x/crypto/argon2"
)

const (
	DataSet  = "!\"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	SaltSize = 16
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var dsSize = big.NewInt(int64(len(DataSet)))

func GenSalt() (string, error) {
	salt := ""
	var n *big.Int
	var err error
	for i := 0; i < SaltSize; i++ {
		n, err = rand.Int(rand.Reader, dsSize) // gen random index
		if err != nil {
			continue
		}

		if n.Int64() >= dsSize.Int64() {
			n = big.NewInt(dsSize.Int64() - 1) // if somehow we are equal tot over the length then just be the last char.
		}

		salt += string(DataSet[n.Int64()])
	}

	if len(salt) != SaltSize {
		return "", errors.New("sroblem generating salt")
	}

	return salt, nil
}
func HashArgon(pass string, salt string) []byte {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	hash := argon2.IDKey([]byte(pass), []byte(salt), p.iterations, p.memory, p.parallelism, p.keyLength)
	return hash
}
