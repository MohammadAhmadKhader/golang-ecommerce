package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mohammadahmadkhader/golang-ecommerce/types"
)

func TestUserServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)
	
	t.Run("Should return 400 status code if payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "john",
			LastName:  "doe",
			Email:     "invalid@gmail.com",
			Password:  "123",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("Should correctly regist the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "john",
			LastName:  "doe",
			Email:     "valid@gmail.com",
			Password:  "123567534",
		}
	
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost,"/register",bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(recorder, req)

		if recorder.Code != 200 {
			t.Errorf("expected status code to be 200 received %d", recorder.Code)
		}
	})
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user was not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}
