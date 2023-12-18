package web

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/rdb"
	. "github.com/doytowin/goquery/test"
	log "github.com/sirupsen/logrus"
	"net/http/httptest"
	"testing"
)

func TestWeb(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	db := rdb.Connect("app.properties")
	InitDB(db)
	defer rdb.Disconnect(db)

	ctx := context.Background()
	tm := rdb.NewTransactionManager(db)

	createUserEntity := func() UserEntity { return UserEntity{} }
	userDataAccess := rdb.NewTxDataAccess[UserEntity](db, createUserEntity)
	service := BuildService[UserEntity, UserQuery](
		"/user/", userDataAccess,
		createUserEntity,
		func() UserQuery { return UserQuery{} },
	)
	rs := &RestService[UserEntity, UserQuery]{Service: service}

	t.Run("Page /user/", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?PageNumber=1&PageSize=2", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"list":[{"id":1,"score":85,"memo":"Good"},{"id":2,"score":40,"memo":"Bad"}],"total":4},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Should return empty array instead of null when no data found.", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?PageNumber=10", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"list":[],"total":4},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?ScoreLt=60&test=test", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"list":[{"id":2,"score":40,"memo":"Bad"},{"id":3,"score":55,"memo":null}],"total":2},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/ ", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?ScoreLt=60&MemoNull=true", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"list":[{"id":3,"score":55,"memo":null}],"total":1},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/?MemoLike=%oo%", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?MemoLike=%25oo%25", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"list":[{"id":1,"score":85,"memo":"Good"}],"total":1},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/?IdIn=1,4", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/?IdIn=1,4", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"list":[{"id":1,"score":85,"memo":"Good"},{"id":4,"score":62,"memo":"Well"}],"total":2},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/1", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/1", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":{"id":1,"score":85,"memo":"Good"},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Get /user/100", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/user/100", nil)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"success":false,"error":"record not found. id: 100"}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("PUT /user/1", func(t *testing.T) {
		tc := tm.StartTransaction(ctx)
		defer tc.Rollback()

		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"score":90}`)
		request := httptest.NewRequest("PUT", "/user/1", body).WithContext(tc)

		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}

		writer = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/user/1", nil).WithContext(tc)

		rs.ServeHTTP(writer, request)

		actual = writer.Body.String()
		expect = `{"data":{"id":1,"score":90,"memo":null},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("DELETE /user/{id}", func(t *testing.T) {
		tc := tm.StartTransaction(ctx)
		defer tc.Rollback()
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/user/1", nil).WithContext(tc)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}

		writer = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/user/", nil).WithContext(tc)
		rs.ServeHTTP(writer, request)

		pageList := core.PageList[UserEntity]{}
		response := core.Response{Data: &pageList}
		_ = json.Unmarshal(writer.Body.Bytes(), &response)
		if pageList.Total != 3 {
			t.Errorf("\nExpected: %d\nBut got : %v", 3, pageList.Total)
		}
	})

	t.Run("PATCH /user/{id}", func(t *testing.T) {
		tc := tm.StartTransaction(ctx)
		defer tc.Rollback()

		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"score":33}`)
		request := httptest.NewRequest("PATCH", "/user/1", body).WithContext(tc)
		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}

		writer = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/user/1", nil).WithContext(tc)
		rs.ServeHTTP(writer, request)

		actual = writer.Body.String()
		expect = `{"data":{"id":1,"score":33,"memo":"Good"},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("POST /user/", func(t *testing.T) {
		tc := tm.StartTransaction(ctx)
		defer tc.Rollback()
		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`[{"score":60, "memo":"Well"}]`)
		request := httptest.NewRequest("POST", "/user/", body).WithContext(tc)
		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}

		writer = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/user/5", nil).WithContext(tc)
		rs.ServeHTTP(writer, request)

		actual = writer.Body.String()
		expect = `{"data":{"id":5,"score":60,"memo":"Well"},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
