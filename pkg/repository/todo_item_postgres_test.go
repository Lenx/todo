package repository

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/Hanqur/todo_app"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoItemPostgres_Create(t *testing.T){
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item todo_app.Item
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name 			string
		mockBehavior 	mockBehavior
		args 			args
		id 				int
		wantErr 		bool
	}{
		{
			name: "OK",
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description, args.item.Deadline).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1,1))

				mock.ExpectCommit()
			},
			args: args{
				listId: 1,
				item: todo_app.Item{
					Title: "Test",
					Description: "Test",
					Deadline: time.Date(2022, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
			},
			id: 2,
			wantErr: false,
		},
		{
			name: "Invalid input",
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description, args.item.Deadline).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args: args{
				listId: 1,
				item: todo_app.Item{
					Title: "",
					Description: "Test",
					Deadline: time.Date(2022, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
			},
			wantErr: true,
		},
		{
			name: "Second insert error",
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description, args.item.Deadline).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			args: args{
				listId: 1,
				item: todo_app.Item{
					Title: "Test",
					Description: "Test",
					Deadline: time.Date(2022, time.November, 10, 23, 0, 0, 0, time.UTC),
				},
			},
			id: 2,
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.CreateItem(testCase.args.listId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}

}