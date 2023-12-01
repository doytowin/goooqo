package goquery

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

func (s *Service[E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	query := createModel[Q](s.queryType)
	queryMap := request.URL.Query()
	pageQuery := query.GetPageQuery()

	sizeStr := queryMap.Get("size")
	size, err := strconv.Atoi(sizeStr)
	if noError(err) {
		pageQuery.PageSize = &size
	}
	pageStr := queryMap.Get("page")
	page, err := strconv.Atoi(pageStr)
	if noError(err) {
		pageQuery.PageNumber = &page
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
