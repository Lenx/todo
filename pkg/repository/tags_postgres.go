package repository

import (
	"fmt"

	"github.com/Hanqur/todo_app"
	"github.com/jmoiron/sqlx"
)

type TagPostgres struct {
	db *sqlx.DB
}

func NewTagPostgres(db *sqlx.DB) *TagPostgres {
	return &TagPostgres{db: db}
}

func (r *TagPostgres) CreateTag(itemId int, input todo_app.Tag) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var tagId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title) values ($1) RETURNING id", tagsTable)
	row := tx.QueryRow(createItemQuery, input.Title)
	err = row.Scan(&tagId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (item_id, tag_id) values ($1, $2)", itemTagsTable)
	_, err = tx.Exec(createListItemQuery, itemId, tagId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return tagId, tx.Commit()
}

func (r *TagPostgres) GetAllTags(userId int, itemId int) ([]todo_app.Tag, error) {
	var tags []todo_app.Tag
	getAllItemsQuery := fmt.Sprintf(`SELECT tag.id, tag.title FROM %s tag 
									INNER JOIN %s tags_item on tags_item.tag_id = tag.id
									INNER JOIN %s items_lists on items_lists.item_id = tags_item.item_id
									INNER JOIN %s lists_users on lists_users.list_id = items_lists.list_id
									WHERE tags_item.item_id=$1 AND lists_users.user_id=$2`, 
						tagsTable, itemTagsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&tags, getAllItemsQuery, itemId, userId); err != nil {
		return nil, err
	}

	return tags, nil
}