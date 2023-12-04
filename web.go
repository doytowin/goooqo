package goquery

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func (s *Service[E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := s.createQuery()
	queryMap := request.URL.Query()

	for name, v := range queryMap {
		field := reflect.ValueOf(query).Elem().FieldByName(name)
		if field.IsValid() {
			if field.Kind() == reflect.Bool {
				v0 := strings.EqualFold(v[0], "TRue")
				field.Set(reflect.ValueOf(v0))
			} else {
				v0, err := strconv.Atoi(v[0])
				if noError(err) {
					field.Set(reflect.ValueOf(&v0))
				}
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
