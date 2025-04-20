package ssh

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh"
)

const (
	RsaPKeyType          = "RSA PRIVATE KEY"
	EncryptedRsaPKeyType = "ENCRYPTED RSA PRIVATE KEY"
)

func GeneratePrivateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  RsaPKeyType,
		Bytes: keyBytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, pemBlock); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodePrivateKeyToEncryptedPEM(privateKey *rsa.PrivateKey, passphrase, salt string) ([]byte, error) {
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	key := pbkdf2.Key([]byte(passphrase), []byte(salt), 100_000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, nonce, keyBytes, nil)

	pemBlock := &pem.Block{
		Type:  EncryptedRsaPKeyType,
		Bytes: append(nonce, cipherText...),
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, pemBlock); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodePEMToPrivateKey(pemBytes []byte, passphrase, salt string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	if block.Type == EncryptedRsaPKeyType {
		key := pbkdf2.Key([]byte(passphrase), []byte(salt), 100_000, 32, sha256.New)

		cipherBlock, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		gcm, err := cipher.NewGCM(cipherBlock)
		if err != nil {
			return nil, err
		}

		nonceSize := gcm.NonceSize()
		if len(block.Bytes) < nonceSize {
			return nil, errors.New("invalid encrypted data")
		}

		nonce := block.Bytes[:nonceSize]
		cipherText := block.Bytes[nonceSize:]
		plainText, err := gcm.Open(nil, nonce, cipherText, nil)
		if err != nil {
			return nil, err
		}

		return x509.ParsePKCS1PrivateKey(plainText)
	} else if block.Type == RsaPKeyType {
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	return nil, errors.New("unsupported key format")
}

func GeneratePublicKey(privateKey *rsa.PrivateKey) (string, error) {
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(pub)), nil
}

func GenerateKeyPair(bits int, passphrase, salt string) (privatePEM []byte, publicSSH string, err error) {
	privateKey, err := GeneratePrivateKey(bits)
	if err != nil {
		return nil, "", err
	}

	var pemBytes []byte
	if passphrase != "" {
		pemBytes, err = EncodePrivateKeyToEncryptedPEM(privateKey, passphrase, salt)
	} else {
		pemBytes, err = EncodePrivateKeyToPEM(privateKey)
	}
	if err != nil {
		return nil, "", err
	}

	pub, err := GeneratePublicKey(privateKey)
	if err != nil {
		return nil, "", err
	}

	return pemBytes, pub, nil
}
