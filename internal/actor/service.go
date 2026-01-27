package actor

import (
	"context"
	"errors"
)

type Service interface {
	GetAll(ctx context.Context, page, limit int, search string) ([]Actor, int64, error)
	GetByID(ctx context.Context, id string) (*Actor, error)
	Create(ctx context.Context, req CreateActorRequest) (*Actor, error)
	Update(ctx context.Context, id string, req UpdateActorRequest) (*Actor, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll(ctx context.Context, page, limit int, search string) ([]Actor, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(ctx, limit, offset, search)
}

func (s *service) GetByID(ctx context.Context, id string) (*Actor, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, req CreateActorRequest) (*Actor, error) {
	actor := &Actor{
		Name:     req.Name,
		PhotoURL: req.PhotoURL,
	}

	if err := s.repo.Create(ctx, actor); err != nil {
		return nil, err
	}

	return actor, nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateActorRequest) (*Actor, error) {
	actor, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if actor == nil {
		return nil, errors.New("actor not found")
	}

	actor.Name = req.Name
	actor.PhotoURL = req.PhotoURL

	if err := s.repo.Update(ctx, actor); err != nil {
		return nil, err
	}

	return actor, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
