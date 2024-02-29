package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

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

func Encrypt(file File, password string, pasteName string) {
	plaintext, err := source.Open()
	if err != nil {
		panic(err)
	}
	plaintextBytes, err := io.ReadAll(plaintext)
	if err != nil {
		panic(err)
	}

	key := deriveKey(password)
	nonce := make([]byte, 12)

	// Randomizing the nonce
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

	ciphertext := aesgcm.Seal(nil, nonce, plaintextBytes, nil)

	ciphertext = append(ciphertext, nonce...)

	f, err := os.Create(Config.DataDir + pasteName + "/" + source.Filename + ".enc")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	if err != nil {
		panic(err)
	}
}

func DecryptFile(source string, password string) {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		panic(err)
	}

	ciphertext, err := os.ReadFile(source)

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

	f, err := os.Create(strings.TrimSuffix(source, ".enc"))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, bytes.NewReader(plaintext))
	if err != nil {
		panic(err)
	}
}

func EncryptText(text string, password string) string {
	key := deriveKey(password)
	nonce := make([]byte, 12)

	// Randomizing the nonce
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
