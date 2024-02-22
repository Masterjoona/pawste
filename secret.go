package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(Config.PasswordSalt + password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func EncryptFile(file *multipart.FileHeader, password string) error {
	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	// Create a new file for writing the encrypted data
	encryptedFile, err := os.Create(file.Filename + ".encrypted")
	if err != nil {
		return err
	}
	defer encryptedFile.Close()

	key := deriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	if _, err := encryptedFile.Write(iv); err != nil {
		return err
	}

	stream, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	if _, err := encryptedFile.Write(stream.Seal(nil, iv, make([]byte, 0), nil)); err != nil {
		return err
	}

	if _, err := io.Copy(encryptedFile, uploadedFile); err != nil {
		return err
	}

	return nil
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
