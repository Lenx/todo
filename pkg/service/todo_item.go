package service

import (
	"fmt"

	"github.com/Hanqur/todo_app"
	"github.com/Hanqur/todo_app/pkg/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) CreateItem(userId int, listId int, input todo_app.Item) (int, error) {
	_, err := s.listRepo.GetListById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateItem(listId, input)
}

func (s *TodoItemService) GetAllItems(userId int, listId int) ([]todo_app.Item, error) {
	_, err := s.listRepo.GetListById(userId, listId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllItems(userId, listId)
}

func (s *TodoItemService) GetItemById(userId int, itemId int) (todo_app.Item, error) {
	return s.repo.GetItemById(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId int, itemId int, input todo_app.UpdateItemInput) error {
	fmt.Println("Hello from service")
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateItem(userId, itemId, input)
}

func (s *TodoItemService) DeleteItem(userId int, itemId int) error {
	return s.repo.DeleteItem(userId, itemId)
}