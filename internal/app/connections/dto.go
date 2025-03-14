package connections

import "time"

type CreateConnectionRequestDto struct {
	NodeId int64  `json:"nodeId"`
	Key    string `json:"key"`
	User   string `json:"user"`
}

type UpdateConnectionRequestDto struct {
	Key  string `json:"key"`
	User string `json:"user"`
}

type ConnectionResponseDto struct {
	Id        int64     `json:"id"`
	Key       string    `json:"key"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}
