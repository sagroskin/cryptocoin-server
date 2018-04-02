package rsa

import (
	"crypto/sha256"
	"testing"
)

func TestGenerateNewKey(t *testing.T) {
	privateKey, err := GenerateNewKey()

	if privateKey == nil || err != nil {
		t.Error("Key not generated:", privateKey, err)
	}
}

func TestPrivateKey(t *testing.T) {
	privateKey, _ := GenerateNewKey()

	privKey := ExportPrivKey(privateKey)

	parsedPrivKey, err := ParsePrivKey(privKey)

	if parsedPrivKey.D.Cmp(privateKey.D) != 0 {
		t.Error("Key does not match:", parsedPrivKey.D, privateKey.D, err)
	}
}

func TestPublicKey(t *testing.T) {
	privateKey, _ := GenerateNewKey()

	pubKey := ExportPubKey(&privateKey.PublicKey)

	parsedPubKey, err := ParsePubKey(pubKey)

	if parsedPubKey.E != privateKey.PublicKey.E {
		t.Error("Key does not match:", parsedPubKey, privateKey.PublicKey, err)
	}
}

func TestSign(t *testing.T) {
	privateKey, _ := GenerateNewKey()
	data := []byte("Test data")
	hash := sha256.Sum256(data)

	result, err := Sign(privateKey, hash[:])

	if err != nil {
		t.Error("Signature failed:", result, err)
	}
}

func TestVerify(t *testing.T) {
	privateKey, _ := GenerateNewKey()
	data := []byte("Test data")
	hash := sha256.Sum256(data)

	signature, err := Sign(privateKey, hash[:])

	if err != nil {
		t.Error("Signature failed:", signature, err)
	}

	result, err := Verify(&privateKey.PublicKey, hash[:], signature)

	if !result {
		t.Error("Verify failed:", result, signature, err)
	}
}
