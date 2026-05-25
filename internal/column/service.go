package column

import (
	"context"
	"errors"

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

type CreateInput struct {
	BoardID  int32
	Name     string
	Position int32
}

func (s *Service) Create(ctx context.Context, in CreateInput) (sqlc.Column, error) {
	if _, err := s.q.GetBoard(ctx, in.BoardID); err != nil {
		return sqlc.Column{}, mapError(err)
	}

	return s.q.CreateColumn(ctx, sqlc.CreateColumnParams{
		BoardID:  in.BoardID,
		Name:     in.Name,
		Position: in.Position,
	})
}

func (s *Service) Get(ctx context.Context, columnID int32) (sqlc.Column, error) {
	column, err := s.q.GetColumn(ctx, columnID)
	return column, mapError(err)
}

func (s *Service) ListByBoard(ctx context.Context, boardID int32) ([]sqlc.Column, error) {
	if _, err := s.q.GetBoard(ctx, boardID); err != nil {
		return nil, mapError(err)
	}

	return s.q.ListColumnsByBoard(ctx, boardID)
}

type UpdateInput struct {
	Name     *string
	Position *int32
}

func (s *Service) Update(ctx context.Context, columnID int32, in UpdateInput) (sqlc.Column, error) {
	existing, err := s.Get(ctx, columnID)
	if err != nil {
		return sqlc.Column{}, err
	}

	name := existing.Name
	if in.Name != nil {
		name = *in.Name
	}

	position := existing.Position
	if in.Position != nil {
		position = *in.Position
	}

	column, err := s.q.UpdateColumn(ctx, sqlc.UpdateColumnParams{
		ID:       columnID,
		Name:     name,
		Position: position,
	})
	return column, mapError(err)
}

func (s *Service) Delete(ctx context.Context, columnID int32) error {
	rows, err := s.q.DeleteColumn(ctx, columnID)
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
