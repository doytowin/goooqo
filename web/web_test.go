/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/doytowin/goooqo/core"
	"github.com/doytowin/goooqo/rdb"
	. "github.com/doytowin/goooqo/test"
	log "github.com/sirupsen/logrus"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestWeb(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	db := rdb.Connect("app.properties")
	InitDB(db)
	defer rdb.Disconnect(db)

	ctx := context.Background()
	tm := rdb.NewTransactionManager(db)

	userDataAccess := rdb.NewTxDataAccess[UserEntity](tm)
	rs := NewRestService[UserEntity, UserQuery]("/user/", userDataAccess)

	t.Run("Should return empty array instead of null when no data found.", func(t *testing.T) {
		tm.SubmitTransaction(ctx, func(tc core.TransactionContext) error {
			writer := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/user/?PageNumber=10", nil).WithContext(tc)

			rs.ServeHTTP(writer, request)

			actual := writer.Body.String()
			expect := `{"data":{"list":[],"total":4},"success":true}`
			if actual != expect {
				t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			}
			return nil
		})
	})

	tests := []struct{ method, url, expect string }{
		{"Get", "/user/?PageNumber=1&PageSize=2", `{"data":{"list":[{"id":1,"score":85,"memo":"Good"},{"id":2,"score":40,"memo":"Bad"}],"total":4},"success":true}`},
		{"Get", "/user/?ScoreLt=60&test=test", `{"data":{"list":[{"id":2,"score":40,"memo":"Bad"},{"id":3,"score":55,"memo":null}],"total":2},"success":true}`},
		{"Get", "/user/?ScoreLt=60&Sort=id,desc", `{"data":{"list":[{"id":3,"score":55,"memo":null},{"id":2,"score":40,"memo":"Bad"}],"total":2},"success":true}`},
		{"Get", "/user/?ScoreLt=60&MemoNull=true", `{"data":{"list":[{"id":3,"score":55,"memo":null}],"total":1},"success":true}`},
		{"Get", "/user/?MemoLike=%25oo%25", `{"data":{"list":[{"id":1,"score":85,"memo":"Good"}],"total":1},"success":true}`},
		{"Get", "/user/?IdIn=1,4", `{"data":{"list":[{"id":1,"score":85,"memo":"Good"},{"id":4,"score":62,"memo":"Well"}],"total":2},"success":true}`},
		{"Get", "/user/?IdIn=1&IdIn=4&IdIn=a5", `{"data":{"list":[{"id":1,"score":85,"memo":"Good"},{"id":4,"score":62,"memo":"Well"}],"total":2},"success":true}`},
		{"Get", "/user/1", `{"data":{"id":1,"score":85,"memo":"Good"},"success":true}`},
		{"Get", "/user/100", `{"success":false,"error":"record not found. id: 100"}`},
	}
	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			writer := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.url, nil)

			rs.ServeHTTP(writer, request)
			actual := writer.Body.String()

			if actual != test.expect {
				t.Errorf("\nExpected: %s\nBut got : %s", test.expect, actual)
			}
		})
	}

	t.Run("PUT /user/1", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"score":90}`)
		request := httptest.NewRequest("PUT", "/user/1", body).WithContext(tc)

		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Fatalf("\nExpected: %s\nBut got : %s", expect, actual)
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
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/user/1", nil).WithContext(tc)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Fatalf("\nExpected: %s\nBut got : %s", expect, actual)
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

	t.Run("DELETE /user/?MemoNull=false", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()
		writer := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/user/?MemoNull=false", nil).WithContext(tc)

		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":3,"success":true}`
		if actual != expect {
			t.Fatalf("\nExpected: %s\nBut got : %s", expect, actual)
		}

		writer = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/user/", nil).WithContext(tc)
		rs.ServeHTTP(writer, request)

		pageList := core.PageList[UserEntity]{}
		response := core.Response{Data: &pageList}
		_ = json.Unmarshal(writer.Body.Bytes(), &response)
		if pageList.Total != 1 {
			t.Errorf("\nExpected: %d\nBut got : %v", 3, pageList.Total)
		}
	})

	t.Run("PATCH /user/{id}", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"score":33}`)
		request := httptest.NewRequest("PATCH", "/user/1", body).WithContext(tc)
		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Fatalf("\nExpected: %s\nBut got : %s", expect, actual)
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

	t.Run("PATCH /user/?memoNull=true", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"memo":"New Memo"}`)
		request := httptest.NewRequest("PATCH", "/user/?memoNull=true", body).WithContext(tc)
		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Fatalf("\nExpected: %s\nBut got : %s", expect, actual)
		}

		writer = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/user/?memoNull=true", nil).WithContext(tc)
		rs.ServeHTTP(writer, request)

		actual = writer.Body.String()
		expect = `{"data":{"list":[],"total":0},"success":true}`
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("POST /user/", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()
		writer := httptest.NewRecorder()
		body := bytes.NewBufferString(`[{"score":60, "memo":"Well"}]`)
		request := httptest.NewRequest("POST", "/user/", body).WithContext(tc)
		request.Header.Set("content-type", "application/json; charset=utf-8")
		rs.ServeHTTP(writer, request)

		actual := writer.Body.String()
		expect := `{"data":1,"success":true}`
		if actual != expect {
			t.Fatalf("\nExpected: %s\nBut got : %s", expect, actual)
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

func Test_resolveQuery(t *testing.T) {
	type Unit struct {
		Name *string `json:"name,omitempty"`
	}
	type SizeQuery struct {
		HLt  *int  `json:"hLt,omitempty"`
		HGe  *int  `json:"hGe,omitempty"`
		Unit *Unit `json:"unit,omitempty"`
	}
	type qo struct {
		Size *SizeQuery `json:"size,omitempty"`
	}
	type args struct {
		queryMap url.Values
		query    qo
	}
	tests := []struct {
		name   string
		expect string
		args   args
	}{
		{"Nested Parameters", `{"size":{"hLt":20}}`,
			args{url.Values{"Size.HLt": {"20"}}, qo{}}},
		{"Two Nested Parameters", `{"size":{"hLt":20,"hGe":10}}`,
			args{url.Values{"Size.HLt": {"20"}, "Size.HGe": {"10"}}, qo{}}},
		{"Level Three of Nested Parameters", `{"size":{"unit":{"name":"cm"}}}`,
			args{url.Values{"Size.Unit.Name": {"cm"}}, qo{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolveQuery(tt.args.queryMap, &tt.args.query)
			data, _ := json.Marshal(tt.args.query)
			actual := string(data)
			if tt.expect != actual {
				t.Errorf("\nExpected: %s\nBut got : %s", tt.expect, actual)
			}
		})
	}
}
