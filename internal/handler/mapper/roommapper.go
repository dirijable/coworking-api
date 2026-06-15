package mapper

import (
	"github.com/dirijable/coworking-api/internal/domain"
	"github.com/dirijable/coworking-api/internal/dto"
)

func RequestToDomain(req dto.RoomRequestDTO) domain.Room {
	return domain.Room{
		Name:     req.Name,
		Capacity: req.Capacity,
	}
}

func DomainToResponse(dom domain.Room) dto.RoomResponseDTO {
	return dto.RoomResponseDTO{
		ID:       dom.ID,
		Name:     dom.Name,
		Capacity: dom.Capacity,
	}
}
