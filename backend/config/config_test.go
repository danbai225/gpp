package config

import "testing"

func TestParsePeer(t *testing.T) {
	err, peer := ParsePeer("Z3BwOi8vdmxlc3NAMS4yLjMuNDozNDU1NS8xMjNiMjJlZi0xMjM0LTEyMzQtMTIzNC1lZmViMjI0ZTAzZTc=")
	if err != nil {
		t.Error(err)
	}
	if peer == nil {
		t.Error("peer is nil")
	}
}
