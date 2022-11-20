package todo_app

type User struct {
	Id int `form:"-" json:"-" db:"id"`
	Name string `form:"name" json:"name" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}