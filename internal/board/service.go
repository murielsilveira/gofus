package board

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/murielsilveira/gofus/internal/db/sqlc"
	"github.com/murielsilveira/gofus/internal/platform/errs"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) Create(ctx context.Context, name string) (sqlc.Board, error) {
	return s.q.CreateBoard(ctx, name)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (sqlc.Board, error) {
	board, err := s.q.GetBoard(ctx, id)
	return board, mapError(err)
}

func (s *Service) List(ctx context.Context) ([]sqlc.Board, error) {
	return s.q.ListBoards(ctx)
}

type UpdateInput struct {
	Name *string
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (sqlc.Board, error) {
	existing, err := s.Get(ctx, id)
	if err != nil {
		return sqlc.Board{}, err
	}

	name := existing.Name
	if in.Name != nil {
		name = *in.Name
	}

	board, err := s.q.UpdateBoard(ctx, sqlc.UpdateBoardParams{
		ID:   id,
		Name: name,
	})
	return board, mapError(err)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	rows, err := s.q.DeleteBoard(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return errs.ErrNotFound
	}
	return err
}
