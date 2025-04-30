package connections

import "github.com/imDrOne/minecraft-server-manager/internal/pkg/ssh"

type ConnectionSshKeyPair struct {
	privatePem []byte
	public     string
}

func NewConnSshKeyPair(privatePem []byte, public string) *ConnectionSshKeyPair {
	return &ConnectionSshKeyPair{privatePem: privatePem, public: public}
}

func ConnSshKeysFromPair(pair ssh.KeyPair) *ConnectionSshKeyPair {
	return &ConnectionSshKeyPair{privatePem: pair.Private, public: pair.Public}
}

func (r ConnectionSshKeyPair) PrivatePem() []byte {
	return r.privatePem
}

func (r ConnectionSshKeyPair) PrivatePemStr() string {
	return string(r.privatePem)
}

func (r ConnectionSshKeyPair) Public() string {
	return r.public
}
