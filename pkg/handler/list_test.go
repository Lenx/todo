package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Hanqur/todo_app"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lenx/todo"
	"github.com/lenx/todo/pkg/service"
	mock_service "github.com/lenx/todo/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandlerList_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId int, list todo.TodoList)

	testTable := []struct {
		name                string
		userId              int
		inputBody           string
		list                todo.TodoList
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			userId:    1,
			inputBody: `{"title": "Test", "description": "Test"}`,
			list: todo.TodoList{
				Title:       "Test",
				Description: "Test",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId int, list todo.TodoList) {
				s.EXPECT().CreateList(userId, list).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Input without title",
			userId:    1,
			inputBody: `{"title": "", "description": "Test"}`,
			list: todo.TodoList{
				Title:       "",
				Description: "Test",
			},
			mockBehavior:        func(s *mock_service.MockTodoList, userId int, list todo_app.TodoList) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"Message":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			userId:    1,
			inputBody: `{"title": "Test", "description": "Test"}`,
			list: todo.TodoList{
				Title:       "Test",
				Description: "Test",
			},
			mockBehavior: func(s *mock_service.MockTodoList, userId int, list todo.TodoList) {
				s.EXPECT().CreateList(userId, list).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"Message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.userId, testCase.list)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/api/lists/", func(c *gin.Context) {
				c.Set(userCtx, testCase.userId)
			}, handler.createList)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists/", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())

		})
	}
}

func TestHandlerList_getAllLists(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, userId int)
	var response = []todo.TodoList{

		{
			Id:          1,
			Title:       "Test",
			Description: "Test",
		},
	}

	testTable := []struct {
		name                string
		userId              int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "OK",
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoList, userId int) {
				s.EXPECT().GetAllLists(userId).Return(response, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"title":"Test","description":"Test"}]}`,
		},
		{
			name:   "Service failure",
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoList, userId int) {
				s.EXPECT().GetAllLists(userId).Return(response, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"Message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			list := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(list, testCase.userId)

			services := &service.Service{TodoList: list}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/api/lists/", func(c *gin.Context) {
				c.Set(userCtx, testCase.userId)
			}, handler.getAllLists)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/lists/", bytes.NewBufferString(""))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())

		})
	}

}
