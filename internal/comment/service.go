package comment

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, userID string, req CreateCommentRequest) (*Comment, error)
	GetByEpisode(ctx context.Context, episodeID string, page, limit int) ([]Comment, int64, error)
	Update(ctx context.Context, userID, commentID string, req UpdateCommentRequest) (*Comment, error)
	Delete(ctx context.Context, userID, commentID string, isAdmin bool) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID string, req CreateCommentRequest) (*Comment, error) {
	comment := &Comment{
		UserID:      userID,
		EpisodeID:   req.EpisodeID,
		CommentText: req.CommentText,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *service) GetByEpisode(ctx context.Context, episodeID string, page, limit int) ([]Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetByEpisodeID(ctx, episodeID, limit, offset)
}

func (s *service) Update(ctx context.Context, userID, commentID string, req UpdateCommentRequest) (*Comment, error) {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("comment not found")
	}
	if comment.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	comment.CommentText = req.CommentText
	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *service) Delete(ctx context.Context, userID, commentID string, isAdmin bool) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("comment not found")
	}
	if !isAdmin && comment.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(ctx, commentID)
}
