package goquery

import (
	"net/http/httptest"
	"testing"
)

func TestWeb(t *testing.T) {
	db := initDB()
	service := BuildController[UserEntity, *UserQuery]("/user/", db, UserEntity{}, &UserQuery{})

	t.Run("Get /user/", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?page=1&size=2", nil)

		service.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"Data":[{"Id":1,"Score":85,"Memo":"Good"},{"Id":2,"Score":40,"Memo":"Bad"}],"Total":4}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
