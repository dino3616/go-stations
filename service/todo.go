package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mattn/go-sqlite3"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	todo := &model.TODO{
		ID: id,
	}
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	return nil, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}
	updates, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if updates == 0 {
		err := &model.ErrNotFound{
			Sqlite3Error: sqlite3.Error{
				Code:         sqlite3.ErrConstraint,
				ExtendedCode: sqlite3.ErrConstraintPrimaryKey,
			},
		}
		return nil, err.Unwrap()
	}

	id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	todo := &model.TODO{
		ID: id,
	}
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
