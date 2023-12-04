package goquery

import (
	log "github.com/sirupsen/logrus"
	"net/http/httptest"
	"testing"
)

func TestWeb(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	db := initDB()
	service := BuildController[UserEntity, *UserQuery](
		"/user/", db,
		func() UserEntity { return UserEntity{} },
		func() *UserQuery { return &UserQuery{} },
	)

	t.Run("Get /user/", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?PageNumber=1&PageSize=2", nil)

		service.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"Data":[{"Id":1,"Score":85,"Memo":"Good"},{"Id":2,"Score":40,"Memo":"Bad"}],"Total":4}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?ScoreLt=60&test=test", nil)

		service.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"Data":[{"Id":2,"Score":40,"Memo":"Bad"},{"Id":3,"Score":55,"Memo":"Bad"}],"Total":2}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
