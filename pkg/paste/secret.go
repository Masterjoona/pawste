package paste

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/Masterjoona/pawste/pkg/shared/config"
	"golang.org/x/crypto/pbkdf2"
)

func SecurePassword(password string) []byte {
	return []byte(config.Config.Salt[:len(password)/2] + password + config.Config.Salt + password[:len(password)/2])
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256(SecurePassword(password))
	return hash[:]
}

func (f *File) Encrypt(password string) error {
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

	f.Blob = aesgcm.Seal(f.Blob[:0], nonce, f.Blob, nil)
	f.Blob = append(f.Blob, nonce...)
	return nil
}

func (f *File) Decrypt(password string) ([]byte, error) {
	fileBlob := f.Blob
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

func (p *Paste) EncryptText() error {
	key := deriveKey(p.Password)
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

	ciphertext := aesgcm.Seal(nil, nonce, []byte(p.Content), nil)
	ciphertext = append(ciphertext, nonce...)

	p.Content = hex.EncodeToString(ciphertext)
	return nil
}

func (p *Paste) DecryptText() string {
	ciphertext, err := hex.DecodeString(p.Content)
	if err != nil {
		panic(err)
	}

	key := deriveKey(p.Password)
	salt := ciphertext[len(ciphertext)-12:]
	str := hex.EncodeToString(salt)

	nonce, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha256.New)

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
