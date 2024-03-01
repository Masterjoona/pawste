package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password)) // wahh so secure
	return hash[:]
}

func Encrypt(file File, password string) ([]byte, error) {
	key := deriveKey(password)
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encryptedFileBytes := aesgcm.Seal(nil, nonce, file.Blob, nil)

	encryptedFileBytes = append(encryptedFileBytes, nonce...)

	return encryptedFileBytes, nil
}

func Decrypt(fileBlob []byte, password string) ([]byte, error) {
	key := deriveKey(password)
	salt := fileBlob[len(fileBlob)-12:]
	onlyFile := fileBlob[:len(fileBlob)-12]
	str := hex.EncodeToString(salt)

	nonce, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	fileBytes, err := aesgcm.Open(nil, nonce, onlyFile, nil)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func EncryptText(text string, password string) string {
	key := deriveKey(password)
	nonce := make([]byte, 12)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)

	ciphertext = append(ciphertext, nonce...)

	return hex.EncodeToString(ciphertext)
}

func DecryptText(text string, password string) string {
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		panic(err)
	}

	key := deriveKey(password)
	salt := ciphertext[len(ciphertext)-12:]
	str := hex.EncodeToString(salt)

	nonce, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext[:len(ciphertext)-12], nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}
