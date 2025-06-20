package connections

import "time"

type CreateConnectionRequestDto struct {
	NodeId int64  `json:"nodeId"`
	User   string `json:"user"`
}

type UpdateConnectionRequestDto struct {
	Key  string `json:"key"`
	User string `json:"user"`
}

type ConnectionResponseDto struct {
	Id        int64     `json:"id"`
	PublicKey string    `json:"publicKey"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}
