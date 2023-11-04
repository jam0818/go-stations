package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

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
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	// トランザクションを開始
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// TODOを挿入
	result, err := tx.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	// 挿入したTODOのIDを取得
	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	// 挿入したTODOの情報をクエリ
	row := tx.QueryRowContext(ctx, confirm, id)
	var todo model.TODO
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	todo.ID = id

	if err != nil {
		return nil, err
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var rows *sql.Rows
	var err error

	if prevID != 0 {
		rows, err = s.db.QueryContext(ctx, readWithID, prevID, size)
	} else {
		rows, err = s.db.QueryContext(ctx, read, size)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.TODO
	for rows.Next() {
		var todo model.TODO
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if len(todos) == 0 {
		return []*model.TODO{}, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		return nil, err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if numRows == 0 {
		return nil, &model.ErrNotFound{}
	}

	row := tx.QueryRowContext(ctx, confirm, id)
	var todo model.TODO
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	todo.ID = id
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	var query string
	var args []interface{}

	if len(ids) == 0 {
		// ids が空の場合はクエリを実行しない
		return nil
	} else if len(ids) == 1 {
		// ids が 1 つの要素を持つ場合は単一のプレースホルダーを使用
		query = "DELETE FROM todos WHERE id = ?"
		args = append(args, ids[0])
	} else {
		// ids が複数の要素を持つ場合は必要な数のプレースホルダーを生成
		query = fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1))
		for _, id := range ids {
			args = append(args, id)
		}
	}
	// ids を []interface{} に変換
	idInterfaces := make([]interface{}, len(ids))
	log.Println(query)
	for i, id := range ids {
		idInterfaces[i] = id
	}

	// トランザクションを開始
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TODO 削除クエリの実行
	result, err := tx.ExecContext(ctx, query, idInterfaces...)
	if err != nil {
		return err
	}

	// 削除された行数を取得
	numRows, err := result.RowsAffected()
	log.Printf("numrows: %v", numRows)
	if err != nil {
		return err
	}

	// 削除された行がなかった場合はエラーを返す
	if numRows == 0 {
		return &model.ErrNotFound{}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
