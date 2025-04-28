package ssh

import "testing"

func TestGenerateAndDecodePlainKey(t *testing.T) {
	salt := "test-salt"
	pair, err := GenerateKeyPair(2048, "", salt)
	if err != nil {
		t.Fatalf("error generating plain key pair: %v", err)
	}

	privateKey, err := DecodePEMToPrivateKey(pair.private, "", salt)
	if err != nil {
		t.Fatalf("failed to decode plain key: %v", err)
	}
	if privateKey == nil {
		t.Fatal("decoded plain private key is nil")
	}
}

func TestGenerateAndDecodeEncryptedKey(t *testing.T) {
	passphrase := "test-passphrase"
	salt := "test-salt"

	pair, err := GenerateKeyPair(2048, passphrase, salt)
	if err != nil {
		t.Fatalf("error generating encrypted key pair: %v", err)
	}

	privateKey, err := DecodePEMToPrivateKey(pair.private, passphrase, salt)
	if err != nil {
		t.Fatalf("failed to decode encrypted key: %v", err)
	}
	if privateKey == nil {
		t.Fatal("decoded encrypted private key is nil")
	}
}

func TestDecodeWithWrongPassphrase(t *testing.T) {
	pair, err := GenerateKeyPair(2048, "correct-pass", "salt")
	if err != nil {
		t.Fatalf("error generating encrypted key: %v", err)
	}

	_, err = DecodePEMToPrivateKey(pair.private, "wrong-pass", "salt")
	if err == nil {
		t.Fatal("expected error when decoding with wrong passphrase")
	}
}

func TestDecodeInvalidPEM(t *testing.T) {
	_, err := DecodePEMToPrivateKey([]byte("not a pem"), "pass", "salt")
	if err == nil {
		t.Fatal("expected error on invalid PEM block")
	}
}

func TestGeneratePublicKeyFromPrivate(t *testing.T) {
	priv, err := GeneratePrivateKey(2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	pub, err := GeneratePublicKey(priv)
	if err != nil {
		t.Fatalf("failed to generate public key: %v", err)
	}

	if pub == "" {
		t.Fatal("generated public key is empty")
	}
}
