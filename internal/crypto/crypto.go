package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	DataSet  = "!\"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	SaltSize = 32
)

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
