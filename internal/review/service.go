package review

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, userID string, req CreateReviewRequest) (*Review, error)
	GetByDrama(ctx context.Context, dramaID string, page, limit int) ([]Review, int64, error)
	Update(ctx context.Context, userID, reviewID string, req UpdateReviewRequest) (*Review, error)
	Delete(ctx context.Context, userID, reviewID string, isAdmin bool) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID string, req CreateReviewRequest) (*Review, error) {
	// Check if already reviewed
	existing, err := s.repo.GetByUserAndDrama(ctx, userID, req.DramaID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("review already exists")
	}

	review := &Review{
		UserID:     userID,
		DramaID:    req.DramaID,
		Rating:     req.Rating,
		ReviewText: req.ReviewText,
	}

	if err := s.repo.Create(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *service) GetByDrama(ctx context.Context, dramaID string, page, limit int) ([]Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetByDramaID(ctx, dramaID, limit, offset)
}

func (s *service) Update(ctx context.Context, userID, reviewID string, req UpdateReviewRequest) (*Review, error) {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}
	if review == nil {
		return nil, errors.New("review not found")
	}

	if review.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	review.Rating = req.Rating
	review.ReviewText = req.ReviewText

	if err := s.repo.Update(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *service) Delete(ctx context.Context, userID, reviewID string, isAdmin bool) error {
	review, err := s.repo.FindByID(ctx, reviewID)
	if err != nil {
		return err
	}
	if review == nil {
		return errors.New("review not found")
	}

	if !isAdmin && review.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.Delete(ctx, reviewID)
}
