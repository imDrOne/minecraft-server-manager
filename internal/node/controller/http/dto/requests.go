package dto

import "github.com/imDrOne/minecraft-server-manager/pkg/pagination/dto"

type CreateNodeDto struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

type UpdateNodeDto = CreateNodeDto

type FindNodeDto = dto.PagePaginationDto
