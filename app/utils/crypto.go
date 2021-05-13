package utils

import (
	"crypto/aes"
	"crypto/cipher"
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

func EncryptKey(key []byte, password string, lgr *zap.Logger) {
	encryptionKey := sha3.Sum256([]byte(password))
	// Error can be safely ignored as it is only thrown if keysize if invalid
	// Which won't happen as we use SHA to generate the encryption key
	blockCipher, _ := aes.NewCipher(encryptionKey[:])
	blockCipher.Encrypt(key[:aes.BlockSize], key[:aes.BlockSize])
	blockCipher.Encrypt(key[aes.BlockSize:], key[aes.BlockSize:])
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

func CFBEncrypt(data []byte, blockCipher cipher.Block) ([]byte, error) {
	// Create dst with length of cipher blocksize + data length
	// And initialize first BlockSize bytes pseudorandom for IV
	dst := make([]byte, blockCipher.BlockSize()+len(data))
	if _, err := io.ReadFull(rand.Reader, dst[:blockCipher.BlockSize()]); err != nil {
		return nil, err
	}

	// dst from 0 to blockSize is the IV
	cfb := cipher.NewCFBEncrypter(blockCipher, dst[:blockCipher.BlockSize()])
	cfb.XORKeyStream(dst[blockCipher.BlockSize():], data)
	return dst, nil
}

func CFBDecrypt(data []byte, blockCipher cipher.Block) []byte {
	// Create CFB Decrypter with cipher, instantiating with IV (first blockSize blocks of data)
	cfb := cipher.NewCFBDecrypter(blockCipher, data[:blockCipher.BlockSize()])
	// Create variable for storing decrypted note of shorter length taking into account IV
	decrypted := make([]byte, len(data)-blockCipher.BlockSize())
	// Decrypt data starting from blockSize to decrypted
	cfb.XORKeyStream(decrypted, data[blockCipher.BlockSize():])
	return decrypted
}
