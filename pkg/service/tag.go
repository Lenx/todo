package service

import (
	"github.com/Hanqur/todo_app"
	"github.com/Hanqur/todo_app/pkg/repository"
)

type TagService struct {
	repo     repository.Tag
	itemRepo repository.TodoItem
}

func NewTagService(repo repository.Tag, itemRepo repository.TodoItem) *TagService {
	return &TagService{repo: repo, itemRepo: itemRepo}
}

func (s *TagService) CreateTag(userId int, itemId int, input todo_app.Tag) (int, error) {
	_, err := s.itemRepo.GetItemById(userId, itemId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateTag(itemId, input)
}

func (s *TagService) GetAllTags(userId int, itemId int) ([]todo_app.Tag, error) {
	_, err := s.itemRepo.GetItemById(userId, itemId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllTags(userId, itemId)
}