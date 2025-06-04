package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	_ "embed"

	"github.com/joho/godotenv"
)

//go:embed ..\.env
var envFile []byte
var (
	secretKey []byte
)

func init() {
	if err := loadEnv(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	passphrase := os.Getenv("PASSPHRASE")
	hash := sha256.Sum256([]byte(passphrase))
	secretKey = hash[:]
}

func encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("cipherText muito curto")
	}

	nonce, cipherTextBytes := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func loadEnv() error {
	envMap, err := godotenv.Parse(bytes.NewReader(envFile))
	if err != nil {
		return fmt.Errorf("failed to parse embedded .env: %w", err)
	}
	// Agora carrega as variáveis no ambiente do processo
	for k, v := range envMap {
		os.Setenv(k, v)
	}
	return nil
}
