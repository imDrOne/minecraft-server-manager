package dto

import "time"

type NodeDto struct {
	Id        int64     `json:"id"`
	Host      string    `json:"host"`
	Port      int32     `json:"port"`
	CreatedAt time.Time `json:"createdAt"`
}
