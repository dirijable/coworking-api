package service

import (
	"context"
	"fmt"

	"github.com/dirijable/coworking-api/internal/domain"
	"github.com/dirijable/coworking-api/internal/errorsx/service"
	"github.com/dirijable/coworking-api/internal/model"
	"github.com/dirijable/coworking-api/internal/service/mapper"
	"github.com/dirijable/coworking-api/internal/service/validator"
	"github.com/dirijable/coworking-api/pkg/postgres/transactor"
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
	repo      Repository
	txManager transactor.Transactor
}

func NewService(repo Repository, txManager transactor.Transactor) *RoomService {
	return &RoomService{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *RoomService) Create(ctx context.Context, dRoom domain.Room) (domain.Room, error) {
	if err := validator.Validate(dRoom); err != nil {
		return domain.Room{}, fmt.Errorf("validate room: %w", err)
	}
	mRoom := mapper.DomainToModel(dRoom)
	var createdDomainRoom domain.Room
	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exist, err := s.repo.ExistByName(txCtx, mRoom)
		if err != nil {
			return fmt.Errorf("conflict check: %w", err)
		}
		if exist {
			return service.ErrConflict
		}
		createdModelRoom, err := s.repo.Create(txCtx, mRoom)
		if err != nil {
			return fmt.Errorf("create room: %w", err)
		}
		createdDomainRoom = mapper.ModelToDomain(createdModelRoom)
		return nil
	})
	if err != nil {
		return domain.Room{}, err
	}
	return createdDomainRoom, nil
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
