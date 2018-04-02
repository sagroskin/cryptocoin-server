package model

import (
	"bytes"
	"testing"
)

func TestSerialize(t *testing.T) {
	t1 := Transaction{Value: 1, PubKey: "123"}
	t2 := Transaction{Value: 1, PubKey: "123"}

	bytes1, err1 := t1.Serialize()
	bytes2, err2 := t2.Serialize()

	if bytes.Compare(bytes1, bytes2) != 0 || err1 != nil || err2 != nil {
		t.Error("Bytes do not match:", bytes1, bytes2, err1, err2)
	}
}

func TestHash(t *testing.T) {
	t1 := Transaction{Value: 1, PubKey: "123"}
	t2 := Transaction{Value: 1, PubKey: "123"}

	h1, err1 := t1.Hash()
	h2, err2 := t2.Hash()

	if bytes.Compare(h1, h2) != 0 || err1 != nil || err2 != nil {
		t.Error("Hash does not match:", h1, h2, err1, err2)
	}
}
