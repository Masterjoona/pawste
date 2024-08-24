package paste

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math"

	"github.com/Masterjoona/pawste/pkg/config"
	"golang.org/x/crypto/pbkdf2"
)

func SecurePassword(password string) []byte {
	salt := config.Vars.Salt
	halfLen := int(math.Ceil(float64(len(password)) / 2.0))
	return []byte(salt[:halfLen] + password + salt + password[:halfLen])
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256(SecurePassword(password))
	return hash[:]
}

func Encrypt(password string, blob *[]byte) error {
	key := deriveKey(password)
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha256.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	*blob = aesgcm.Seal((*blob)[:0], nonce, *blob, nil)
	*blob = append(*blob, nonce...)

	return nil
}

func Decrypt(password string, fileBlob []byte) ([]byte, error) {
	key := deriveKey(password)
	salt := fileBlob[len(fileBlob)-12:]
	onlyFile := fileBlob[:len(fileBlob)-12]
	str := hex.EncodeToString(salt)

	nonce, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha256.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedFileBytes, err := aesgcm.Open(nil, nonce, onlyFile, nil)
	if err != nil {
		return nil, err
	}

	copy(fileBlob[:len(decryptedFileBytes)], decryptedFileBytes)

	return fileBlob[:len(decryptedFileBytes)], nil
}

func EncryptText(password, content string) (string, error) {
	key := deriveKey(password)
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha256.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(content), nil)
	ciphertext = append(ciphertext, nonce...)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptText(password, content string) (string, error) {
	ciphertext, err := hex.DecodeString(content)
	if err != nil {
		return "", err
	}

	key := deriveKey(password)
	salt := ciphertext[len(ciphertext)-12:]
	str := hex.EncodeToString(salt)

	nonce, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha256.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext[:len(ciphertext)-12], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
