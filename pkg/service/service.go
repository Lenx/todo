package service

import (
	"github.com/Hanqur/todo_app"
	"github.com/Hanqur/todo_app/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todo_app.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	CreateList(userId int, list todo_app.TodoList) (int, error)
	GetAllLists(userId int) ([]todo_app.TodoList, error)
	GetListById(userId int, listId int) (todo_app.TodoList, error)
	DeleteList(userId int, listId int) error
	UpdateList(userId int, listId int, input todo_app.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(userId int, listId int, input todo_app.Item) (int, error)
	GetAllItems(userId int, listId int) ([]todo_app.Item, error)
	GetItemById(userId int, itemId int) (todo_app.Item, error)
	UpdateItem(userId int, itemId int, input todo_app.UpdateItemInput) error
	DeleteItem(userId int, itemId int) error
}

type Tag interface {
	CreateTag(userId int, itemId int, input todo_app.Tag) (int, error)
	GetAllTags(userId int, itemId int) ([]todo_app.Tag, error)
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	Tag
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList: NewTodoListService(repos.TodoList),
		TodoItem: NewTodoItemService(repos.TodoItem, repos.TodoList),
		Tag: NewTagService(repos.Tag, repos.TodoItem),
	}
}