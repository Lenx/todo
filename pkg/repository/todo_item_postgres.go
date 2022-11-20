package repository

import (
	"fmt"
	"strings"

	"github.com/Hanqur/todo_app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) CreateItem(listId int, input todo_app.Item) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description, deadline) values ($1, $2, $3) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, input.Title, input.Description, input.Deadline)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItems(userId int, listId int) ([]todo_app.Item, error) {
	var items []todo_app.Item
	getAllItemsQuery := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description, tl.done, tl.deadline FROM %s tl 
						INNER JOIN %s li on li.item_id = tl.id
						INNER JOIN %s ul on ul.list_id = li.list_id 
						WHERE li.list_id=$1 AND ul.user_id=$2`, 
						todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, getAllItemsQuery, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetItemById(userId int, itemId int) (todo_app.Item, error) {
	var item todo_app.Item
	getItemsByIdQuery := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description, tl.done, tl.deadline FROM %s tl 
									INNER JOIN %s li on li.item_id = tl.id
									INNER JOIN %s ul on ul.list_id = li.list_id 
									WHERE tl.id=$1 AND ul.user_id=$2`, 
						todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, getItemsByIdQuery, itemId, userId); err != nil {
		return item, err 
	}

	return item, nil
}

func (r *TodoItemPostgres) UpdateItem(userId int, itemId int, input todo_app.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId ++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId ++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId ++
	}

	if input.Deadline.String() != "2001-01-01 00:00:00 +0000 UTC" {
		setValues = append(setValues, fmt.Sprintf("deadline = $%d", argId))
		args = append(args, input.Deadline)
		argId ++
	}


	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s tl SET %s FROM %s li, %s ul
						WHERE tl.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND tl.id = $%d`, 
						todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoItemPostgres) DeleteItem(userId int, itemId int) error {
	deleteItemQuery := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
									todoItemsTable, listsItemsTable, usersListsTable)
	if _, err := r.db.Exec(deleteItemQuery, userId, itemId); err != nil {
		return err
	}
	
	return nil
}