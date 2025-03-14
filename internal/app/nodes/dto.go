package nodes

import (
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"time"
)

// CreateNodeRequestDto todo: add validation
type CreateNodeRequestDto struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

type UpdateNodeRequestDto = CreateNodeRequestDto

type FindNodeRequestDto = pagination.PagePaginationRequestDto

type NodeResponseDto struct {
	Id        int64     `json:"id"`
	Host      string    `json:"host"`
	Port      int32     `json:"port"`
	CreatedAt time.Time `json:"createdAt"`
}
