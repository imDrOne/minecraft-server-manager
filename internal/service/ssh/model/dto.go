package model

import "github.com/imDrOne/minecraft-server-manager/pkg/ssh"

type NodeSSHConnectionTO struct {
	NodeId int64
	Auth   ssh.Auth
	Host   string
	Port   uint
	User   string
}

func (r *NodeSSHConnectionTO) WithAuth(auth ssh.Auth) NodeSSHConnectionTO {
	return NodeSSHConnectionTO{
		Auth: auth,
		Host: r.Host,
		Port: r.Port,
		User: r.User,
	}
}
