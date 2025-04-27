package vault

const (
	privateKey = "private"
	publicKey  = "public"
)

type KeyPair struct {
	Private string
	Public  string
}
