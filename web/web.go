package web

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type restService[E Entity, Q Query] struct {
	DataAccess[context.Context, E]
	createQuery  func() Q
	createEntity func() E
	idRgx        *regexp.Regexp
}

func NewRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[context.Context, E],
	createEntity func() E,
	createQuery func() Q,
) http.Handler {
	return &restService[E, Q]{
		DataAccess:   dataAccess,
		createQuery:  createQuery,
		createEntity: createEntity,
		idRgx:        regexp.MustCompile(prefix + `(\d+)$`),
	}
}

func (s *restService[E, Q]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match := s.idRgx.FindStringSubmatch(request.URL.Path)
	var data any
	var err error
	if len(match) > 0 {
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

	query := s.createQuery()
	queryMap := request.URL.Query()
	resolveQuery(queryMap, &query)
	pageList, err := s.Page(request.Context(), query)
	writeResult(writer, err, pageList)
}

func (s *restService[E, Q]) process(request *http.Request, id string) (any, error) {
	var err error
	var data any
	switch request.Method {
	case "PUT":
		body, _ := io.ReadAll(request.Body)
		entity := s.createEntity()
		err = json.Unmarshal(body, &entity)
		if NoError(err) {
			entity.SetId(&entity, id)
			return s.Update(request.Context(), entity)
		}
	case "PATCH":
		body, _ := io.ReadAll(request.Body)
		entity := s.createEntity()
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
		field := elem.FieldByName(name)
		if !field.IsValid() {
			continue
		}
		log.Debug("field.Kind: ", field.Kind())
		if field.Kind() == reflect.Bool {
			v0 := strings.EqualFold(v[0], "TRue")
			field.Set(reflect.ValueOf(v0))
		} else if field.Kind() == reflect.Pointer {
			resolvePointer(field, v)
		}
	}
}

func resolvePointer(field reflect.Value, v []string) {
	log.Debug("field.Type: ", field.Type())
	fieldType := field.Type().String()
	if fieldType == "*[]int" {
		strArr := strings.Split(v[0], ",")
		var v0 []int
		for _, s := range strArr {
			num, err := strconv.Atoi(s)
			if NoError(err) {
				v0 = append(v0, num)
			}
		}
		field.Set(reflect.ValueOf(&v0))
	} else if fieldType == "*int" {
		v0, err := strconv.Atoi(v[0])
		if NoError(err) {
			field.Set(reflect.ValueOf(&v0))
		}
	} else if fieldType == "int" {
		v0, err := strconv.Atoi(v[0])
		if NoError(err) {
			field.Set(reflect.ValueOf(v0))
		}
	} else {
		field.Set(reflect.ValueOf(&v[0]))
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
