package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

	var rows *sql.Rows
	var err error
	if size == 0 {
		size = -1
	}
	if prevID == 0 {
		rows, err = s.db.QueryContext(ctx, read, size)
	} else {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := make([]*model.TODO, 0)
	for rows.Next() {
		todo := &model.TODO{}
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return todos, nil
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
	if len(ids) == 0 {
		return nil
	}

	const deleteFmt = `DELETE FROM todos WHERE id IN (%s)`

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	res, err := s.db.ExecContext(ctx, fmt.Sprintf(deleteFmt, strings.TrimRight(strings.Repeat("?,", len(ids)), ",")), args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return &model.ErrNotFound{
			Sqlite3Error: sqlite3.Error{
				Code:         sqlite3.ErrConstraint,
				ExtendedCode: sqlite3.ErrConstraintUnique,
			},
		}
	}

	return nil
}
