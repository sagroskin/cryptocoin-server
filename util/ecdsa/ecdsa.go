package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"math/big"
)

// GenerateNewKey generates a new ECDSA private and public key using the P-256 curve.
func GenerateNewKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// ExportPrivKey exports the ECDSA private key as a Base64 encoded string.
func ExportPrivKey(key *ecdsa.PrivateKey) string {
	keyBytes, _ := x509.MarshalECPrivateKey(key)
	return base64.StdEncoding.EncodeToString(keyBytes)
}

// ParsePrivKey parses a Base64 string to an ECDSA private key.
func ParsePrivKey(key string) (*ecdsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	keyParsed, err := x509.ParseECPrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}

	return keyParsed, nil
}

// ExportPubKey exports the ECDSA public key as a Base64 encoded string.
func ExportPubKey(key *ecdsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(append(key.X.Bytes(), key.Y.Bytes()...))
}

// ParsePubKey parses a Base64 string to an ECDSA public key.
func ParsePubKey(key string) (*ecdsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	x := big.Int{}
	y := big.Int{}
	sigLen := len(keyBytes)
	x.SetBytes(keyBytes[:(sigLen / 2)])
	y.SetBytes(keyBytes[(sigLen / 2):])

	keyParsed := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	if err != nil {
		return nil, err
	}

	return &keyParsed, nil
}

// Sign creates a signature using a SHA256 hash and ECDSA private key.
func Sign(privKey *ecdsa.PrivateKey, hash []byte) (string, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])

	if err != nil {
		return "", err
	}

	signatureBytes := r.Bytes()
	signatureBytes = append(signatureBytes, s.Bytes()...)

	signature := base64.StdEncoding.EncodeToString(signatureBytes)

	return signature, nil
}

// Verify verifies a signature using a SHA256 hash, provided signature, and ECDSA public key.
func Verify(pubKey *ecdsa.PublicKey, hash []byte, signature string) (bool, error) {
	signatureBlocks, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	r := big.Int{}
	s := big.Int{}
	sigLen := len(signatureBlocks)
	r.SetBytes(signatureBlocks[:(sigLen / 2)])
	s.SetBytes(signatureBlocks[(sigLen / 2):])

	result := ecdsa.Verify(pubKey, hash[:], &r, &s)

	return result, nil
}
