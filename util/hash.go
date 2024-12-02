package util

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

func Hash(password string, salt string) ([]byte, error) {
	saltBytes := []byte(salt)
	hash, err := scrypt.Key([]byte(password), saltBytes, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("Error generating hash: ", err)
		return nil, err
	}
	fmt.Println("Hash successfully generated: ", string(hash))
	return hash, nil
}
