package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var ErrUnknownAuthType = fmt.Errorf("unknown SSH auth type")

type AuthMethodType int

const (
	AuthPassword AuthMethodType = iota
	AuthPrivateKey
)

type Auth struct {
	Type       AuthMethodType
	Password   string
	PrivateKey []byte
}

func (a *Auth) ToSSHAuthMethod() (goph.Auth, error) {
	switch a.Type {
	case AuthPassword:
		return goph.Password(a.Password), nil
	case AuthPrivateKey:
		signer, err := ssh.ParsePrivateKey(a.PrivateKey)
		if err != nil {
			return nil, err
		}
		return goph.Auth{ssh.PublicKeys(signer)}, nil
	default:
		return nil, ErrUnknownAuthType
	}
}
