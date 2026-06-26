package handler

import (
	"context"
	"fmt"
	"net/http"

	httpreq "github.com/dirijable/coworking-api/internal/core/http/request"
	httpresp "github.com/dirijable/coworking-api/internal/core/http/response"
	"github.com/dirijable/coworking-api/internal/features/room/domain"
	"github.com/dirijable/coworking-api/internal/features/room/dto"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, room domain.Room) (domain.Room, error)
	FindById(ctx context.Context, id uuid.UUID) (domain.Room, error)
	FindAll(ctx context.Context) ([]domain.Room, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type RoomHTTPHandler struct {
	srv Service
}

func NewRoomHTTPHandler(service Service) *RoomHTTPHandler {
	return &RoomHTTPHandler{
		srv: service,
	}
}

func (h *RoomHTTPHandler) Create(rw http.ResponseWriter, r *http.Request) {
	var reqDTO dto.RoomRequestDTO
	if err := httpreq.ValidateAndDecode(r, &reqDTO); err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}
	domainRoom := RequestToDomain(reqDTO)
	createdRoom, err := h.srv.Create(r.Context(), domainRoom)
	if err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}

	locationPath := fmt.Sprintf("%s/%s", r.URL.Path, createdRoom.ID.String())
	rw.Header().Set("Location", locationPath)
	if err = httpresp.SendJSONResponse(rw, http.StatusCreated, DomainToResponse(createdRoom)); err != nil {
		return
	}
}

func (h *RoomHTTPHandler) FindById(rw http.ResponseWriter, r *http.Request) {
	id, err := pathUUID(r, "id")
	if err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}
	domainRoom, err := h.srv.FindById(r.Context(), id)
	if err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}
	if err = httpresp.SendJSONResponse(rw, http.StatusOK, DomainToResponse(domainRoom)); err != nil {
		return
	}
}

func (h *RoomHTTPHandler) FindAll(rw http.ResponseWriter, r *http.Request) {
	domainRooms, err := h.srv.FindAll(r.Context())
	if err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}
	dtoRooms := make([]dto.RoomResponseDTO, 0, len(domainRooms))
	for _, dr := range domainRooms {
		dtoRooms = append(dtoRooms, DomainToResponse(dr))
	}
	if err = httpresp.SendJSONResponse(rw, http.StatusOK, dtoRooms); err != nil {
		return
	}
}

func (h *RoomHTTPHandler) DeleteById(rw http.ResponseWriter, r *http.Request) {
	id, err := pathUUID(r, "id")
	if err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}
	if err := h.srv.DeleteById(r.Context(), id); err != nil {
		httpresp.SendErrorResponse(rw, err)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
