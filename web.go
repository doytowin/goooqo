package goquery

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"reflect"
	"strconv"
)

func (s *Service[E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := s.createQuery()
	queryMap := request.URL.Query()

	for name, v := range queryMap {
		field := reflect.ValueOf(query).Elem().FieldByName(name)
		if field.IsValid() {
			integer, err := strconv.Atoi(v[0])
			if noError(err) {
				field.Set(reflect.ValueOf(&integer))
			}
		}
	}

	pageList, err := s.Page(query)
	if noError(err) {
		marshal, err := json.Marshal(pageList)
		if noError(err) {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(marshal)
		}
	}

}
