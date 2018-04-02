package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

// GenerateNewKey ?
func GenerateNewKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 512)
}

// ExportPrivKey ?
func ExportPrivKey(key *rsa.PrivateKey) string {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	return base64.StdEncoding.EncodeToString(keyBytes)
}

// ParsePrivKey ?
func ParsePrivKey(key string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	keyParsed, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return keyParsed, nil
}

// ExportPubKey ?
func ExportPubKey(key *rsa.PublicKey) string {
	keyBytes := x509.MarshalPKCS1PublicKey(key)
	return base64.StdEncoding.EncodeToString(keyBytes)
}

// ParsePubKey ?
func ParsePubKey(key string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	keyParsed, err := x509.ParsePKCS1PublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return keyParsed, nil
}

// Sign ?
func Sign(privKey *rsa.PrivateKey, hash []byte) (string, error) {
	// hashed := sha256.Sum256(data)
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])

	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(signatureBytes)

	return signature, nil
}

// Verify ?
func Verify(pubKey *rsa.PublicKey, hash []byte, signature string) (bool, error) {
	signatureBlocks, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	// hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signatureBlocks)

	if err != nil {
		return false, err
	}

	return true, nil
}
