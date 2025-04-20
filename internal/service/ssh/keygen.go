package ssh

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/pkg/ssh"
)

type KeygenService struct {
	cfg config.SSHKeygen
}

func NewKeygenService(cfg config.SSHKeygen) *KeygenService {
	return &KeygenService{cfg: cfg}
}

func (r KeygenService) Generate() string {
	ssh.GenerateKeyPair(r.cfg.Bits, r.cfg.Passphrase, r.cfg.Salt)
	return ""
}
