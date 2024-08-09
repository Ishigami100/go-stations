package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
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
		confirm = `SELECT id,subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the insert statement
	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	// Get the inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Prepare the confirm statement
	stmtConfirm, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}
	defer stmtConfirm.Close()

	// Query the inserted record
	row := stmtConfirm.QueryRowContext(ctx, id)

	// Scan the result into a TODO object
	todo := &model.TODO{}
	err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
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
		confirm = `SELECT id,subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the update statement
	result, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, &model.ErrNotFound{}
	}

	// Prepare the confirm statement
	stmtConfirm, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}
	defer stmtConfirm.Close()

	// Query the inserted record
	row := stmtConfirm.QueryRowContext(ctx, id)

	// Scan the result into a TODO object
	todo := &model.TODO{}
	err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
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
