package web

import (
	"encoding/json"
	"fmt"
	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

type restService[E Entity, Q Query] struct {
	DataAccess[E]
	idRgx *regexp.Regexp
}

func NewRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[E],
) http.Handler {
	return &restService[E, Q]{
		DataAccess: dataAccess,
		idRgx:      regexp.MustCompile(prefix + `([\da-fA-F]+)$`),
	}
}

func (s *restService[E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match := s.idRgx.FindStringSubmatch(request.URL.Path)
	var data any
	var err error
	if len(match) > 0 {
		// Process requests for /<model>/{id}
		id := match[1]
		data, err = s.process(request, id)
		writeResult(writer, err, data)
		return
	}
	if request.Method == "POST" {
		body, _ := io.ReadAll(request.Body)
		var entity []E
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			data, err = s.CreateMulti(request.Context(), entity)
		}
		writeResult(writer, err, data)
		return
	}

	query := *new(Q)
	queryMap := request.URL.Query()
	resolveQuery(queryMap, &query)
	if request.Method == "PATCH" {
		body, _ := io.ReadAll(request.Body)
		var entity E
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			data, err = s.PatchByQuery(request.Context(), entity, query)
		}
		writeResult(writer, err, data)
		return
	}
	pageList, err := s.Page(request.Context(), query)
	writeResult(writer, err, pageList)
}

func (s *restService[E, Q]) process(request *http.Request, id string) (any, error) {
	var err error
	var data any
	switch request.Method {
	case "PUT":
		body, _ := io.ReadAll(request.Body)
		entity := *new(E)
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			entity.SetId(&entity, id)
			return s.Update(request.Context(), entity)
		}
	case "PATCH":
		body, _ := io.ReadAll(request.Body)
		entity := *new(E)
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			entity.SetId(&entity, id)
			return s.Patch(request.Context(), entity)
		}
	case "DELETE":
		return s.Delete(request.Context(), id)
	default:
		var entity *E
		entity, err = s.Get(request.Context(), id)
		if entity == nil && NoError(err) {
			err = fmt.Errorf("record not found. id: %s", id)
		} else {
			data = entity
		}
	}
	return data, err
}

func resolveQuery(queryMap url.Values, query any) {
	elem := reflect.ValueOf(query).Elem()
	for name, v := range queryMap {
		path := strings.Split(name, ".")
		field := elem.FieldByName(path[0])
		for i := 1; i < len(path); i++ {
			if !field.IsValid() {
				break
			}
			if field.IsNil() {
				fieldType := field.Type().Elem()
				field.Set(reflect.New(fieldType))
			}
			field = field.Elem().FieldByName(path[i])
		}

		if field.IsValid() {
			convertAndSet(field, v)
		}
	}
}

func convertAndSet(field reflect.Value, v []string) {
	log.Debug("field.Type: ", field.Type())
	fieldType := field.Type()
	v0, err := converterMap[fieldType](v)
	if NoError(err) || v0 != nil {
		field.Set(reflect.ValueOf(v0))
	}
}

func writeResult(writer http.ResponseWriter, err error, data any) {
	response := Response{Data: data, Success: NoError(err), Error: ReadError(err)}
	marshal, err := json.Marshal(response)
	if NoError(err) {
		writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = writer.Write(marshal)
	}
}
