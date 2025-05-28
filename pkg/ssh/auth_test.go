package ssh

import (
	_ "embed"
	"golang.org/x/crypto/ssh"
	"testing"
)

//go:embed test_ssh_key
var testSshKey string

func TestSSHAuth_ToSSHAuthMethod(t *testing.T) {
	type fields struct {
		Type       AuthMethodType
		Password   string
		PrivateKey []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    ssh.AuthMethod
		wantErr bool
	}{
		{
			name: "AuthPassword - success",
			fields: fields{
				Type:     AuthPassword,
				Password: "my-password",
			},
			wantErr: false,
		},
		{
			name: "AuthPrivateKey - success",
			fields: fields{
				Type:       AuthPrivateKey,
				PrivateKey: []byte(testSshKey),
			},
			wantErr: false,
		},
		{
			name: "AuthPrivateKey - invalid key",
			fields: fields{
				Type:       AuthPrivateKey,
				PrivateKey: []byte("not-a-real-key"),
			},
			wantErr: true,
		},
		{
			name: "Unknown auth type",
			fields: fields{
				Type: 99,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				Type:       tt.fields.Type,
				Password:   tt.fields.Password,
				PrivateKey: tt.fields.PrivateKey,
			}
			got, err := a.ToSSHAuthMethod()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToSSHAuthMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ToSSHAuthMethod() got is nil, expected valid AuthMethod")
			}
		})
	}
}
