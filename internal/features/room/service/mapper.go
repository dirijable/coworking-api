package service

import (
	"github.com/dirijable/coworking-api/internal/features/room/domain"
	"github.com/dirijable/coworking-api/internal/features/room/model"
)

func DomainToModel(room domain.Room) model.Room {
	return model.Room{
		Name:     room.Name,
		Capacity: room.Capacity,
	}
}

func ModelToDomain(room model.Room) domain.Room {
	return domain.Room{
		ID:       room.ID,
		Name:     room.Name,
		Capacity: room.Capacity,
	}
}
