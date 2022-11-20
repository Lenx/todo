package repository

import (
	"github.com/Hanqur/todo_app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo_app.User) (int, error)
	GetUser(username, password string) (todo_app.User, error)
}

type TodoList interface {
	CreateList(userId int, list todo_app.TodoList) (int, error)
	GetAllLists(userId int) ([]todo_app.TodoList, error)
	GetListById(userId int, listId int) (todo_app.TodoList, error)
	DeleteList(userId int, listId int) error
	UpdateList(userId int, listId int, input todo_app.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(listId int, input todo_app.Item) (int, error)
	GetAllItems(userId int, listId int) ([]todo_app.Item, error)
	GetItemById(userId int, itemId int) (todo_app.Item, error)
	UpdateItem(userId int, itemId int, input todo_app.UpdateItemInput) error
	DeleteItem(userId int, itemId int) error
}

type Tag interface {
	CreateTag(itemId int, input todo_app.Tag) (int, error)
	GetAllTags(userId int, itemId int) ([]todo_app.Tag, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
	Tag
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList: NewTodoListPostgres(db),
		TodoItem: NewTodoItemPostgres(db),
		Tag: NewTagPostgres(db),
	}
}