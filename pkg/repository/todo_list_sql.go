package repository

import (
	"fmt"
	"todo"

	"github.com/jmoiron/sqlx"
)

type TodoListSql struct {
	db *sqlx.DB
}

func NewTodoListSql(db *sqlx.DB) *TodoListSql {
	return &TodoListSql{db: db}
}

func (r *TodoListSql) Create(userId int, list todo.TodoList)(int, error)  {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row :=tx.QueryRow(createListQuery, list.Title, list.Description)
	if err:= row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}