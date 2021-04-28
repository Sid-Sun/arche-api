package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"io"
)

func GenerateEncryptionKey(lgr *zap.Logger) ([]byte, [32]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		// TODO: Add Logging
		return nil, [32]byte{}, err
	}

	return key, sha3.Sum256(key), nil
}

func EncryptKey(key []byte, password string, lgr *zap.Logger) error {
	encryptionKey := sha3.Sum256([]byte(password))
	blockCipher, err := aes.NewCipher(encryptionKey[:])
	if err != nil {
		// TODO: Add Logging
		return err
	}

	blockCipher.Encrypt(key[:aes.BlockSize], key[:aes.BlockSize])
	blockCipher.Encrypt(key[aes.BlockSize:], key[aes.BlockSize:])

	return nil
}

func VerifyKeyIntegrity(key []byte, hash []byte) {
	bytes.Equal(key, hash)
}

func DecryptKey(key []byte, password string, lgr *zap.Logger) error {
	encryptionKey := sha3.Sum256([]byte(password))
	blockCipher, err := aes.NewCipher(encryptionKey[:])
	if err != nil {
		// TODO: Add Logging
		return err
	}

	blockCipher.Decrypt(key[:aes.BlockSize], key[:aes.BlockSize])
	blockCipher.Decrypt(key[aes.BlockSize:], key[aes.BlockSize:])

	return nil
}
