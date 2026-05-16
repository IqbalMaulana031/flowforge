package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func DecryptPassword(ciphertextB64, privateKeyB64 string) (string, error) {
	if privateKeyB64 == "" {
		return "", errors.New("password encryption private key is not configured")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	privateKey, err := parseRSAPrivateKey(privateKeyB64)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func parseRSAPrivateKey(privateKeyB64 string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return nil, err
	}

	if block, _ := pem.Decode(keyBytes); block != nil {
		keyBytes = block.Bytes
	}

	if key, err := x509.ParsePKCS8PrivateKey(keyBytes); err == nil {
		privateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		return privateKey, nil
	}

	return x509.ParsePKCS1PrivateKey(keyBytes)
}
