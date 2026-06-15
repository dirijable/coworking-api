package service

import (
	"context"
	"fmt"

	"github.com/dirijable/coworking-api/internal/features/domain"
	"github.com/dirijable/coworking-api/internal/features/model"
	service2 "github.com/dirijable/coworking-api/internal/features/room/service"
	"github.com/dirijable/coworking-api/internal/features/service/errorx"
	"github.com/dirijable/coworking-api/internal/features/service/mapper"
	"github.com/google/uuid"
)

type Repository interface {
	ExistByName(ctx context.Context, room model.Room) (bool, error)
	Create(ctx context.Context, room model.Room) (model.Room, error)
	FindById(ctx context.Context, id uuid.UUID) (model.Room, error)
	FindAll(ctx context.Context) ([]model.Room, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type RoomService struct {
	repo Repository
}

func NewService(repo Repository) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (s *RoomService) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	if err := service2.validate(room); err != nil {
		return domain.Room{}, fmt.Errorf("validate room: %w", err)
	}
	dRoom := mapper.DomainToModel(room)
	exist, err := s.repo.ExistByName(ctx, dRoom)
	if err != nil {
		return domain.Room{}, fmt.Errorf("conflict check: %w", err)
	}
	if exist {
		return domain.Room{}, errorx.ErrConflict
	}
	createdRoom, err := s.repo.Create(ctx, dRoom)
	if err != nil {
		return domain.Room{}, fmt.Errorf("create room: %w", err)
	}
	return mapper.ModelToDomain(createdRoom), nil
}

func (s *RoomService) FindById(ctx context.Context, id uuid.UUID) (domain.Room, error) {
	room, err := s.repo.FindById(ctx, id)
	if err != nil {
		return domain.Room{}, fmt.Errorf("find by id: %w", err)
	}
	return mapper.ModelToDomain(room), nil
}

func (s *RoomService) FindAll(ctx context.Context) ([]domain.Room, error) {
	mRooms, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("find all: %w", err)
	}
	dRooms := make([]domain.Room, 0, len(mRooms))
	for _, room := range mRooms {
		dRooms = append(dRooms, mapper.ModelToDomain(room))
	}
	return dRooms, nil
}

func (s *RoomService) DeleteById(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteById(ctx, id); err != nil {
		return fmt.Errorf("delete by id: %w", err)
	}
	return nil
}
