package ssh

import (
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/pkg/ssh"
)

type KeygenService struct {
	cfg config.SSHKeygen
}

func NewKeygenService(cfg config.SSHKeygen) *KeygenService {
	return &KeygenService{cfg: cfg}
}

func (r *KeygenService) GeneratePair() (ssh.KeyPair, error) {
	pair, err := ssh.GenerateKeyPair(r.cfg.Bits, r.cfg.Passphrase, r.cfg.Salt)
	if err != nil {
		return ssh.KeyPair{}, fmt.Errorf("error on generating key-pair: %w", err)
	}
	return pair, nil
}
