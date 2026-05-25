package task

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
	ColumnID    int32
	Title       string
	Description string
	Position    int32
}

func (s *Service) Create(ctx context.Context, in CreateInput) (sqlc.Task, error) {
	if _, err := s.q.GetColumn(ctx, in.ColumnID); err != nil {
		return sqlc.Task{}, mapError(err)
	}

	return s.q.CreateTask(ctx, sqlc.CreateTaskParams{
		ColumnID:    in.ColumnID,
		Title:       in.Title,
		Description: in.Description,
		Position:    in.Position,
	})
}

func (s *Service) Get(ctx context.Context, taskID int32) (sqlc.Task, error) {
	task, err := s.q.GetTask(ctx, taskID)
	return task, mapError(err)
}

func (s *Service) ListByColumn(ctx context.Context, columnID int32) ([]sqlc.Task, error) {
	if _, err := s.q.GetColumn(ctx, columnID); err != nil {
		return nil, mapError(err)
	}

	return s.q.ListTasksByColumn(ctx, columnID)
}

type UpdateInput struct {
	ColumnID    *int32
	Title       *string
	Description *string
	Position    *int32
}

func (s *Service) Update(ctx context.Context, taskID int32, in UpdateInput) (sqlc.Task, error) {
	existing, err := s.Get(ctx, taskID)
	if err != nil {
		return sqlc.Task{}, err
	}

	columnID := existing.ColumnID
	if in.ColumnID != nil {
		if _, err := s.q.GetColumn(ctx, *in.ColumnID); err != nil {
			return sqlc.Task{}, mapError(err)
		}
		columnID = *in.ColumnID
	}

	title := existing.Title
	if in.Title != nil {
		title = *in.Title
	}

	description := existing.Description
	if in.Description != nil {
		description = *in.Description
	}

	position := existing.Position
	if in.Position != nil {
		position = *in.Position
	}

	task, err := s.q.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:          taskID,
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Position:    position,
	})
	return task, mapError(err)
}

func (s *Service) Delete(ctx context.Context, taskID int32) error {
	rows, err := s.q.DeleteTask(ctx, taskID)
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
