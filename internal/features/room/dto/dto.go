package dto

import "github.com/google/uuid"

type RoomRequestDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Capacity int    `json:"capacity" validate:"required,gte=1,lte=1000"`
}

type RoomResponseDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Capacity int       `json:"capacity"`
}
