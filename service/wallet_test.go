package service

import "testing"

func TestGetWallet(t *testing.T) {
	w, err := GetWallet("J2nHLtdwZFmxbAe3oniv40NOrekJ1B/tRxu1J2xDJ+n7vGvYoqm4EJLoJUSC9pnTSNHh3dMKBpumEkfynd1huA==")

	if w == nil || err != nil {
		t.Error("GetWallet failed:", w, err)
	}
}
