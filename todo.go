package todo_app

import (
	"errors"
	"fmt"
	"time"
)

type TodoList struct {
	Id          int    `form:"-" json:"id" db:"id"`
	Title       string `form:"title" json:"title" db:"title" binding:"required"`
	Description string `form:"description" json:"description" db:"description"`
}

type UserList struct {
	Id     int
	UserId int
	ListId int
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateListInput struct {
	Title       *string `form:"title" json:"title" db:"title" binding:"required"`
	Description *string `form:"description" json:"description" db:"description"`
}

func (i UpdateListInput) Validate() error {
	fmt.Println(i.Title)
	if *i.Title == "" && *i.Description == "" {
		return errors.New("update structure has no values")
	}

	return nil
}

type Item struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done bool `json:"done" db:"done"`
	Deadline time.Time `json:"deadline" db:"deadline"`
}

type UpdateItemInput struct {
	Title       *string `json:"title" binding:"required"`
	Description *string `json:"description"`
	Done *bool `json:"done"`
	Deadline time.Time `json:"deadline" db:"deadline"`
}

func (i UpdateItemInput) Validate() error {
	fmt.Println(*i.Title, *i.Description)
	if *i.Title == "" && *i.Description == "" && i.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type Tag struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" binding:"required"`
}
