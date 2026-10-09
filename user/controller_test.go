package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func newTestController(t *testing.T) (*Controller, *sql.DB) {
	t.Helper()
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "users.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	schema, err := os.ReadFile("../sql/create-database.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatal(err)
	}
	return NewController(NewServiceFacade(NewService(NewUserRepo(db)))), db
}

func TestCreateUserJSONFieldCasing(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int
	}{
		{"lowercase", `{"email":"alice@example.com","name":"Alice","password":"test-password"}`, http.StatusOK},
		{"empty password", `{"email":"alice@example.com","name":"Alice","password":""}`, http.StatusOK},
		{"capital email", `{"Email":"alice@example.com","name":"Alice","password":"test-password"}`, http.StatusBadRequest},
		{"uppercase email", `{"EMAIL":"alice@example.com","name":"Alice","password":"test-password"}`, http.StatusBadRequest},
		{"mixed email", `{"eMaIl":"alice@example.com","name":"Alice","password":"test-password"}`, http.StatusBadRequest},
		{"capital name", `{"email":"alice@example.com","Name":"Alice","password":"test-password"}`, http.StatusBadRequest},
		{"mixed name", `{"email":"alice@example.com","nAmE":"Alice","password":"test-password"}`, http.StatusBadRequest},
		{"capital password", `{"email":"alice@example.com","name":"Alice","Password":"test-password"}`, http.StatusBadRequest},
		{"mixed password", `{"email":"alice@example.com","name":"Alice","pAsSwOrD":"test-password"}`, http.StatusBadRequest},
		{"Unicode case variant", `{"email":"alice@example.com","name":"Alice","paſſword":"test-password"}`, http.StatusBadRequest},
		{"escaped capital key", `{"email":"alice@example.com","name":"Alice","\u0050assword":"test-password"}`, http.StatusBadRequest},
		{"duplicate casing", `{"email":"alice@example.com","Email":"bob@example.com","name":"Alice","password":"test-password"}`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller, db := newTestController(t)
			router := gin.New()
			router.POST("/user", controller.CreateUser)
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			if resp.Code != tc.want {
				t.Fatalf("status = %d, want %d; body = %s", resp.Code, tc.want, resp.Body.String())
			}
			var count int
			if err := db.QueryRow("select count(*) from users").Scan(&count); err != nil {
				t.Fatal(err)
			}
			wantCount := 0
			if tc.want == http.StatusOK {
				wantCount = 1
				if resp.Body.Len() != 0 {
					t.Fatalf("registration body = %s, want empty", resp.Body.String())
				}
			}
			if count != wantCount {
				t.Fatalf("created users = %d, want %d", count, wantCount)
			}
		})
	}
}

func TestGetUserJSONFields(t *testing.T) {
	for _, key := range []string{"email", "Email", "EMAIL"} {
		t.Run(key, func(t *testing.T) {
			controller, db := newTestController(t)
			if _, err := db.Exec("insert into users(id, email, name, password_hash) values(?, ?, ?, ?)",
				"user-id", "alice@example.com", "Alice", []byte("test-hash")); err != nil {
				t.Fatal(err)
			}
			router := gin.New()
			router.POST("/user/getbyemail", controller.GetUserByEmail)
			req := httptest.NewRequest(http.MethodPost, "/user/getbyemail", strings.NewReader(`{"`+key+`":"alice@example.com"}`))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			if resp.Code != http.StatusOK {
				t.Fatalf("status = %d; body = %s", resp.Code, resp.Body.String())
			}
			var got map[string]any
			if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			want := map[string]any{"user": map[string]any{
				"id": "user-id", "name": "Alice", "email": "alice@example.com",
			}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("response = %#v, want %#v", got, want)
			}
		})
	}
}
